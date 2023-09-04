/*
author: siyu
date: 2023/05/10
*/

package advanced

import (
	"Collie/src/graph"
	"fmt"
)

//type Leaf struct {
//	Id string
//	Dist int
//}
//
//type ResultTree struct {
//	Root string
//	Leaves map[string]Leaf
//	//Score int
//}
//
//type ResTrees []ResultTree

const IDLEN =  10

type VO struct {
	/*
	RS: the result trees
	*/
	NodeList []Node
	NodeListB map[Node]bool
	Buckets map[string]Bucket
	RS graph.ResTrees
	TargetK string
	UnRoot map[string][]string  // key: un-reached keyword, value: the list can not reach to the keyword of 'key'
}

func (g *GraphAdv) AuthFindResTrees(q []string) VO {
	/*
	Authenticated finding all result trees
	*/
	var vo VO
	// step 1: finding all target buckets
	vo.Buckets, vo.NodeList, vo.NodeListB = g.MPBT.AuthFindAllBuckets(q)

	// step 2: finding all result trees; finding top k result trees
	vo.RS, vo.UnRoot, vo.TargetK = g.FindAllResTrees(vo.Buckets)
	//fmt.Println(vo.RS)
	return vo
}

func (g *GraphAdv) FindAllResTrees(buckets map[string]Bucket) (graph.ResTrees, map[string][]string, string) {
	unRoot := make(map[string][]string)
	res := graph.ResTrees{}
	min := 10000000
	tarK := ""
	for k, b := range buckets {
		if len(b.Cont) < min {
			min = len(b.Cont)
			tarK = k
		}
	}
	for _, id := range buckets[tarK].SV {
		yes, proofK := IsRoot(id, buckets)
		if yes {
			reT := graph.ResultTree{}
			reT.Root = id
			reT.Leaves = make(map[string]graph.Leaf)
			for k, b := range buckets {
				reT.Leaves[k] = graph.Leaf{
					Id:   b.Cont[id].TargetN,
					Dist: b.Cont[id].Dis,
				}
			}
			res.PushRes(reT)
		} else {
			unRoot[proofK] = append(unRoot[proofK], id)
		}
	}
	return res, unRoot, tarK
}

func (vo VO) Size() float32 {
	size := 0
	NodeSize := 0
	hashL := 16
	intL := 4
	if len(vo.NodeList) != 0 {
		nodeMap := make(map[Node]bool)
		for _, node := range vo.NodeList {
			if hs, ok := node.(HashNode); ok {
				hashSize := len(hs.Hash())
				NodeSize = NodeSize + hashSize
			} else {
				if !nodeMap[node] {
					nodeMap[node] = true
					if leaf, ok := node.(*LeafNode); ok {
						buckSize := len(leaf.Value.Cont) * IDLEN
						buckSize += len(leaf.Value.Cont) * (intL + IDLEN)
						leafSize := len(leaf.Path) + buckSize
						NodeSize = NodeSize + leafSize
					} else if branch, ok := node.(*BranchNode); ok {
						branchSize := len(branch.Branches) * (intL+hashL)
						NodeSize = NodeSize + branchSize
					} else if ext, ok := node.(*ExtensionNode); ok {
						extSize := len(ext.Path)
						NodeSize = NodeSize + extSize
					}
				}
			}
		}
	}
	size += NodeSize
	for _, vtx := range vo.UnRoot {
		size += IDLEN * len(vtx)
	}
	fmt.Println(float32(size)/1024, "KB")
	return float32(size)/1024
}

func IsRoot(id string, bucks map[string]Bucket) (bool, string){
	yes := true
	for k, b := range bucks {
		_, exist := b.Cont[id]; if !exist {
			yes = false
			return yes, k
		}
	}
	return yes, ""
}

