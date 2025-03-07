package coin

import (
	"context"
	"fmt"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)

func TestGetBalance(t *testing.T) {
	balances, err := devCli.MgoXGetBalance(ctx, request.MgoXGetBalanceRequest{
		Owner:    "0x376e4c0da89168eb3cf0cda6910fff575d78b99296b3bb80c7a64a99449a1a25",
		CoinType: "0xff57fcbe56c70f7f32d9addb21a94f6277f528bfd07555f0eb4513ea50597449::meth::METH",
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	t.Log(balances)
}

func TestGetCoins(t *testing.T) {
	coins, err := devCli.MgoXGetCoins(ctx, request.MgoXGetCoinsRequest{
		Owner:    "0x376e4c0da89168eb3cf0cda6910fff575d78b99296b3bb80c7a64a99449a1a25",
		CoinType: "0x2::mgo::MGO",
		Limit:    10,
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	t.Log(coins)
}

func TestGetTotalSupply(t *testing.T) {
	coins, err := devCli.MgoXGetTotalSupply(ctx, request.MgoXGetTotalSupplyRequest{
		CoinType: "0x2::mgo::MGO",
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	t.Log(coins)
}
