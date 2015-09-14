package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/marinatb/marina"
	"github.com/marinatb/marina/embedders"
	"github.com/marinatb/marina/netdl"
	"github.com/marinatb/marina/protocol"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func embed(net *netdl.Network, mapper string) (
	error, *protocol.MaterializationEmbedding) {
	switch mapper {
	case "mcl":
		return embedders.DefaultEmbed(net)
	default:
		err := fmt.Errorf("unkown mapper '%s'", mapper)
		log.Println(err)
		return err, nil
	}
}

func materialize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	log.Println("[materialize]")

	rq := new(protocol.NetworkMaterializationRequest)
	err := protocol.Unpack(r.Body, rq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		d := protocol.Diagnostic{"error", "malformed json"}
		w.Write(protocol.PackWire(d))
		return
	}

	err, em := embed(&rq.Net, rq.Mapper)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		d := protocol.Diagnostic{"error", fmt.Sprintf("materialization: %s", err)}
		w.Write(protocol.PackWire(d))
		return
	}

	xpdir := "/marina/xp/" + rq.Net.Name
	os.MkdirAll(xpdir, 0755)
	ioutil.WriteFile(xpdir+"/net.json", protocol.PackLegible(rq.Net), 0644)
	ioutil.WriteFile(xpdir+"/map.json", protocol.PackLegible(em), 0644)

}

func main() {

	log.Printf("harbormaster v%d.%d\n", marina.MajorVersion, marina.MinorVersion)

	router := httprouter.New()
	router.POST("/materialize", materialize)

	log.Printf("listening on https://::0:4676/")
	log.Fatal(
		http.ListenAndServeTLS(":4676",
			"/marina/keys/cert.pem",
			"/marina/keys/key.pem",
			router))

}
