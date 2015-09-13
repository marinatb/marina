package first

import (
	"bytes"
	"github.com/marinatb/marina/protocol"
	"io/ioutil"
	"testing"
)

func TestLocal(t *testing.T) {

	//Build the network
	net := Build()
	t.Log(net)

	//Pack the network into a json file
	js := protocol.PackLegible(net)
	t.Log(string(js))

	//Unpack the network from a json file
	_net, err := protocol.UnpackNetwork(js)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(_net)

}

func TestHarbor(t *testing.T) {

	//Build the network
	net := Build()

	netrq := protocol.NetworkMaterializationRequest{*net, "mcl"}

	//Pack the network into a json file
	js := protocol.PackLegible(netrq)

	//Request Materialization from harbormaster
	client := protocol.InsecureClient()
	resp, err := client.Post("https://localhost:4676/materialize", "application/json",
		bytes.NewBuffer(js))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Log(string(body))
	} else {
		t.Logf("Response Status %d", resp.StatusCode)
		var diag protocol.Diagnostic
		protocol.Unpack(resp.Body, &diag)
		t.Fatal(diag)
	}

}
