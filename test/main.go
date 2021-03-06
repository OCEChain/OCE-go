/* package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/OCEChain/oce-go/models/auth"
	"github.com/OCEChain/oce-go/models/bank"
	"github.com/OCEChain/oce-go/models/types"

	"github.com/OCEChain/oce-go/models/auth/txbuilder"
	"github.com/OCEChain/oce-go/models/bank/client"
	"github.com/OCEChain/oce-go/models/crypto/hd"
	bip39 "github.com/cosmos/go-bip39"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tmlibs/bech32"
)

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterConcrete(bank.MsgSend{}, "cosmos-sdk/Send", nil)
	cdc.RegisterConcrete(auth.StdTx{}, "auth/StdTx", nil)
}

func main() {

	// 生成公私钥
	// 发送交易
	// 查询地址余额
	// 查询交易(txhash、height)

	//genKey()
	//sendTX()
	//getAccount()
	//getTX()
}

func genKey() {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		log.Fatalln(err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatalln(err)
	}
	seed := bip39.NewSeed(mnemonic, "")
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, "44'/118'/0'/0/0")
	if err != nil {
		log.Fatalln(err)
	}

	prik := secp256k1.PrivKeySecp256k1(derivedPriv)
	pubk := prik.PubKey()
	Addr, err := bech32.ConvertAndEncode("adr", pubk.Address().Bytes())
	if err != nil {
		log.Fatalln(err)
	}
	PubKey, _ := bech32.ConvertAndEncode("pub", pubk.Bytes())
	fmt.Println("Address:   " + Addr)
	fmt.Println("PublicKey: " + PubKey)
	fmt.Println("Mnemonic:  " + mnemonic)
}

func getAccount() {
	_, bz, err := bech32.DecodeAndConvert("adr1ttsph4qv93hllu8spl026s0rfmwhfl9d6fenyw")
	hexPubKey := append([]byte("account:"), bz...)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("Hex PublicKey:" + hex.EncodeToString(hexPubKey))

	url := `http://120.132.120.245/abci_query?path="/store/acc/key"&data=0x` + hex.EncodeToString(hexPubKey)
	res := httpGet(url)

	accRes := AccountResponse{}
	err = json.Unmarshal(res, &accRes)
	if err != nil {
		log.Fatalln(err)
	}

	br, err := base64.StdEncoding.DecodeString(accRes.Result.Response.Value)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(br))
}

func getTX() {
	url := `http://120.132.120.245/tx?hash=0xBB83B9A3A0D41CF0FAB1933F08CD6FD7000F28CB04AAEAD30FDF70BE466D3714`
	res := httpGet(url)

	tranRes := TranResponse{}
	err := json.Unmarshal(res, &tranRes)
	if err != nil {
		log.Fatalln(err)
	}

	br, err := base64.StdEncoding.DecodeString(tranRes.Result.Tx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(br))
}

func sendTX() {
	//fromAdr := "adr12fxqmhv9steldtqykkjm2emql8eqfvw6am76xj"
	fromAdr := "adr1ttsph4qv93hllu8spl026s0rfmwhfl9d6fenyw"
	toAdr := "adr1yrd22rg0hq3wkj4jwv0s8z8xp9fpnah8dd5u59"
	from, err := types.AccAddressFromBech32(fromAdr)
	if err != nil {
		log.Fatalln(err)
	}
	to, err := types.AccAddressFromBech32(toAdr)
	if err != nil {
		log.Fatalln(err)
	}

	coins, err := types.ParseCoins("8888coin1")
	if err != nil {
		log.Fatalln(err)
	}
	msg := client.CreateMsg(from, to, coins)

	tb := txbuilder.StdSignMsg{
		ChainID:  "oce",
		Sequence: 0,
		Memo:     "",
		Msgs:     []types.Msg{msg},
		Fee:      auth.NewStdFee(200000, types.Coin{}),
	}
	sign, err := buildAndSign(tb)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(hex.EncodeToString(sign))

	//Commit (Waiting for a new block)
	//url := "http://120.132.120.245/broadcast_tx_commit?tx=0x" + hex.EncodeToString(sign)

	//Propose (Waiting for the proposal result)
	url := "http://120.132.120.245/broadcast_tx_sync?tx=0x" + hex.EncodeToString(sign)
	fmt.Println(string(httpGet(url)))
}

func buildAndSign(msg txbuilder.StdSignMsg) ([]byte, error) {
	//mnemonic := "bounce prevent cross remind lunch pitch project dragon firm stove labor bicycle phrase giggle cliff huge betray mask ecology gloom access alarm yellow tuna"
	mnemonic := "unfair subway explain reward shrug cement dial junk twin vital badge sing lift chair cage interest rack fault feature original acoustic vote sheriff car"
	seed := bip39.NewSeed(mnemonic, "")
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, "44'/118'/0'/0/0")
	if err != nil {
		log.Fatalln(err)
	}

	priv := secp256k1.PrivKeySecp256k1(derivedPriv)

	sigBytes, err := priv.Sign(msg.Bytes())
	if err != nil {
		log.Fatalln(err)
	}
	pubkey := priv.PubKey()

	sig := auth.StdSignature{
		Sequence:  msg.Sequence,
		PubKey:    pubkey,
		Signature: sigBytes,
	}
	return cdc.MarshalJSON(auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo))
}

func httpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

type TranResponse struct {
	Result TxResult `json:"result"`
}

type TxResult struct {
	Tx string `json:"tx"`
}

type AccountResponse struct {
	Result Result
}

type Result struct {
	Response Response
}

type Response struct {
	Value string `json:"value"`
}

type TxResponse struct {
	Value Value
}

type Value struct {
	Msg Msg
}

type Msg struct {
	Value MsgValue
}

type MsgValue struct {
	Inputs Inputs
}

type Inputs struct {
	Address string `json:"address"`
	Coins   Coins
}

type Outputs struct {
	Address string `json:"address"`
	Coins   Coins
}

type Coins struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
*/