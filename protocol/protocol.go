package protocol

import (
	"bytes"
	"encoding/json"
	"github.com/marinatb/marina/netdl"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)

type Diagnostic struct {
	Level, Message string
}

type NetElementRef struct {
	Id   uuid.UUID `json:"id"`
	Type string    `json:"type"`
}

type MarinerConfig struct {
	NodeId   uuid.UUID                   `json:"nodeId"`
	Elements map[uuid.UUID]NetElementRef `json:"elements"`
	LXC      string                      `json:"lxc"`
	OVS      string                      `json:"ovs"`
	Click    string                      `json:"click"`
}

type NetworkMaterializationRequest struct {
	Net    netdl.Network
	Mapper string
}

type MaterializationMap struct {
	Net *netdl.Network
	Map map[uuid.UUID]MarinerConfig `json:"map"`
}

func Unpack(r *http.Request, x interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), &x)
	if err != nil {
		log.Println("[unpack] bad message")
		log.Println(err)
		log.Println(buf.String())
		return nil
	}
	return nil
}

func UnpackNetwork(js []byte) (*netdl.Network, error) {
	net := new(netdl.Network)
	err := json.Unmarshal(js, net)
	if err != nil {
		log.Printf("failed to unmarshal json: %s", err)
		return nil, err
	}
	net.Init()
	return net, nil
}

func pack(x interface{}, pretty bool) []byte {
	var js []byte
	var err error
	if !pretty {
		js, err = json.Marshal(x)
	} else {
		js, err = json.MarshalIndent(x, "", "  ")
	}
	if err != nil {
		log.Printf("[pack] %v", err)
	}
	return js
}

func PackWire(x interface{}) []byte {
	return pack(x, false)
}

func PackLegible(x interface{}) []byte {
	return pack(x, true)
}
