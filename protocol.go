package marina

import (
	"bytes"
	"encoding/json"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)

type Diagnostic struct {
	Level, Message string
}

type MarinerConfig struct {
	NodeId uuid.UUID `json:"nodeId"`
	LXC    string    `json:"lxc"`
	OVS    string    `json:"ovs"`
	Click  string    `json:"click"`
}

type MaterializationMap map[uuid.UUID]MarinerConfig

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

func Pack(x interface{}) []byte {
	js, _ := json.Marshal(x)
	return js
}
