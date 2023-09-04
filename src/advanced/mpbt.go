/*
author: siyu
date: 2023/05/08
*/

package advanced

import (
	"errors"
	"fmt"
)

const BranchSize = 16

type Node interface {
	Hash() []byte
	Serialize() []byte
}

type HashNode struct {
	Hs []byte
}

func (n HashNode) Hash() []byte {
	return n.Hs
}

func (n HashNode) Serialize() []byte {
	return n.Hs
}

func NewHashNode(hash []byte) HashNode {
	hashNode := HashNode{Hs: hash}
	return hashNode
}

type Trie struct {
	root Node
}

func (t *Trie) GetRoot() Node {
	return t.root
}

func NewTrie() *Trie {
	return &Trie{}
}

func IsEmptyNode(node Node) bool {
	return node == nil
}

func (t *Trie) Insert(key []byte, bucket Bucket) error {
	/*
		Inserting (key, value) into trie
		key: the key to be inserted
		bucket: the bucket corresponding to the key
	*/

	if len(key) == 0 {
		return errors.New("the key is empty")
	}
	node := &t.root
	var pre = node
	var recordB byte
	for {
		if IsEmptyNode(*node) {
			leaf := NewLeafNode(key, bucket)
			*node = leaf
			return nil
		}
		// leaf node case
		if leaf, ok := (*node).(*LeafNode); ok {
			matched := PrefixMatchedLen(leaf.Path, key)
			// first case: full matched
			if matched == len(key) && matched == len(leaf.Path) {
				// It's wrong, because each key is different
				return errors.New("the key is already exist")
			}
			// second case: no matched
			branch := NewBranchNode()
			if matched == 0 {
				if preBranch, yes := (*pre).(*BranchNode); yes {
					preBranch.SetBranch(recordB, branch)
				}
				*node = branch
				if len(key) == 0 {
					branch.SetValue(bucket)
					oldLeaf := NewLeafNode(leaf.Path[1:], leaf.Value)
					branch.SetBranch(leaf.Path[0], oldLeaf)
					return nil
				}
				if len(leaf.Path) == 0 {
					branch.SetValue(leaf.Value)
					newLeaf := NewLeafNode(key[1:], bucket)
					branch.SetBranch(key[0], newLeaf)
					return nil
				}
				oldLeaf := NewLeafNode(leaf.Path[1:], leaf.Value)
				branch.SetBranch(leaf.Path[0],oldLeaf)
				newLeaf := NewLeafNode(key[1:], bucket)
				branch.SetBranch(key[0], newLeaf)
				return nil
			}
			// third case: part matched
			ext := NewExtensionNode(leaf.Path[:matched], branch)
			*node = ext
			if preBranch, yes := (*pre).(*BranchNode); yes {
				preBranch.SetBranch(recordB, ext)
			}
			if matched == len(leaf.Path) {
				branch.SetValue(leaf.Value)
				branchKey, leafKey := key[matched], key[matched+1:]
				newLeaf := NewLeafNode(leafKey, bucket)
				branch.SetBranch(branchKey, newLeaf)
			} else if matched == len(key) {
				branch.SetValue(bucket)
				oldBranchKey, oldLeafKey := leaf.Path[matched], leaf.Path[matched+1:]
				oldLeaf := NewLeafNode(oldLeafKey, leaf.Value)
				branch.SetBranch(oldBranchKey, oldLeaf)
			} else {
				oldBranchKey, oldLeafKey := leaf.Path[matched], leaf.Path[matched+1:]
				oldLeaf := NewLeafNode(oldLeafKey, leaf.Value)
				branch.SetBranch(oldBranchKey, oldLeaf)
				branchKey, leafKey := key[matched], key[matched+1:]
				newLeaf := NewLeafNode(leafKey, bucket)
				branch.SetBranch(branchKey, newLeaf)
			}
			return nil
		}
		// branch node case
		if branch, ok := (*node).(*BranchNode); ok {
			if len(key) == 0 {
				branch.SetValue(bucket)
				return nil
			}
			pre = node
			recordB = key[0]
			b, remaining := key[0], key[1:]
			key = remaining
			tmp := branch.GetBranch(b)
			if tmp == nil {
				leaf := NewLeafNode(key, bucket)
				branch.SetBranch(b, leaf)
				return nil
			} else {
				node = &tmp
				continue
			}
		}
		// extension node case
		if ext, ok := (*node).(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Path, key)
			// first case: full matched
			if  matched == len(ext.Path) {
				key = key[matched:]
				node = &ext.Next
				continue
			}
			// second case: no matched
			branch := NewBranchNode()
			if matched == 0 {
				if preBranch, ok := (*pre).(*BranchNode); ok {
					preBranch.SetBranch(recordB, branch)
				}
				extBranchKey, extRemainingKey := ext.Path[0], ext.Path[1:]
				if len(extRemainingKey) == 0 {
					branch.SetBranch(extBranchKey, ext.Next)
				} else {
					newExt := NewExtensionNode(extRemainingKey, ext.Next)
					branch.SetBranch(extBranchKey, newExt)
				}
				if len(key) == 0 {
					branch.SetValue(bucket)
					*node = branch
				} else {
					leafBranchKey, leafRemainingKey := key[0], key[1:]
					newLeaf := NewLeafNode(leafRemainingKey, bucket)
					branch.SetBranch(leafBranchKey, newLeaf)
					*node = branch
				}
				return nil
			}
			// third case: part matched
			commonKey, branchKey, extRemainingKey := ext.Path[:matched], ext.Path[matched], ext.Path[matched+1:]
			oldExt := NewExtensionNode(commonKey, branch)
			if preBranch, ok := (*pre).(*BranchNode); ok {
				preBranch.SetBranch(recordB, oldExt)
			}
			if len(extRemainingKey) == 0 {
				branch.SetBranch(branchKey, ext.Next)
			} else {
				newExt := NewExtensionNode(extRemainingKey, ext.Next)
				branch.SetBranch(branchKey, newExt)
			}
			if len(commonKey) == len(key) {
				branch.SetValue(bucket)
			} else {
				leafBranchKey, leafRemainingKey := key[matched], key[matched+1:]
				newLeaf := NewLeafNode(leafRemainingKey, bucket)
				branch.SetBranch(leafBranchKey, newLeaf)
			}
			*node = oldExt
			return nil
		}
		return errors.New("unknown type")
	}
}

func (t *Trie) Get(key []byte) (bool, Bucket) {
	/*
		Get the target node depends on the given key
	*/
	node := t.root
	for {
		if IsEmptyNode(node) {
			return false, Bucket{}
		}

		if leaf, ok := node.(*LeafNode); ok {
			fmt.Println("leaf node") // for test
			matched := PrefixMatchedLen(leaf.Path, key)
			if matched != len(leaf.Path) || matched != len(key) {
				return false, Bucket{}
			}
			return true, leaf.Value
		}

		if branch, ok := node.(*BranchNode); ok {
			fmt.Println("branch node") // for test
			if len(key) == 0 {
				return true, branch.Value
			}
			b, remaining := key[0], key[1:]
			key = remaining
			node = branch.GetBranch(b)
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
			fmt.Println("extension node") // for test
			matched := PrefixMatchedLen(ext.Path, key)
			if matched < len(ext.Path) {
				return false, Bucket{}
			}
			key = key[matched:]
			node = ext.Next
			continue
		}
		return false, Bucket{}
	}
}

func (t *Trie) HashRoot() []byte {
	/*
		computing the root hash
	*/
	if t.root == nil {
		return nil
	}
	hashed := hash(&t.root)
	return hashed
}

func (t *Trie) Print() {
	if t.root == nil {
		return
	}
	PrintNode(t.root)
	return
}

func (t *Trie) AuthFindAllBuckets(keys []string) (map[string]Bucket, []Node, map[Node]bool) {
	/*
		obtaining all buckets for the given keywords
	*/
	//nodeListB := make(map[Node]bool)
	var nodeList []Node
	resBucks := make(map[string]Bucket)
	nodeListB := make(map[Node]bool)
	for _, k := range keys {
		flag, buck, nodes := t.AuthSearch([]byte(TransK2Hex(k)))
		if flag {
			resBucks[k] = buck
			for _, n := range nodes {
				if !nodeListB[n] {
					nodeListB[n] = true
					nodeList = append(nodeList, n)
				}
			}
		}
	}

	// Building the bound
	for _, node := range nodeList {
		switch node.(type) {
		case *LeafNode:
			continue
		case *BranchNode:
			branch, _ := (node).(*BranchNode)
			for i, child := range branch.Branches {
				if IsEmptyNode(child) || nodeListB[child]{
					continue
				} else {
					var hashNode HashNode
					if leaf, ok := child.(*LeafNode); ok {
						hashNode = NewHashNode(leaf.Hs)
					} else if bran, y := child.(*BranchNode); y {
						hashNode = NewHashNode(bran.Hs)
					} else {
						hashNode = NewHashNode(child.(*ExtensionNode).Hs)
					}
					branch.SetBranch(i, hashNode)
					nodeList = append(nodeList, hashNode)
				}
			}
		case *ExtensionNode:
			ext, _ := (node).(*ExtensionNode)
			if !nodeListB[ext.Next] {
				hashNode := NewHashNode(ext.Next.(*BranchNode).Hs)
				ext.SetNext(hashNode)
				nodeList = append(nodeList, hashNode)
			}
		}
	}
	return resBucks, nodeList, nodeListB
}

func (t *Trie) AuthSearch(key []byte) (bool, Bucket, []Node) {
	/*
		obtaining the target bucket of the given key and its merkle proof
	*/
	var vNodes []Node
	if len(key) == 0 {
		return false, Bucket{}, vNodes
	}
	node := t.root
	for {
		if IsEmptyNode(node) {
			return false, Bucket{}, vNodes
		}

		if leaf, ok := node.(*LeafNode); ok {
			//fmt.Println("leaf node")
			vNodes = append(vNodes, leaf)
			matched := PrefixMatchedLen(leaf.Path, key)
			if matched == len(key) {
				return true, leaf.Value, vNodes
			} else {
				return false, Bucket{}, vNodes
			}
		}

		if branch, ok := node.(*BranchNode); ok {
			//fmt.Println("branch node")
			vNodes = append(vNodes, branch)
			if len(key) == 0 {
				return true, branch.Value, vNodes
			}
			b, remaining := key[0], key[1:]
			key = remaining
			node = branch.GetBranch(b)
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
			//fmt.Println("extension node")
			vNodes = append(vNodes, ext)
			matched := PrefixMatchedLen(ext.Path, key)
			if matched < len(ext.Path) && matched < len(key){
				return false, Bucket{}, vNodes
			}
			key = key[matched:]
			node = ext.Next
			continue
		}
		return false, Bucket{}, vNodes
	}
}

func PrintNode(node Node) {
	switch (node).(type) {
	case *LeafNode:
		leaf, _ := (node).(*LeafNode)
		fmt.Println("LeafNode hash: ", leaf.Hs)
		return
	case *ExtensionNode:
		ext, _ := (node).(*ExtensionNode)
		fmt.Println("ExtensionNode hash: ", ext.Hs)
		PrintNode(ext.Next)
		return
	case *BranchNode:
		branch, _ := (node).(*BranchNode)
		fmt.Println("BranchNode hash: ", branch.Hs)
		branchLst := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
		for _, i := range branchLst {
			if child := branch.Branches[i]; child != nil {
				PrintNode(child)
			}
		}
		return
	}
	return
}

func hash(node *Node) []byte {
	/*
		computing root hash of the subtree corresponding to the given node
	*/
	switch (*node).(type) {
	case *LeafNode:
		leaf, _ := (*node).(*LeafNode)
		leaf.Hs = leaf.Hash()
		return leaf.Hs
	case *ExtensionNode:
		ext, _ := (*node).(*ExtensionNode)
		hash(&ext.Next)
		ext.Hs = ext.Hash()
		return ext.Hs
	case *BranchNode:
		branch, _ := (*node).(*BranchNode)
		branchLst := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
		for _, i := range branchLst {
			if child := branch.Branches[i]; child != nil {
				hash(&child)
			}
		}
		branch.Hs = branch.Hash()
		return branch.Hs
	}
	return nil
}

func PrefixMatchedLen(node1, node2 []byte) int {
	matched := 0
	for i := 0; i < len(node1) && i < len(node2); i++ {
		n1, n2 := node1[i], node2[i]
		if n1 == n2 {
			matched++
		} else {
			break
		}
	}
	return matched
}






