/*
author: siyu
date: 2023/04/24
*/

package advanced

import (
	"Collie/src/graph"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

const HOP = 4

type GraphAdv struct {
	graph.Graph
	InvertList map[string][]string
	MPBT *Trie
}

func (g *GraphAdv) LoadGraph(gF, kF string) {
	g.Vertices = graph.ReadGJsonFile(gF)
	g.InvertList, _ = graph.ReadKJsonFile(kF)
	g.MPBT = NewTrie()
}

func (g *GraphAdv) Preprocessing(preFile string) {
	/*
	Generate the keyword.json, which saves the Bucket for each keyword,
	the Buckets saves the vertices that can reach to the keyword
	 */
	buckets := make(map[string]Bucket) // key: bucket
	for key, _ := range g.InvertList {
		buckets[key] = Bucket{
			Cont: make(map[string]Path),
		}
	}
	for v, _ := range g.Vertices {
		paths := g.ExpandVertex(v)
		for k, p := range paths {
			if _, exist := buckets[k].Cont[v]; !exist {
				buckets[k].Cont[v] = p
			} else if buckets[k].Cont[v].Dis < p.Dis {
				buckets[k].Cont[v] = p
			}
		}
	}

	// write to json file
	filePtr, err := os.Create(preFile)
	if err != nil {
		fmt.Println("create failed", err.Error())
		return
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(buckets)
}

func (g *GraphAdv) ExpandVertex(id string) map[string]Path {
	paths := make(map[string]Path)
	for _, k := range g.Vertices[id].Keywords {
		p := Path{
			TargetN: id,
			Dis: 0,
		}
		paths[k] = p
	}
	hops := make(map[int][]string)
	hops[1] = g.Vertices[id].OutEdges

	for i := 1; i < HOP; i++ {
		for _, v := range hops[i] {
			for _, k := range g.Vertices[v].Keywords {
				if _, visited := paths[k]; !visited {
					p := Path{
						TargetN: v,
						Dis: i,
					}
					paths[k] = p
				}
			}
			hops[i+1] = append(hops[i+1], g.Vertices[v].OutEdges...)
		}
	}
	return paths
}

func StatisticVClasses(vScopeF string) {
	vertices := make(map[string]map[string]int)

	// load vertices
	vFile, err := os.Open(vScopeF)
	if err != nil {
		fmt.Println(err)
	}
	defer vFile.Close()
	decoderFilter := json.NewDecoder(vFile)
	decoderFilter.Decode(&vertices)

	// 统计有多少种vertex

	i := 0
	for id, _ := range vertices {
		if len(vertices[id]) > 100 {
			i++
		}
	}
	fmt.Println(i)
}

// Dominate filter

type VerK struct {
	Id string
	Num int
}

type VerKLst []VerK

func (p VerKLst) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p VerKLst) Len() int           { return len(p) }
func (p VerKLst) Less(i, j int) bool { return p[i].Num < p[j].Num }

func LoadFilterPreProcess(preProF, verKF, outF string) {
	rawBuckets := make(map[string]Bucket)
	filter := make(map[string]bool)
	fineBuckets := make(map[string]Bucket)

	// load raw buckets
	buckFile, err := os.Open(preProF)
	if err != nil {
		fmt.Println(err)
	}
	defer buckFile.Close()
	decoderRaw := json.NewDecoder(buckFile)
	decoderRaw.Decode(&rawBuckets)

	// load filter
	verKFile, err := os.Open(verKF)
	if err != nil {
		fmt.Println(err)
	}
	defer buckFile.Close()
	decoderFilter := json.NewDecoder(verKFile)
	decoderFilter.Decode(&filter)

	// filtering buckets
	for k, b := range rawBuckets {
		fmt.Println(k, len(b.Cont))
		fineB := FilterBucket(b, filter)
		fineBuckets[TransK2Hex(k)] = fineB
		fmt.Println(k, len(fineB.Cont))
	}

	// write filtered buckets
	fineFile, err := os.Create(outF)
	if err != nil {
		fmt.Println("create failed", err.Error())
		return
	}
	defer fineFile.Close()
	encoder := json.NewEncoder(fineFile)
	err = encoder.Encode(fineBuckets)
}

func FilterBucket(b Bucket, filter map[string]bool) Bucket{
	fineBuck := Bucket{Cont: make(map[string]Path)}

	for v, _ := range b.Cont {
		_, yes := filter[v]; if !yes {
			fineBuck.Cont[v] = b.Cont[v]
		}
	}
	return fineBuck
}

// dominate test

func StatisticDominate(file string) {
	/*
	Statistic the number of dominant vertex
	 */
	vKScope := make(map[string]map[string]int)
	DomFlag := make(map[string]bool)

	verKFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer verKFile.Close()
	decoderFilter := json.NewDecoder(verKFile)
	decoderFilter.Decode(&vKScope)

	vKL := VerKLst{}
	for v, ks := range vKScope {
		vKL = append(vKL, VerK{Id: v, Num: len(ks)})
	}
	sort.Sort(vKL)

	t := time.Now()
	for i := 0; i < len(vKL); i++ {
		_, exist := DomFlag[vKL[i].Id]; if !exist {
			for j := i+1; j < len(vKL); j++ {
				if IsDominate(vKScope[vKL[i].Id], vKScope[vKL[j].Id]) {
					DomFlag[vKL[j].Id] = true
				}
			}
		}
	}
	e := time.Since(t)
	fmt.Println(e)
	fmt.Println("All: ", len(vKL))
	fmt.Println("Now: ", len(vKL)-len(DomFlag))

	// write to file
	fineFile, err := os.Create("removed.json")
	if err != nil {
		fmt.Println("create failed", err.Error())
		return
	}
	defer fineFile.Close()
	encoder := json.NewEncoder(fineFile)
	err = encoder.Encode(DomFlag)
}

func (g *GraphAdv) GenKDistribution(v string) map[string]int {
	/*
	Statistic the distribution of keywords in the specific hops of v
	 */
	scope := make(map[string]int)
	for _, k := range g.Vertices[v].Keywords {
		scope[k] = 0
	}
	hops := make(map[int][]string)
	hops[1] = g.Vertices[v].OutEdges

	for i := 1; i < HOP; i++ {
		for _, n := range hops[i] {
			for _, k := range g.Vertices[n].Keywords {
				if _, visited := scope[k]; !visited {
					scope[k] = i
				}
			}
			hops[i+1] = append(hops[i+1], g.Vertices[v].OutEdges...)
		}
	}

	return scope
}

func IsDominate(v1, v2 map[string]int) bool {
	for k, _ := range v2 {
		_, exist := v1[k]; if !exist {
			return false
		} else {
			if v1[k] > v2[k] {
				return false
			}
		}
	}
	return true
}