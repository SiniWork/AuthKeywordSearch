/*
author: siyu
date: 2023/03/16
*/

package verification

import (
	"Collie/src/backtracking"
	"Collie/src/graph"
	"math/big"
)

type VOB struct {
	/*
	MP: the proof of searching MBT
	KeySets: the keyword sets
	SubG: the subgraph where performing the backwards search
	GH: the xor of the complements of the SubG
	NonExiP: the non exist proof of each keyword
	*/
	MP []byte
	KeySets map[string][]string
	SubG map[string]graph.Vertex
	GH []byte
	NonExiP map[string][]*big.Int
	RS      graph.ResTrees
}

type OnChainInfo struct {
	OnRH []byte
	OnGH []byte
	OnAc *big.Int
	OnBase *big.Int
	OnN *big.Int
}

func (vo *VOB) AuthenticationB(q []string, oCI OnChainInfo) bool{
	// step 1: verifying the processing of searching MBT
	reRH := vo.MP
	if len(vo.NonExiP) != 0 {
		for key, xd := range vo.NonExiP {
			if !VerifyNonExist(xd[0], xd[1], oCI.OnBase, oCI.OnN, oCI.OnAc, key) {
				return false
			}
		}
	}
	for _, key := range q {
		if _, exist := vo.KeySets[key]; exist {
			reRH = backtracking.XOR(reRH, backtracking.CompBucketHash(key, vo.KeySets[key]))
		}
	}
	if string(backtracking.Hash(reRH)) != string(oCI.OnRH) {
		return false
	} else if len(vo.KeySets) < len(q) {
		return true
	}

	// step 2: verifying the subgraph where performing the backwards search
	reGH := vo.GH
	for _, vtx := range vo.SubG {
		reGH = backtracking.XOR(reGH, vtx.CompuVtxHash())
	}
	if string(backtracking.Hash(reGH)) != string(oCI.OnGH) {
		return false
	}

	// step 3: backward search based on VO
	bookKeep := make(map[string]map[string]graph.Leaf) // key1: the id of the explored vertex; key2: the keyword
	priQue := &graph.PriorityQueue{}
	for kd, ids := range vo.KeySets {
		// Initializes the iterators
		iterHp := &graph.IterHeap{}
		for _, id := range ids {
			iter := graph.Iterator{id}
			iterHp.PushIter(iter)
			// Initializes the bookKeep
			if _, explored := bookKeep[id]; explored {
				bookKeep[id][kd] = graph.Leaf{Id: id, Dist: 0}
			} else {
				bookKeep[id] = make(map[string]graph.Leaf)
				bookKeep[id][kd] = graph.Leaf{Id: id, Dist: 0}
			}
		}
		priQue.PushItem(graph.Item{Keyword: kd, Priority: len(iterHp.Top()), IHp: *iterHp})
	}
	//backtracking.CheckStatus(bookKeep, *priQue)
	for !backtracking.Terminal(*priQue) {
		item := priQue.PopItem()
		iter := item.IHp.PopIter()
		for _, vtx := range vo.SubG[iter[len(iter)-1]].InEdges {
			if !backtracking.Exist(vtx, iter) && len(iter) < backtracking.HOP { // explored a node
				updIter := append(iter, vtx)
				item.IHp.PushIter(updIter)
				// update bookKeep...
				if _, explored := bookKeep[vtx]; explored {
					if _, keyReachable := bookKeep[vtx][item.Keyword]; keyReachable {
						if bookKeep[vtx][item.Keyword].Dist > (len(updIter)-1) {
							bookKeep[vtx][item.Keyword] = graph.Leaf{Id: updIter[0], Dist: len(updIter)-1}
						}
					} else {
						bookKeep[vtx][item.Keyword] = graph.Leaf{Id: updIter[0], Dist: len(updIter)-1}
					}
				} else {
					bookKeep[vtx] = make(map[string]graph.Leaf)
					bookKeep[vtx][item.Keyword] = graph.Leaf{Id: updIter[0], Dist: len(updIter)-1}
				}
			}
		}
		// update priorityQueue
		if len(item.IHp) > 0 {
			priQue.PushItem(item)
		}
	}
	//backtracking.CheckStatus(bookKeep, *priQue)
	res := graph.ResTrees{}
	for id, tree := range bookKeep {
		if len(tree) == len(q) {
			res.PushRes(graph.ResultTree{Root: id, Leaves: tree})
		}
	}
	if len(res) != len(vo.RS) {
		return false
	}
	return true
}

func VerifyNonExist(x, d, g, n, ac *big.Int, ele string) bool {
	// (ac^x)*(d)^b mod n = g
	b := backtracking.Str2BInt(ele)
	d.Exp(d, b, n)
	acc := new(big.Int).Exp(ac, x, n)
	flag := new(big.Int).Mul(acc, d)
	if flag.Mod(flag, n).Cmp(g) == 0 {
		return true
	}
	return false
}