package testsite

import (
	"fmt"
	"github.com/marinatb/marina/netdl"
)

func Build() *netdl.Network {

	net := netdl.NewNetwork("testsite")
	nSwitches, nodesPerSwitch := 4, 4

	var _sw *netdl.PSwitch = nil
	for i := 0; i < nSwitches; i++ {
		sw := net.NewPSwitch(fmt.Sprintf("sw%d", i), 10000)
		for j := 0; j < nodesPerSwitch; j++ {
			pn := net.NewPNode(fmt.Sprintf("node%d", i*nSwitches+j))
			net.NewPLink(sw, pn, fmt.Sprintf("sw%d_node%d", i, i*nSwitches+j), 1000, 0)
		}
		if _sw != nil {
			net.NewPLink(sw, sw, fmt.Sprintf("sw%d_sw%d", i, i-1), 10000, 0)
		}
		_sw = sw
	}

	return net
}
