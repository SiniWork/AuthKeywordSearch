/*
author: siyu
date: 2023/03/08
*/

package backtracking

import (
	"Collie/src/graph"
	"crypto/md5"
	"fmt"
	"math/big"
)

const IDLEN =  10

type VO struct {
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

func (g *GraphBack) computeGHash() {
	g.VtxHashes = make(map[string][]byte)
	for id, vtx := range g.Vertices {
		g.VtxHashes[id] = vtx.CompuVtxHash()
		if len(g.GTag) == 0 {
			g.GTag = g.VtxHashes[id]
		} else {
			g.GTag = XOR(g.GTag, g.VtxHashes[id])
		}
	}
	g.GHash = Hash(g.GTag)
}

func (g *GraphBack) GenAuthInfo(invertF string) {
	/*
	Generate authentication information for the graph, including GHash and MBT
	*/
	g.computeGHash()
	g.Mbt.BldMBT(invertF)
	g.Mbt.compRootHashes()
	//keys := []string{}
	//for k, _ := range g.Mbt.Buckets {
	//	keys = append(keys, k)
	//}
	//g.Acc.Initial(keys)
}

func (g *GraphBack) AuthBackwards(q []string) VO {
	/*
	step1: authenticated search the MBT
	step2: authenticated backtracking search
	 */
	var vo VO
	// step1
	nonEK := []string{}
	vo.KeySets, vo.MP, nonEK = g.Mbt.AuthSearch(q)
	vo.GH = g.GTag
	vo.SubG = make(map[string]graph.Vertex)
	if len(nonEK) != 0 {
		vo.NonExiP = make(map[string][]*big.Int)
		for _, key := range nonEK {
			x, y := g.Acc.NonMemberProof(key)
			vo.NonExiP[key] = append(vo.NonExiP[key], x)
			vo.NonExiP[key] = append(vo.NonExiP[key], y)
		}
		return vo
	}
	for _, vertices := range vo.KeySets {
		for _, vtx := range vertices {
			if _, exist := vo.SubG[vtx]; !exist {
				vo.SubG[vtx] = g.Vertices[vtx]
				vo.GH = XOR(vo.GH, g.VtxHashes[vtx])
			}
		}
	}
	// step2
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
	//CheckStatus(bookKeep, *priQue)
	for !Terminal(*priQue) {
		item := priQue.PopItem()
		iter := item.IHp.PopIter()
		if len(iter) != 0 {
			for _, vtx := range g.Vertices[iter[len(iter)-1]].InEdges {
				// add vo
				if _, exist := vo.SubG[vtx]; !exist {
					vo.SubG[vtx] = g.Vertices[vtx]
					vo.GH = XOR(vo.GH, g.VtxHashes[vtx])
				}
				if !Exist(vtx, iter) && len(iter) < HOP { // explored a node
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
	}
	//CheckStatus(bookKeep, *priQue)
	res := graph.ResTrees{}
	for id, tree := range bookKeep {
		if len(tree) == len(q) {
			res.PushRes(graph.ResultTree{Root: id, Leaves: tree})
		}
	}
	vo.RS = res
	return vo
}

func (vo VO) Size() float32 {
	size := 0
	size += len(vo.MP)
	for key, ids := range vo.KeySets {
		size += len(key)
		size += len(ids) * IDLEN
	}
	for _, vtx := range vo.SubG {
		size += IDLEN
		size += len(vtx.InEdges) * IDLEN
		size += len(vtx.OutEdges) * IDLEN
		for _, key := range vtx.Keywords {
			size += len(key)
		}
	}
	size += len(vo.GH)
	fmt.Println(float32(size)/1024, "KB")
	return float32(size)
}

func XOR(str1, str2 []byte) []byte {
	/*
	Computing the XOR result of the given two byte array
	*/
	var res []byte
	if len(str1) != len(str2) {
		return res
	} else {
		for i:=0; i<len(str1); i++ {
			res = append(res, str1[i] ^ str2[i])
		}
	}
	return res
}

func Hash(raw []byte) []byte{
	h := md5.New()
	h.Write(raw)
	return h.Sum(nil)
}