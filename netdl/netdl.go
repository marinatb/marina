package netdl

import (
	"fmt"
	"github.com/marinatb/marina"
	"github.com/satori/go.uuid"
	"strings"
)

type Endpoint interface {
	GetNetHost() *NetHost
}

type NetHost struct {
	Id         uuid.UUID             `json:"-"`
	Name       string                `json:"name"`
	Interfaces map[string]*Interface `json:"interfaces"`
}

func (h *NetHost) AddInterface(latency, capacity uint) *Interface {
	ifx := new(Interface)
	ifx.Id = uuid.NewV4()
	ifx.Name = fmt.Sprintf("eth%d", len(h.Interfaces))
	ifx.Latency = latency
	ifx.Capacity = capacity
	h.Interfaces[ifx.Id.String()] = ifx
	return ifx
}

type Interface struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	PacketConductor
}

type Computer struct {
	NetHost
	Container string `json:"container"`
}

func (c *Computer) String() string {
	return fmt.Sprintf("%s %s", c.Name, c.Container)
}

func (c *Computer) GetNetHost() *NetHost {
	return &c.NetHost
}

type PacketConductor struct {
	Capacity uint `json:"capacity"`
	Latency  uint `json:"latency"`
}

type Switch struct {
	NetHost
	PacketConductor
	Net       *Network `json:"-"`
	Endpoints map[string]NetIfRef
}

func (sw *Switch) String() string {
	s := fmt.Sprintf("%s %dmbps [", sw.Name, sw.Capacity)
	for _, ref := range sw.Endpoints {
		s += sw.Net.ResolveEndpoint(ref).GetNetHost().Name + ","
	}
	return strings.TrimSuffix(s, ",") + "]"
}

func (sw *Switch) Connect(endpoints ...Endpoint) {
	for _, e := range endpoints {
		ifx := e.GetNetHost().AddInterface(sw.Capacity, 0)
		sw.Endpoints[e.GetNetHost().Id.String()] = NetIfRef{
			e.GetNetHost().Id, ifx.Id, marina.SimpleTypename(e)}
	}
}

func (s *Switch) GetNetHost() *NetHost {
	return &s.NetHost
}

type Router struct {
	NetHost
	PacketConductor
}

func (r *Router) String() string {
	return fmt.Sprintf("%s %dms %dmbps", r.Name, r.Latency, r.Capacity)
}

func (r *Router) GetNetHost() *NetHost {
	return &r.NetHost
}

type NetIfRef struct {
	Id   uuid.UUID `json:"id"`
	IfId uuid.UUID `json:"ifname"`
	Type string    `json:"type"`
}

type Link struct {
	Id        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Net       *Network    `json:"-"`
	Endpoints [2]NetIfRef `json:"endpoints"`
	PacketConductor
}

func (l *Link) String() string {
	s := fmt.Sprintf("%s %dms %dmbps [", l.Name, l.Latency, l.Capacity)
	for _, ref := range l.Endpoints {
		s += l.Net.ResolveEndpoint(ref).GetNetHost().Name + ","
	}
	return strings.TrimSuffix(s, ",") + "]"
}

type Network struct {
	Id        uuid.UUID            `json:"id"`
	Name      string               `json:"name"`
	Computers map[string]*Computer `json:"computers"`
	Routers   map[string]*Router   `json:"routers"`
	Switches  map[string]*Switch   `json:"switches"`
	Links     map[string]*Link     `json:"links"`
}

func NewNetwork(name string) *Network {
	net := new(Network)
	net.Id = uuid.NewV4()
	net.Name = name
	net.Computers = make(map[string]*Computer)
	net.Routers = make(map[string]*Router)
	net.Switches = make(map[string]*Switch)
	net.Links = make(map[string]*Link)
	return net
}

func (net *Network) String() string {

	s := fmt.Sprintf("Network: %s\n", net.Name)
	s += "  Computers:\n"
	for _, c := range net.Computers {
		s += fmt.Sprintf("    %v\n", c)
	}
	s += "  Routers:\n"
	for _, r := range net.Routers {
		s += fmt.Sprintf("    %v\n", r)
	}
	s += "  Switches:\n"
	for _, sw := range net.Switches {
		s += fmt.Sprintf("    %v\n", sw)
	}
	s += "  Links:\n"
	for _, l := range net.Links {
		s += fmt.Sprintf("    %v\n", l)
	}

	return s
}

func (net *Network) Init() {
	for k, _ := range net.Switches {
		net.Switches[k].Net = net
	}
	for k, _ := range net.Links {
		net.Links[k].Net = net
	}
}

func (net *Network) NewRouter(name string, capacity, latency uint) *Router {
	router := new(Router)
	router.Id = uuid.NewV4()
	router.Name = name
	router.Interfaces = make(map[string]*Interface)
	router.Latency = latency
	router.Capacity = capacity

	net.Routers[router.Id.String()] = router

	return router
}

func (net *Network) NewLink(a, b Endpoint, name string, capacity, latency uint) *Link {
	link := new(Link)
	link.Id = uuid.NewV4()
	link.Name = name
	link.Net = net
	link.Latency = latency
	link.Capacity = capacity

	ifxa := a.GetNetHost().AddInterface(capacity, latency)
	ifxb := b.GetNetHost().AddInterface(capacity, latency)

	link.Endpoints[0] = NetIfRef{a.GetNetHost().Id, ifxa.Id, marina.SimpleTypename(a)}
	link.Endpoints[1] = NetIfRef{b.GetNetHost().Id, ifxb.Id, marina.SimpleTypename(b)}

	net.Links[link.Id.String()] = link

	return link
}

func (net *Network) NewSwitch(name string, capacity uint) *Switch {
	sw := new(Switch)
	sw.Net = net
	sw.Id = uuid.NewV4()
	sw.Name = name
	sw.Latency = 0
	sw.Capacity = capacity
	sw.Endpoints = make(map[string]NetIfRef)
	net.Switches[sw.Id.String()] = sw
	return sw
}

func (net *Network) NewComputer(name string) *Computer {
	c := new(Computer)
	c.Id = uuid.NewV4()
	c.Name = name
	c.Container = "Ubuntu-15.04-Base"
	c.Interfaces = make(map[string]*Interface)
	net.Computers[c.Id.String()] = c
	return c
}

func (net *Network) ResolveEndpoint(r NetIfRef) Endpoint {
	switch r.Type {
	case "Computer":
		return net.Computers[r.Id.String()]
	case "Switch":
		return net.Switches[r.Id.String()]
	case "Router":
		return net.Routers[r.Id.String()]
	default:
		return nil
	}
}
