package embedders

import (
	"github.com/marinatb/marina/netdl"
	"github.com/marinatb/marina/protocol"
)

func MclEmbed(net *netdl.Network) (error, *protocol.MaterializationEmbedding) {
	eb := protocol.NewMaterializationEmbedding(net)

	//magic happens

	return nil, eb

}
