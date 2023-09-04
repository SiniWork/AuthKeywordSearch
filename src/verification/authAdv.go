/*
author: siyu
date: 2023/05/10
*/

package verification

import (
	"Collie/src/advanced"
	"Collie/src/graph"
	"fmt"
)

type VOA struct {
	/*
	NodeList: the proof of searching MPBT
	*/
	NodeList []advanced.Node
	NodeListB map[advanced.Node]bool
	Buckets map[string]advanced.Bucket
	UnRoot map[string][]string
	TarK string
	RS graph.ResTrees
}

func (vo *VOA) AuthenticationA(q []string, oCI OnChainInfo) bool {
	// step 1: verifying the processing of searching MPBT
	for _, k := range q {
		flag := reSearch([]byte(advanced.TransK2Hex(k)), vo.NodeList, vo.NodeListB)
		if !flag {
			return false
		}
	}
	reRH := vo.NodeList[0].Hash()
	if string(reRH) != string(oCI.OnRH) {
		return false
	}
	// step 2: verifying the processing of searching result trees
	for _, t := range vo.RS { // verify the root nodes come from original graph
		_, exist := vo.Buckets[vo.TarK].Cont[t.Root]; if !exist {
			return false
		}
	}
	unRNum := 0
	for k, vL := range vo.UnRoot { // verify the un-root nodes come from original graph and are indeed un-root
		unRNum += len(vL)
		for _, v := range vL {
			_, exist1 := vo.Buckets[vo.TarK].Cont[v]  // v should exist in the target keyword bucket
			_, exist2 := vo.Buckets[k].Cont[v]        // v should not exist in other buckets
			if !exist1 || exist2 {
				return false
			}
		}
	}
	if (len(vo.RS) + unRNum) != len(vo.Buckets[vo.TarK].SV) { // verify no valid root is maliciously discarded
		return false
	}
	return true
}

func reSearch(key []byte, nodeL []advanced.Node, nodeLB map[advanced.Node]bool) bool {
	if len(key) == 0 {
		return false
	}
	if len(nodeL) == 0 {
		return false
	}
	node := nodeL[0]
	for {
		if advanced.IsEmptyNode(node) {
			return false
		}

		if leaf, ok := node.(*advanced.LeafNode); ok {
			matched := advanced.PrefixMatchedLen(leaf.Path, key)
			if matched == len(key) {
				return true
			} else {
				return false
			}
		}

		if branch, ok := node.(*advanced.BranchNode); ok {
			if len(key) == 0 {
				return true
			}
			b, remaining := key[0], key[1:]
			key = remaining
			node = branch.GetBranch(b)
			continue
		}

		if ext, ok := node.(*advanced.ExtensionNode); ok {
			matched := advanced.PrefixMatchedLen(ext.Path, key)
			if matched < len(ext.Path) && matched < len(key){
				return false
			}
			key = key[matched:]
			node = ext.Next
			continue
		}
		fmt.Println("keyword not exist")
		return false
	}
}