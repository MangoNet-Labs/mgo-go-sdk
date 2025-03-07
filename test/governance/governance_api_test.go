package governance

import (
	"context"
	"testing"

	"github.com/mangonet-labs/mgo-go-sdk/client"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/utils"
)

var ctx = context.Background()
var devCli = client.NewMgoClient(config.RpcMgoTestnetEndpoint)

func TestGetReferenceGasPrice(t *testing.T) {
	price, err := devCli.MgoXGetReferenceGasPrice(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(price)
}

func TestGetStakes(t *testing.T) {
	stakes, err := devCli.MgoXGetStakes(ctx, request.MgoXGetStakesRequest{
		Owner: "0x6d5ae691047b8e55cb3fc84da59651c5bae57d2970087038c196ed501e00697b",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(stakes))
	if len(stakes) > 0 {
		t.Log(stakes[0].ValidatorAddress)
		t.Log(stakes[0].StakingPool)
		t.Log(len(stakes[0].Stakes))
		if len(stakes[0].Stakes) > 0 {
			t.Log(stakes[0].Stakes[0])
			t.Logf("stakeId:%v", stakes[0].Stakes[0].StakedMgoId)
		}
	}
}

func TestGetStakesByIds(t *testing.T) {
	stakes, err := devCli.MgoXGetStakesByIds(ctx, request.MgoXGetStakesByIdsRequest{
		StakedMgoIds: []string{"0x70a3040054dede54d0e99be74ca80e22be5cd5710c57a725d55c2c7640b0028b"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(stakes))
	if len(stakes) > 0 {
		t.Log(stakes[0].ValidatorAddress)
		t.Log(stakes[0].StakingPool)
		t.Log(len(stakes[0].Stakes))
		if len(stakes[0].Stakes) > 0 {
			t.Log(stakes[0].Stakes[0])
			t.Logf("stakeId:%v", stakes[0].Stakes[0].StakedMgoId)
		}
	}
}

func TestGetCommitteeInfo(t *testing.T) {
	info, err := devCli.MgoXGetCommitteeInfo(ctx, request.MgoXGetCommitteeInfoRequest{
		Epoch: "50",
	})
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(info)
}

func TestGetLatestMgoSystemState(t *testing.T) {
	state, err := devCli.MgoXGetLatestMgoSystemState(ctx)
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(state)
}

func TestGetValidatorsApy(t *testing.T) {
	apy, err := devCli.MgoXGetValidatorsApy(ctx)
	if err != nil {
		t.Fatal(err)
	}
	utils.JsonPrint(apy)
}
