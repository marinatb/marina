package testsite

import (
	"github.com/marinatb/marina/protocol"
	"io/ioutil"
	"testing"
)

func TestLocal(t *testing.T) {

	net := Build()
	t.Log(net)

	js := protocol.PackLegible(net)
	ioutil.WriteFile("/marina/site/testsite.json", js, 0644)
}
