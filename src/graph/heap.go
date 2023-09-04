/*
author: siyu
date: 2023/03/10
*/

package graph

import (
	"container/heap"
	"fmt"
)


// Define a min-heap to maintain iterators

type Iterator []string
type IterHeap []Iterator

func (h IterHeap) Len() int {
	return len(h)
}

func (h IterHeap) Less(i, j int) bool {
	return len(h[i]) < len(h[j])
	//return h[i].Dist < h[j].Dist
}

func (h IterHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IterHeap) Push(iter interface{}) {
	*h = append(*h, iter.(Iterator))
}

func (h *IterHeap) Pop() (iter interface{}) {
	old := *h
	n := len(*h)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (h *IterHeap) PushIter(iter Iterator) {
	heap.Push(h, iter)
}

func (h *IterHeap) PopIter() (iter Iterator) {
	if len(*h) == 0 {
		return Iterator{}
	} else {
		iter = heap.Pop(h).(Iterator)
		return iter
	}
}

func (h *IterHeap) Top() (iter Iterator) {
	if len(*h) == 0 {
		return Iterator{}
	} else {
		iter = (*h)[0]
		return iter
	}
}

// Define a min-heap to maintain iterator heaps

type Item struct {
	Keyword string
	Priority int
	IHp      IterHeap
}

type PriorityQueue []Item

func (p PriorityQueue) Len() int {
	return len(p)
}

func (p PriorityQueue) Less(i, j int) bool {
	return p[i].Priority < p[j].Priority
}

func (p PriorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *PriorityQueue) Push(item interface{}) {
	*p = append(*p, item.(Item))
}

func (p *PriorityQueue) Pop() (item interface{}) {
	old := *p
	n := len(*p)
	x := old[n-1]
	*p = old[:n-1]
	return x
}

func (p *PriorityQueue) PushItem(item Item) {
	heap.Push(p, item)
}

func (p *PriorityQueue) PopItem() (item Item) {
	item = heap.Pop(p).(Item)
	return item
}

func (p *PriorityQueue) Top() (item Item) {
	item = (*p)[0]
	return item
}


// Define a min-heap to maintain all result trees

type Leaf struct {
	Id string
	Dist int
}

type ResultTree struct {
	Root string
	Leaves map[string]Leaf  // key: the keyword
	//Score int
}

func (t *ResultTree) computeScore() int {
	score := 0
	for _, leaf := range t.Leaves {
		score += leaf.Dist
	}
	return score
}

type ResTrees []ResultTree

func (r ResTrees) Len() int {
	return len(r)
}

func (r ResTrees) Less(i, j int) bool {
	return r[i].computeScore() < r[j].computeScore()
}

func (r ResTrees) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r *ResTrees) Push(tree interface{}) {
	*r = append(*r, tree.(ResultTree))
}

func (r *ResTrees) Pop() (tree interface{}) {
	old := *r
	n := len(*r)
	x := old[n-1]
	*r = old[:n-1]
	return x
}

func (r *ResTrees) PushRes(tree ResultTree) {
	heap.Push(r, tree)
}

func (r *ResTrees) PopRes() (tree ResultTree) {
	tree = heap.Pop(r).(ResultTree)
	return tree
}

func (r *ResTrees) Top() (tree ResultTree) {
	tree = (*r)[0]
	return tree
}

func TestHeap() {
	hp1 := &IterHeap{}
	heap.Init(hp1)
	p1 := Iterator{"0001", "0002", "0003"}
	p2 := Iterator{"0001", "0003"}
	p3 := Iterator{"0001", "0002", "0003", "0006", "0009"}
	hp1.PushIter(p1)
	hp1.PushIter(p2)
	hp1.PushIter(p3)
	fmt.Println(hp1.Top())

	hp2 := &IterHeap{}
	heap.Init(hp2)
	p11 := Iterator{"0001", "0002", "0003", "0006", "0009", "0006", "0009"}
	p22 := Iterator{"0001"}
	p33 := Iterator{"0001", "0002", "0003", "0006", "0009"}
	hp2.PushIter(p11)
	hp2.PushIter(p22)
	hp2.PushIter(p33)
	fmt.Println(hp2.Top())

	pq := &PriorityQueue{}
	pq.PushItem(Item{Priority: len(hp1.Top()), IHp: *hp1})
	pq.PushItem(Item{Priority: len(hp2.Top()), IHp: *hp2})
	fmt.Println(pq.Top())
}

