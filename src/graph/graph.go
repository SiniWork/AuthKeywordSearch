/*
author: siyu
date: 2023/05/10
*/

package graph

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"os"
	"sort"
)

type Vertex struct {
	/*
		InEdges: the in edges of the vertex
		OutEdges: the out edges of the vertex
		Keywords: the keyword set of the vertex
	*/
	Id string
	InEdges []string
	OutEdges []string
	Keywords []string
}

func (v *Vertex) CompuVtxHash() []byte{
	raw := v.Raw()
	return Hash(raw)
}

func (v *Vertex) Raw() []byte {
	var valueStr []string
	sort.Strings(v.InEdges)
	sort.Strings(v.OutEdges)
	sort.Strings(v.Keywords)
	for _, edge := range v.InEdges {
		valueStr = append(valueStr, edge)
	}
	for _, edge := range v.OutEdges {
		valueStr = append(valueStr, edge)
	}
	for _, key := range v.Keywords {
		valueStr = append(valueStr, key)
	}
	raw := []interface{}{v.Id, valueStr}
	rlp, _ := rlp.EncodeToBytes(raw)
	return rlp
}

type Graph struct {
	/*
	Vertices: the dict to store the vertices and the key is the id
	Matrix: the matrix of the graph
	*/
	Vertices map[string]Vertex
	Matrix   map[string]map[string]bool
}

func Hash(raw []byte) []byte{
	h := md5.New()
	h.Write(raw)
	return h.Sum(nil)
}

func ReadGJsonFile(filePath string) map[string]Vertex {
	vertices := make(map[string]Vertex)
	jsonFile, err :=  os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	content, _ := io.ReadAll(jsonFile)
	var result []map[string]interface{}
	json.Unmarshal(content, &result)
	for _, vtx := range result {
		vertex := Vertex{
			Id:       vtx["id"].(string),
			InEdges:  nil,
			OutEdges: nil,
			Keywords: nil,
		}
		for _, val := range vtx["in"].([]interface{}) {
			vertex.InEdges = append(vertex.InEdges, val.(string))
		}
		for _, val := range vtx["out"].([]interface{}) {
			vertex.OutEdges = append(vertex.OutEdges, val.(string))
		}
		for _, val := range vtx["keywords"].([]interface{}) {
			vertex.Keywords = append(vertex.Keywords, val.(string))
		}
		vertices[vertex.Id] = vertex
	}
	return vertices
}

func ReadKJsonFile(filePath string) (map[string][]string, error) {
	invertIndex := make(map[string][]string)
	jsonFile, err :=  os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	content, _ := io.ReadAll(jsonFile)
	var result []map[string]interface{}
	json.Unmarshal(content, &result)
	for _, key := range result {
		for _, val := range key["ids"].([]interface{}) {
			invertIndex[key["key"].(string)] = append(invertIndex[key["key"].(string)], val.(string))
		}
	}
	return invertIndex, err
}