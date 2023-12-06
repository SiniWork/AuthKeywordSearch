package main

import (
	"Collie/src/backtracking"
	"Collie/src/verification"
	"fmt"
	"time"
)

func main() {
	fmt.Println("----------------Loading Graph----------------")
	var g backtracking.GraphBack
	dataset := "imdb"
	g.LoadGraph("./data/"+dataset+".json")
	g.GenAuthInfo("./data/"+dataset+"_invert_key.json")
	key1 := "RootDigest"
	eth.CommitEth(key1, string(g.Mbt.RootHash))
	key2 := "GraphDigest"
	eth.CommitEth(key1, string(g.GHash))


	//fmt.Println("Jude: ", len(g.Mbt.Buckets["Jude"]))
	//fmt.Println("Randy: ", len(g.Mbt.Buckets["Randy"]))
	//fmt.Println("Thomas: ", len(g.Mbt.Buckets["Thomas"]))
	//fmt.Println("Matrix: ", len(g.Mbt.Buckets["Matrix"]))
	//fmt.Println("Keanu: ", len(g.Mbt.Buckets["Keanu"]))

	// dealing with query
	fmt.Println("----------------Dealing with queries---------")
	//q := []string{"a", "b", "c"}
	q := []string{"Keanu", "Matrix", "Thomas"}
	voSP := backtracking.VO{}
	tStart := time.Now()
	voSP = g.AuthBackwards(q)
	tEnd := time.Since(tStart)
	fmt.Println(tEnd)
	voSP.Size()
	fmt.Println("the number of results: ", len(voSP.RS))
	voClient := verification.VOB{
		MP: voSP.MP,
		KeySets: voSP.KeySets,
		SubG: voSP.SubG,
		GH: voSP.GH,
		NonExiP: voSP.NonExiP,
		RS: voSP.RS,
	}
	onChainInfo := verification.OnChainInfo{
		OnRH: []byte(eth.QueryEth("RootDigest").([]interface{})[1].(string)),
		OnGH: []byte(eth.QueryEth("GraphDigest").([]interface{})[1].(string)),
		OnAc: g.Acc.Ac,
		OnBase: g.Acc.G,
		OnN: g.Acc.N,
	}

	// verification
	fmt.Println("----------------Verification-----------------")
	fmt.Println(voClient.AuthenticationB(q, onChainInfo))
}
