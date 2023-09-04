/*
author: siyu
date: 2023/03/08
*/

package backtracking

import (
	"Collie/src/graph"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"sort"
)

type MBT struct {
	/*
	Buckets: the key represents a keyword and the value stores the vertices containing that keyword
	BucketHashes: the hashes of each bucket
	RootHash: the root hash (saved on chain)
	 */
	Buckets map[string][]string
	BucketHashes map[string][]byte
	RootTag []byte
	RootHash []byte
}

func (t *MBT) BldMBT(file string) {
	t.Buckets, _ = graph.ReadKJsonFile(file)
	//t.LookInvert()
}

func (t *MBT) LookInvert() {
	for key, vtxs := range t.Buckets {
		fmt.Println("\"", key, "\"", ", number: ", len(vtxs))
	}
}

func (t *MBT) compRootHashes() {
	t.BucketHashes = make(map[string][]byte)
	for key, ids := range t.Buckets {
		t.BucketHashes[key] = CompBucketHash(key, ids)
		if len(t.RootTag) == 0 {
			t.RootTag = t.BucketHashes[key]
		} else {
			t.RootTag = XOR(t.RootTag, t.BucketHashes[key])
		}
	}
	t.RootHash = Hash(t.RootTag)
}

func (t *MBT) AuthSearch(keyList []string) (map[string][]string, []byte, []string) {
	/*
	Searching the MBT to obtain the target buckets as well as generate the VO
	 */
	results := make(map[string][]string)
	remainTag := t.RootTag
	nonExiK := []string{}
	for _, key := range keyList {
		if _, exist := t.Buckets[key]; exist {
			results[key] = t.Buckets[key]
			remainTag = XOR(remainTag, t.BucketHashes[key])
		} else {
			nonExiK = append(nonExiK, key)
		}
	}
	return results, remainTag, nonExiK
}

func CompBucketHash(key string, ids []string) []byte {
	sort.Strings(ids)
	raw := []interface{}{key, ids}
	rlp, err := rlp.EncodeToBytes(raw)
	if err != nil {
		fmt.Println(err)
	}
	return Hash(rlp)
}


