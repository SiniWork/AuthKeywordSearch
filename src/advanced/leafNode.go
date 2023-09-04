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

type LeafNode struct {
	Path []byte
	Value Bucket
	Hs []byte
}

func NewLeafNode(key []byte, value Bucket) *LeafNode {
	return &LeafNode{
		Path:  key,
		Value: value,
	}
}

func (l LeafNode) Hash() []byte {
	//h := md5.New()
	//h.Write(l.Serialize())
	//return h.Sum(nil)
	return crypto.Keccak256(l.Serialize())
}

func (l LeafNode) Serialize() []byte {
	path := l.Path
	var valueStr []string
	var keys []string
	for k, _ := range l.Value.Cont {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		valueStr = append(valueStr, l.Value.Cont[k].TargetN)
		valueStr = append(valueStr, strconv.Itoa( l.Value.Cont[k].Dis))
	}
	for _, v := range l.Value.SV {
		valueStr = append(valueStr, v)
	}

	raw := []interface{}{path, valueStr}
	rlp, _ := rlp.EncodeToBytes(raw)
	return rlp
}