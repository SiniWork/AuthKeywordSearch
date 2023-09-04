/*
author: siyu
date: 2023/03/07
 */

package backtracking

import (
	"Collie/src/graph"
	"fmt"
)


const HOP = 4

type GraphBack struct {
	/*
		VtxHashes: the dict to store the hash value of each vertex
		GTag: the XOR of all hashes of the vertices
		GHash: the hash value of the graph (saved on chain)
		Mbt: the merkle bucket tree that based on the keywords of the graph
	*/
	graph.Graph
	VtxHashes map[string][]byte
	GTag []byte
	GHash    []byte
	Mbt      MBT
	Acc      Accumulator
}

func (g *GraphBack) LoadGraph(graphF string) {
	g.Vertices = graph.ReadGJsonFile(graphF)
}

func (g *GraphBack) ObtainKeySets(keywords []string) map[string][]string {
	/*
	Obtain the keywords vertex sets for each keyword
	 */
	keywordSets := make(map[string][]string)
	for _, kd := range keywords {
		keywordSets[kd] = g.Mbt.Buckets[kd]
	}
	return keywordSets
}

func (g *GraphBack) Backwards(keywords []string) graph.ResTrees {
	/*
	The Backtracking Search Algorithm
	 */
	KS := g.ObtainKeySets(keywords)
	// step 1: initialization
	bookKeep := make(map[string]map[string]graph.Leaf) // key1: the id of the explored vertex; key2: the keyword
	priQue := &graph.PriorityQueue{}
	for kd, ids := range KS {
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
	// check the initialization
	fmt.Println("-----------------initial:")
	CheckStatus(bookKeep, *priQue)

	// step2: backwards search
	for !Terminal(*priQue) {
		item := priQue.PopItem()
		iter := item.IHp.PopIter()
		for _, vtx := range g.Vertices[iter[len(iter)-1]].InEdges {
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
	fmt.Println("-----------------result:")
	res := graph.ResTrees{}
	for id, tree := range bookKeep {
		if len(tree) == len(keywords) {
			res.PushRes(graph.ResultTree{Root: id, Leaves: tree})
		}
	}
	fmt.Println(res)
	//checkStatus(bookKeep, *priQue)

	return res
}

func CheckStatus(bk map[string]map[string]graph.Leaf, pq graph.PriorityQueue) {
	/*
	Check the status of the "Backwards" processing
	 */
	fmt.Println("-----------bookKeep:")
	for id, book := range bk {
		fmt.Println(id, ": ", book)
	}
	fmt.Println("----------iterators:")
	for pq.Len() > 0 {
		item := pq.PopItem()
		fmt.Println(item.Keyword, ": ")
		for item.IHp.Len() > 0 {
			iter := item.IHp.PopIter()
			fmt.Println("(", iter[len(iter)-1], ",", iter, ",", len(iter)-1, ")")
		}
	}
}

func Terminal(p graph.PriorityQueue) bool {
	/*
	Determines whether the "Backwards" should terminal
	 */
	if p.Len() > 0 {
		return false
	} else {
		return true
	}
}

func Exist(target string, lis []string) bool {
	/*
	Determines whether one str existed in a string list
	 */
	for _, each := range lis {
		if target == each {
			return true
		}
	}
	return false
}

