/*
author: siyu
date: 2023/05/08
*/

package advanced

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type ExtensionNode struct {
	Path []byte
	Next Node
	Hs []byte
}

func NewExtensionNode(path []byte, next Node) *ExtensionNode {
	return &ExtensionNode{
		Path: path,
		Next: next,
	}
}

func (e *ExtensionNode) SetNext(node Node) {
	e.Next = node
}

func (e ExtensionNode) Hash() []byte {
	return crypto.Keccak256(e.Serialize())
}

func (e ExtensionNode) Serialize() []byte {
	raw := make([]interface{}, 2)
	raw[0] = e.Path
	raw[1] = e.Next.Hash()

	rlp, _ := rlp.EncodeToBytes(raw)
	return rlp
}