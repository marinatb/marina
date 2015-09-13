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

func buildMap(net *netdl.Network) protocol.MaterializationMap {

	mm := make(protocol.MaterializationMap)

	//magic happens

	return mm

}

func materialize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	net := new(netdl.Network)
	err := protocol.Unpack(r, net)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		d := protocol.Diagnostic{"error", "malformed json"}
		w.Write(protocol.Pack(d))
		return
	}

	mm := buildMap(net)

	xpdir := "/marina/xp/" + net.Name
	os.MkdirAll(xpdir, 0755)
	ioutil.WriteFile(xpdir+"/net.json", protocol.Pack(*net), 0644)
	ioutil.WriteFile(xpdir+"/map.json", protocol.Pack(mm), 0644)

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
