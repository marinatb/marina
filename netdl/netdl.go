package netdl

import (
	"fmt"
	"github.com/marinatb/marina"
	"github.com/satori/go.uuid"
	"strings"
)

func Fingerprint(id *uuid.UUID) string {
	s := id.String()
	return s[len(s)-8:]
}

type Endpoint interface {
	GetNetHost() *NetHost
}
type PEndpoint interface {
	GetPNetHost() *NetHost
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
	ifx.Link = uuid.Nil
	h.Interfaces[ifx.Id.String()] = ifx
	return ifx
}

func (h *NetHost) PlugInterface(ifid uuid.UUID, link *Link) {
	h.Interfaces[ifid.String()].Link = link.Id
}

type Interface struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Link uuid.UUID `json:"link"`
	PacketConductor
}

func (ifx *Interface) Plug(link *Link) {
	ifx.Link = link.Id
}
func (ifx *Interface) PPlug(link *PLink) {
	ifx.Link = link.Id
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

type PNode struct {
	NetHost
	NetId    uuid.UUID `json:"netid"`
	Elements []NetRef  `json:"elements"`
}

func (p *PNode) String() string {
	s := p.Name + " " + Fingerprint(&p.NetId) + " ["
	for _, nr := range p.Elements {
		s += Fingerprint(&nr.Id) + ","
	}

	return strings.TrimSuffix(s, ",") + "]"
}

func (p *PNode) GetPNetHost() *NetHost {
	return &p.NetHost
}

type PacketConductor struct {
	Capacity uint `json:"capacity"`
	Latency  uint `json:"latency"`
}

type Switch struct {
	NetHost
	PacketConductor
	Net *Network `json:"-"`
}

func (sw *Switch) String() string {
	s := fmt.Sprintf("%s %dmbps { ", sw.Name, sw.Capacity)
	for _, ifx := range sw.Interfaces {
		s += ifx.Name + "[" + sw.Net.ResolveLink(ifx.Link).Name + "], "
	}
	return strings.TrimSuffix(s, ", ") + " }"
}

func (sw *Switch) GetNetHost() *NetHost {
	return &sw.NetHost
}

type PSwitch struct {
	NetHost
	Net               *Network `json:"-"`
	AllocatedCapacity uint     `json:"allocated_capacity"`
	PacketConductor
}

func (sw *PSwitch) String() string {
	s := fmt.Sprintf("%s %dmbps %dmbps { ", sw.Name, sw.Capacity, sw.AllocatedCapacity)
	for _, ifx := range sw.Interfaces {
		s += ifx.Name + "[" + sw.Net.ResolvePLink(ifx.Link).Name + "], "
	}
	return strings.TrimSuffix(s, ", ") + " }"
}

func (s *PSwitch) GetPNetHost() *NetHost {
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

type NetRef struct {
	Id   uuid.UUID `json:"id"`
	Type string    `json:"type"`
}

type NetIfRef struct {
	NetRef
	IfId uuid.UUID `json:"ifname"`
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

type PLink struct {
	Id        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Net       *Network    `json:"-"`
	Endpoints [2]NetIfRef `json:"endpoints"`
	PacketConductor
}

func (l *PLink) String() string {
	s := fmt.Sprintf("%s %dms %dmbps [", l.Name, l.Latency, l.Capacity)
	for _, ref := range l.Endpoints {
		s += l.Net.ResolvePEndpoint(ref).GetPNetHost().Name + ","
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
	PNodes    map[string]*PNode    `json:"pnodes"`
	PSwitches map[string]*PSwitch  `json:"pswitches"`
	PLinks    map[string]*PLink    `json:"plinks"`
}

func NewNetwork(name string) *Network {
	net := new(Network)
	net.Id = uuid.NewV4()
	net.Name = name
	net.Computers = make(map[string]*Computer)
	net.Routers = make(map[string]*Router)
	net.Switches = make(map[string]*Switch)
	net.Links = make(map[string]*Link)
	net.PNodes = make(map[string]*PNode)
	net.PSwitches = make(map[string]*PSwitch)
	net.PLinks = make(map[string]*PLink)
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
	s += "  PNodes:\n"
	for _, p := range net.PNodes {
		s += fmt.Sprintf("    %v\n", p)
	}
	s += "  PSwitches:\n"
	for _, sw := range net.PSwitches {
		s += fmt.Sprintf("    %v\n", sw)
	}
	s += "  PLinks:\n"
	for _, l := range net.PLinks {
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
	for k, _ := range net.PSwitches {
		net.PSwitches[k].Net = net
	}
	for k, _ := range net.PLinks {
		net.PLinks[k].Net = net
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
	ifxa.Plug(link)
	ifxb := b.GetNetHost().AddInterface(capacity, latency)
	ifxb.Plug(link)

	nir := NetIfRef{}
	nir.Id = a.GetNetHost().Id
	nir.IfId = ifxa.Id
	nir.Type = marina.SimpleTypename(a)
	link.Endpoints[0] = nir

	nir = NetIfRef{}
	nir.Id = b.GetNetHost().Id
	nir.IfId = ifxb.Id
	nir.Type = marina.SimpleTypename(b)
	link.Endpoints[1] = nir

	net.Links[link.Id.String()] = link

	return link
}

func (net *Network) NewPLink(a, b PEndpoint, name string, capacity, latency uint) *PLink {
	link := new(PLink)
	link.Id = uuid.NewV4()
	link.Name = name
	link.Net = net
	link.Latency = latency
	link.Capacity = capacity

	ifxa := a.GetPNetHost().AddInterface(capacity, latency)
	ifxa.PPlug(link)
	ifxb := b.GetPNetHost().AddInterface(capacity, latency)
	ifxb.PPlug(link)

	nir := NetIfRef{}
	nir.Id = a.GetPNetHost().Id
	nir.IfId = ifxa.Id
	nir.Type = marina.SimpleTypename(a)
	link.Endpoints[0] = nir

	nir = NetIfRef{}
	nir.Id = b.GetPNetHost().Id
	nir.IfId = ifxb.Id
	nir.Type = marina.SimpleTypename(b)
	link.Endpoints[1] = nir

	net.PLinks[link.Id.String()] = link

	return link
}

func (net *Network) NewSwitch(name string, capacity uint) *Switch {
	sw := new(Switch)
	sw.Net = net
	sw.Id = uuid.NewV4()
	sw.Name = name
	sw.Latency = 0
	sw.Capacity = capacity
	sw.Interfaces = make(map[string]*Interface)
	net.Switches[sw.Id.String()] = sw
	return sw
}

func (net *Network) NewPSwitch(name string, capacity uint) *PSwitch {
	sw := new(PSwitch)
	sw.Net = net
	sw.Id = uuid.NewV4()
	sw.Name = name
	sw.Latency = 0
	sw.Capacity = capacity
	sw.Interfaces = make(map[string]*Interface)
	net.PSwitches[sw.Id.String()] = sw
	return sw
}

func (net *Network) NewPNode(name string) *PNode {
	pn := new(PNode)
	pn.Id = uuid.NewV4()
	pn.Name = name
	pn.NetId = uuid.Nil
	pn.Interfaces = make(map[string]*Interface)
	net.PNodes[pn.Id.String()] = pn
	return pn
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

func (net *Network) ResolveLink(id uuid.UUID) *Link {
	if id == uuid.Nil {
		return nil
	}
	link, _ := net.Links[id.String()]

	return link
}

func (net *Network) ResolvePLink(id uuid.UUID) *PLink {
	if id == uuid.Nil {
		return nil
	}
	link, _ := net.PLinks[id.String()]

	return link
}

func (net *Network) ResolvePEndpoint(r NetIfRef) PEndpoint {
	switch r.Type {
	case "PNode":
		return net.PNodes[r.Id.String()]
	case "PSwitch":
		return net.PSwitches[r.Id.String()]
	default:
		return nil
	}
}
