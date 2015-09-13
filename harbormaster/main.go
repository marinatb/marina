package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/marinatb/marina"
	"github.com/marinatb/marina/netdl"
	"github.com/marinatb/marina/protocol"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func mclMap(net *netdl.Network) *protocol.MaterializationMap {
	mm := new(protocol.MaterializationMap)

	//magic happens

	return mm
}

func buildMap(net *netdl.Network, mapper string) *protocol.MaterializationMap {
	switch mapper {
	case "mcl":
		return mclMap(net)
	default:
		log.Printf("unkown mapper '%s", mapper)
		return nil
	}
}

func materialize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	rq := new(protocol.NetworkMaterializationRequest)
	err := protocol.Unpack(r, rq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		d := protocol.Diagnostic{"error", "malformed json"}
		w.Write(protocol.PackWire(d))
		return
	}

	mm := buildMap(&rq.Net, rq.Mapper)

	xpdir := "/marina/xp/" + rq.Net.Name
	os.MkdirAll(xpdir, 0755)
	ioutil.WriteFile(xpdir+"/net.json", protocol.PackLegible(rq.Net), 0644)
	ioutil.WriteFile(xpdir+"/map.json", protocol.PackLegible(mm), 0644)

}

func main() {

	log.Printf("harbormaster v%d.%d\n", marina.MajorVersion, marina.MinorVersion)

	router := httprouter.New()
	router.POST("/materialize", materialize)

	log.Printf("listening on https://::0:8080/")
	log.Fatal(
		http.ListenAndServeTLS(":8080",
			"/marina/keys/cert.pem",
			"/marina/keys/key.pem",
			router))

}
