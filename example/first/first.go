package first

import (
	"github.com/marinatb/marina/netdl"
)

func Build() *netdl.Network {

	net := netdl.NewNetwork("first")

	ra := net.NewRouter("ra", 100, 1)
	rb := net.NewRouter("rb", 100, 3)
	rc := net.NewRouter("rc", 100, 2)
	/*ra_rb := */ net.NewLink(ra, rb, "ra_rb", 100, 7)
	/*rb_rc := */ net.NewLink(rb, rc, "rb_rc", 100, 11)
	/*rc_ra := */ net.NewLink(rc, ra, "rc_ra", 100, 9)
	sa := net.NewSwitch("sa", 1000)
	sb := net.NewSwitch("sb", 1000)
	sc := net.NewSwitch("sc", 1000)
	a0 := net.NewComputer("a0")
	a1 := net.NewComputer("a1")
	sa.Connect(ra, a0, a1)
	b0 := net.NewComputer("b0")
	b1 := net.NewComputer("b1")
	sb.Connect(rb, b0, b1)
	c0 := net.NewComputer("c0")
	c1 := net.NewComputer("c1")
	sc.Connect(rc, c0, c1)

	return net

}
