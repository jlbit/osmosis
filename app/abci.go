package app

import (
	"encoding/json"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
)

// Query implements the ABCI interface. It delegates to CommitMultiStore if it
// implements Queryable.
func (app *OsmosisApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if req.Path == "/cosmos.staking.v1beta1.Query/ValidatorDelegations" {
		return abci.ResponseQuery{
			Code:      1,
			Log:       "This query is too resource intensive. Please run your node",
			Codespace: "forbidden",
		}
	}

	if req.Path == "/simple-metric" || req.Path == "simple-metric" {
		metricRes := app.SimpleMetrics.CalcAllAverageResponses()
		jsonRes, err := json.Marshal(metricRes)
		if err != nil {
			panic(err)
		}
		return abci.ResponseQuery{
			Value: jsonRes,
		}
	}

	start := time.Now()
	res = app.BaseApp.Query(req)
	elapsed := time.Since(start)

	app.SimpleMetrics.Measure(req.Path, elapsed)

	return res
}
