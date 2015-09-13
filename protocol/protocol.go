package protocol

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/marinatb/marina/netdl"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
)

type Diagnostic struct {
	Level, Message string
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s: %s", d.Level, d.Message)
}

type NetElementRef struct {
	Id   uuid.UUID `json:"id"`
	Type string    `json:"type"`
}

type MarinerConfig struct {
	NodeId   uuid.UUID                `json:"nodeId"`
	Elements map[string]NetElementRef `json:"elements"`
	LXC      string                   `json:"lxc"`
	OVS      string                   `json:"ovs"`
	Click    string                   `json:"click"`
}

type NetworkMaterializationRequest struct {
	Net    netdl.Network
	Mapper string
}

type MaterializationEmbedding struct {
	Net *netdl.Network           `json:"net"`
	Map map[string]MarinerConfig `json:"map"`
}

func NewMaterializationEmbedding(net *netdl.Network) *MaterializationEmbedding {
	mm := new(MaterializationEmbedding)
	mm.Net = net
	mm.Map = make(map[string]MarinerConfig)
	return mm
}

func Unpack(r io.ReadCloser, x interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
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

func InsecureClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
