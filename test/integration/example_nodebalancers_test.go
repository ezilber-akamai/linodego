package integration

import (
	"context"
	"fmt"
	"log"

	"github.com/linode/linodego"
)

func ExampleClient_CreateNodeBalancer() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleCreateNodeBalancer")
	defer teardown()

	fmt.Println("## NodeBalancer create")
	var nbID int
	nb := &linodego.NodeBalancer{
		ClientConnThrottle: 20,
		Region:             "us-southeast",
	}

	createOpts := nb.GetCreateOptions()
	nb, err := linodeClient.CreateNodeBalancer(context.Background(), createOpts)
	if err != nil {
		log.Fatal(err)
	}
	nbID = nb.ID

	fmt.Println("### Get")
	nb, err = linodeClient.GetNodeBalancer(context.Background(), nbID)
	if err != nil {
		log.Fatal(err)
	}

	updateOpts := nb.GetUpdateOptions()
	*updateOpts.Label += "_renamed"
	nb, err = linodeClient.UpdateNodeBalancer(context.Background(), nbID, updateOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Delete")
	if err := linodeClient.DeleteNodeBalancer(context.Background(), nbID); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ## NodeBalancer create
	// ### Get
	// ### Delete
}

func ExampleClient_CreateNodeBalancerConfig() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleCreateNodeBalancerConfig")
	defer teardown()

	fmt.Println("## NodeBalancer create")
	nb, err := linodeClient.CreateNodeBalancer(context.Background(), linodego.NodeBalancerCreateOptions{
		ClientConnThrottle: linodego.Pointer(20),
		Region:             linodego.Pointer("us-southeast"),
		FirewallID:         linodego.Pointer(GetFirewallID()),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("## NodeBalancer Config create")
	createOpts := linodego.NodeBalancerConfigCreateOptions{
		Port: 80,
		/*
		   Protocol:      linodego.ProtocolHTTP,
		   Algorithm:     linodego.AlgorithmLeastConn,
		   Stickiness:    linodego.StickinessHTTPCookie,
		   Check:         linodego.CheckHTTP,
		   CipherSuite:   linodego.CipherRecommended,
		*/
		CheckTimeout:  linodego.Pointer(5),
		CheckInterval: linodego.Pointer(30),
		CheckAttempts: linodego.Pointer(5),
	}
	nbc, err := linodeClient.CreateNodeBalancerConfig(context.Background(), nb.ID, createOpts)
	if err != nil {
		log.Fatal(err)
	}
	nbcID := nbc.ID

	fmt.Println("## NodeBalancer Config update")
	updateOpts := nbc.GetUpdateOptions()
	updateOpts.Port += 8000
	nbc, err = linodeClient.UpdateNodeBalancerConfig(context.Background(), nb.ID, nbc.ID, updateOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### List")
	configs, err := linodeClient.ListNodeBalancerConfigs(context.Background(), nb.ID, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Get")
	nbc, err = linodeClient.GetNodeBalancerConfig(context.Background(), nb.ID, configs[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Delete")
	if nbc.ID != nbcID {
		log.Fatalf("Unexpected Nodebalancer Config ID %d != %d", nbc.ID, nbcID)
	}
	if err := linodeClient.DeleteNodeBalancerConfig(context.Background(), nb.ID, nbc.ID); err != nil {
		log.Fatal(err)
	}

	if err := linodeClient.DeleteNodeBalancer(context.Background(), nb.ID); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ## NodeBalancer create
	// ## NodeBalancer Config create
	// ## NodeBalancer Config update
	// ### List
	// ### Get
	// ### Delete
}

func ExampleClient_CreateNodeBalancerNode() {
	// Example readers, Ignore this bit of setup code needed to record test fixtures
	linodeClient, teardown := createTestClient(nil, "fixtures/ExampleCreateNodeBalancerNode")
	defer teardown()

	fmt.Println("## NodeBalancer create")
	nb, err := linodeClient.CreateNodeBalancer(context.Background(), linodego.NodeBalancerCreateOptions{
		ClientConnThrottle: linodego.Pointer(20),
		Region:             linodego.Pointer("us-southeast"),
		FirewallID:         linodego.Pointer(GetFirewallID()),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("## NodeBalancer Config create")

	nbc, err := linodeClient.CreateNodeBalancerConfig(context.Background(), nb.ID, linodego.NodeBalancerConfigCreateOptions{
		Port: 80,
	})
	if err != nil {
		log.Fatal(err)
	}

	instanceOpts := linodego.InstanceCreateOptions{
		Label:      linodego.Pointer("nodebalancer-example-instance"),
		RootPass:   linodego.Pointer(randPassword()),
		Region:     "us-southeast",
		Type:       "g6-nanode-1",
		Image:      linodego.Pointer("linode/debian9"),
		Booted:     linodego.Pointer(false),
		FirewallID: linodego.Pointer(GetFirewallID()),
	}
	instance, err := linodeClient.CreateInstance(context.Background(), instanceOpts)
	if err != nil {
		log.Fatal(err)
	}

	opts := linodego.InstanceIPAddOptions{
		Public: false,
	}
	ip, err := linodeClient.AddInstanceIPAddress(context.Background(), instance.ID, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("## NodeBalancer Node create")
	nodeOpts := linodego.NodeBalancerNodeCreateOptions{
		Address: fmt.Sprintf("%s:%d", ip.Address, nbc.Port),
		Label:   "node-" + ip.Address,
	}
	nbn, err := linodeClient.CreateNodeBalancerNode(context.Background(), nb.ID, nbc.ID, nodeOpts)
	if err != nil {
		log.Fatal(err)
	}
	nbnID := nbn.ID

	fmt.Println("## NodeBalancer Node update")
	updateOpts := nbn.GetUpdateOptions()
	updateOpts.Address = linodego.Pointer(fmt.Sprintf("%s:8080", ip.Address))
	nbn, err = linodeClient.UpdateNodeBalancerNode(context.Background(), nb.ID, nbc.ID, nbn.ID, updateOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### List")
	nodes, err := linodeClient.ListNodeBalancerNodes(context.Background(), nb.ID, nbc.ID, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Get")
	nbn, err = linodeClient.GetNodeBalancerNode(context.Background(), nb.ID, nbc.ID, nodes[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Delete")
	if nbn.ID != nbnID {
		log.Fatalf("Unexpected Nodebalancer Node ID %d != %d", nbn.ID, nbnID)
	}
	if err := linodeClient.DeleteNodeBalancerNode(context.Background(), nb.ID, nbc.ID, nbn.ID); err != nil {
		log.Fatal(err)
	}

	if err := linodeClient.DeleteNodeBalancerConfig(context.Background(), nb.ID, nbc.ID); err != nil {
		log.Fatal(err)
	}

	if err := linodeClient.DeleteNodeBalancer(context.Background(), nb.ID); err != nil {
		log.Fatal(err)
	}

	if err := linodeClient.DeleteInstance(context.Background(), instance.ID); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ## NodeBalancer create
	// ## NodeBalancer Config create
	// ## NodeBalancer Node create
	// ## NodeBalancer Node update
	// ### List
	// ### Get
	// ### Delete
}
