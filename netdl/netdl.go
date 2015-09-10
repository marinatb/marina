package netdl

import (
	"github.com/satori/go.uuid"
)

type NetHost struct {
	Id         uuid.UUID            `json:"id"`
	Interfaces map[string]Interface `json:"interfaces"`
}

type Interface struct {
	Name string `json:"name"`
	PacketConductor
}

type Computer struct {
	NetHost
	OS           string `json:"os"`
	Start_script string `json:"start_script"`
}

type PacketConductor struct {
	Capacity int `json:"capacity"`
	Latency  int `json:"latency"`
}

type Switch struct {
	NetHost
	PacketConductor
}

type Router struct {
	NetHost
	PacketConductor
}

type NetIfRef struct {
	Id     uuid.UUID `json:"id"`
	IfName string    `json:"ifname"`
}

type Link struct {
	Id uuid.UUID `json:"id"`
	PacketConductor
	Endpoints [2]NetIfRef `json:"endpoints"`
}

type Network struct {
	Id       uuid.UUID                 `json:"id"`
	Name     string                    `json:"name"`
	Elements map[uuid.UUID]interface{} `json:"elements"`
}
