/*
author: siyu
date: 2023/05/05
*/

package advanced

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"sort"
	"strconv"
)

type BranchNode struct {
	Branches map[byte]Node
	Value Bucket
	Hs []byte
}

func NewBranchNode() *BranchNode {
	return &BranchNode{
		Branches: make(map[byte]Node),
	}
}

func (b *BranchNode) SetBranch(bit byte, node Node) {
	b.Branches[bit] = node
}

func (b *BranchNode) GetBranch(bit byte) Node {
	return b.Branches[bit]
}

func (b *BranchNode) RemoveBranch(bit byte) {
	b.Branches[bit] = nil
}

func (b *BranchNode) SetValue(value Bucket) {
	b.Value = value
}

func (b *BranchNode) RemoveValue() {
	b.Value = Bucket{}
}

func (b BranchNode) Hash() []byte {
	return crypto.Keccak256(b.Serialize())
}

func (b BranchNode) Serialize() []byte {
	raw := make([]interface{}, BranchSize+1)
	branchLst := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	for i, c := range branchLst {
		if b.Branches[c] == nil {
			raw[i] = " "
		} else {
			node := b.Branches[c]
			raw[i] = node.Hash()
		}
	}
	var valueStr []string
	var keys []string
	for k, _ := range b.Value.Cont {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		valueStr = append(valueStr, b.Value.Cont[k].TargetN)
		valueStr = append(valueStr, strconv.Itoa( b.Value.Cont[k].Dis))
	}
	for _, v := range b.Value.SV {
		valueStr = append(valueStr, v)
	}
	raw[BranchSize] = valueStr
	rlp, _ := rlp.EncodeToBytes(raw)
	return rlp
}