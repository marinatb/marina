package first

import (
	"github.com/marinatb/marina/protocol"
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
