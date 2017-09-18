package main

import (
	"container/heap"
	"flag"
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	Char               string
	Weight             byte
	SubtreeR, SubtreeL *Node
}

type Nodes []*Node

type Code struct {
	Char   string
	Value  string
	Weight byte
}

// Satisfying container/heap interface
func (ch Nodes) Len() int           { return len(ch) }
func (ch Nodes) Less(i, j int) bool { return ch[i].Weight < ch[j].Weight }
func (ch Nodes) Swap(i, j int)      { ch[i], ch[j] = ch[j], ch[i] }

func (n *Nodes) Push(x interface{}) {
	*n = append(*n, x.(*Node))
}

func (ch *Nodes) Pop() interface{} {
	old := *ch
	n := len(old)
	x := old[n-1]
	*ch = old[0 : n-1]
	return x
}

func main() {
	TreeHeap := &Nodes{}
	heap.Init(TreeHeap)

	input := flag.String("string", "", "")
	flag.Parse()

	str := strings.ToLower(*input)
	hs := handleString(str)

	for k, v := range hs {
		heap.Push(TreeHeap, &Node{
			Char:   k,
			Weight: byte(v),
		})
	}

	for TreeHeap.Len() > 1 {
		st1 := heap.Pop(TreeHeap).(*Node)
		st2 := heap.Pop(TreeHeap).(*Node)

		heap.Push(TreeHeap, &Node{
			SubtreeR: st1,
			SubtreeL: st2,
			Weight:   st1.Weight + st2.Weight,
		})
	}

	// Pop heap and call traverse
	hp := heap.Pop(TreeHeap).(*Node)
	codes := hp.recursiveTraversal(nil, nil)

	// Sort 'codes' back into order, matching string char index
	sort.Slice(codes, func(i, j int) bool {
		return strings.Index(str, codes[i].Char) <= strings.Index(str, codes[j].Char)
	})

	encoded := encode(str, codes)
	decoded := hp.Decode(encoded)

	for _, c := range codes {
		fmt.Println(c)
	}

	fmt.Printf("%v\n%v\n%v\n", encoded, decoded, string(decoded))
}

// Create character:weight map
func handleString(str string) map[string]int {
	mapped := make(map[string]int)
	for _, char := range strings.Split(str, "") {
		mapped[char] = strings.Count(str, char)
	}

	return mapped
}

// Recursively traverse through nodes, build node subtrees
var cp []*Code
func (n *Node) recursiveTraversal(list []byte, c []*Code) []*Code {
	if n.SubtreeR == nil && n.SubtreeL == nil {
		c = append(c, &Code{
			Char:   n.Char,
			Value:  string(list),
			Weight: byte(n.Weight),
		})
	} else {
		n.SubtreeR.recursiveTraversal(append(list, byte('1')), nil) // Right
		n.SubtreeL.recursiveTraversal(append(list, byte('0')), nil) // Left
	}

	if len(c) != 0 {
		cp = append(cp, c...)
	}

	return cp
}

// Encode string to binary
func encode(str string, codes []*Code) string {
	var encoded string
	for _, k := range str {
		for _, kv := range codes {
			if string(k) == kv.Char {
				encoded += kv.Value
			}
		}
	}

	encoded += " " // Empty bit
	return encoded
}

// Traversing back through the encoded string, left to right
func (node *Node) Decode(encoded string) []byte {
	var decoded []byte
	n := node

	for i := 0; i < len(encoded); i++ {
		if n.SubtreeL == nil && n.SubtreeR == nil {
			decoded = append(decoded, []byte(n.Char)...)
			n = node
		}

		switch byte(encoded[i]) {
		case byte('0'):
			n = n.SubtreeL
		case byte('1'):
			n = n.SubtreeR
		}
	}

	return decoded
}
