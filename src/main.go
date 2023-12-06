package main

import (
	"Collie/src/advanced"
	"fmt"
)

type Info struct {
	Time string
	Size float32
	Num int
}

func main() {
	dataset := "dblp"
	fmt.Println("----------------Loading Graph " + dataset + "----------------")
	advanced.StatisticVClasses("./data/" + dataset + "_scope.json")
	var g advanced.GraphAdv
	g.LoadGraph("./data/" + dataset + ".json", "./data/" + dataset + "_invert_key.json")
	g.Preprocessing(dataset + ".json")

	fmt.Println("----------------Building Trie----------------")
	buckets := advanced.LoadPreProcess("./data/" + dataset + "_precom.json")
	for k, b := range buckets {
		byteKey := []byte(k)
		g.MPBT.Insert(byteKey,b)
	}
	RD := g.MPBT.HashRoot()
	key := "RootDigest"
	eth.CommitEth(key, string(RD))
	
	fmt.Println("----------------Dealing with Query-----------")
	//q := []string{"a", "b", "c"}
	q := []string{"Keanu", "Matrix", "Thomas"}
	tStart := time.Now()
	voSP := g.AuthFindResTrees(q)
	tEnd := time.Since(tStart)
	fmt.Println(tEnd)
	size := voSP.Size()
	fmt.Println("the number of results: ", len(voSP.RS))
	
	info := Info{
		Time: tEnd.String(),
		Size: size,
		Num: len(voSP.RS),
	}
	fmt.Println(info)
	file, _ := os.OpenFile("dblp_info.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	defer file.Close()
	encoder := json.NewEncoder(file)
	err := encoder.Encode(info)
	if err != nil {
		fmt.Println("error")
	}
	
	voClient := verification.VOA{
		Buckets: voSP.Buckets,
		NodeList: voSP.NodeList,
		NodeListB: voSP.NodeListB,
		TarK: voSP.TargetK,
		UnRoot: voSP.UnRoot,
		RS: voSP.RS,
	}
	onChainInfo := verification.OnChainInfo{
		OnRH: []byte(eth.QueryEth("RootDigest").([]interface{})[1].(string)),
	}
	
	fmt.Println("----------------Verification-----------------")
	fmt.Println(voClient.AuthenticationA(q, onChainInfo))
}
