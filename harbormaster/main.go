package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/marinatb/marina"
	"github.com/marinatb/marina/netdl"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func buildMap(net *netdl.Network) marina.MaterializationMap {

	mm := make(marina.MaterializationMap)

	//magic happens

	return mm

}

func materialize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	net := new(netdl.Network)
	err := marina.Unpack(r, net)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		d := marina.Diagnostic{"error", "malformed json"}
		w.Write(marina.Pack(d))
		return
	}

	mm := buildMap(net)

	xpdir := "/marina/xp/" + net.Name
	os.MkdirAll(xpdir, 0755)
	ioutil.WriteFile(xpdir+"/net.json", marina.Pack(*net), 0644)
	ioutil.WriteFile(xpdir+"/map.json", marina.Pack(mm), 0644)

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
