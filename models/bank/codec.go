package bank

import (
	"github.com/OCEChain/oce-go/models/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/Send", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
