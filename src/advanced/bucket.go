/*
author: siyu
date: 2023/05/04
*/

package advanced

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Path struct {
	TargetN string
	Dis int
}

type Bucket struct {
	Cont map[string]Path  // v_id: path
	SV []string  // sorted v_id according to distance
}

type VerPair struct {
	ID string
	Dis int
}

type VLst []VerPair

func LoadPreProcess(filePath string) map[string]Bucket {
	readCont := make(map[string]Bucket)
	buckets := make(map[string]Bucket)
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	decoder.Decode(&readCont)
	for k, b := range readCont {
		b.SV = SortV(b.Cont)
		codK := TransK2Hex(k)
		buckets[codK] = b
	}
	//fmt.Println(buckets)
	return buckets
}

func (p VLst) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p VLst) Len() int           { return len(p) }
func (p VLst) Less(i, j int) bool { return p[i].Dis < p[j].Dis }

func SortV(cont map[string]Path) []string {
	var sortedV []string
	vL := make(VLst, len(cont))
	i := 0
	for k, v := range cont {
		vL[i] = VerPair{k, v.Dis}
		i++
	}
	sort.Sort(vL)
	for _, vp := range vL {
		sortedV = append(sortedV, vp.ID)
	}
	return sortedV
}

func TransK2Hex(name string) string {
	return hex.EncodeToString([]byte(name))
}