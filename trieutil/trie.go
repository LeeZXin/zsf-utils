package trieutil

import "github.com/LeeZXin/zsf-utils/collections/hashmap"

// TNode 前缀树

type TrieNode[T any] struct {
	label    string
	children *hashmap.HashMap[rune, *TrieNode[T]]
	data     T
	has      bool
}

const (
	// LongestMatchType 最长匹配
	LongestMatchType = iota + 1
	// ShortestMatchType 最短匹配
	ShortestMatchType
)

// Trie 通用前缀匹配树
type Trie[T any] struct {
	root *TrieNode[T]
}

// Insert 插入
func (r *Trie[T]) Insert(key string, data T) {
	if r.root == nil {
		r.root = &TrieNode[T]{
			children: hashmap.NewHashMap[rune, *TrieNode[T]](),
		}
	}
	if key == "" {
		return
	}
	node := r.root
	for i, k := range key {
		if !node.children.Contains(k) {
			c := &TrieNode[T]{
				label:    key[:i+1],
				children: hashmap.NewHashMap[rune, *TrieNode[T]](),
			}
			node.children.Put(k, c)
		}
		node, _ = node.children.Get(k)
	}
	node.data = data
	node.has = true
}

// FullSearch 完全匹配
func (r *Trie[T]) FullSearch(key string) (T, bool) {
	if r.root == nil {
		var t T
		return t, false
	}
	node := r.root
	for _, k := range key {
		var ok bool
		node, ok = node.children.Get(k)
		if !ok {
			var t T
			return t, false
		}
	}
	if !node.has {
		var t T
		return t, false
	}
	return node.data, true
}

// PrefixSearch 前缀匹配
func (r *Trie[T]) PrefixSearch(key string, matchType int) (T, bool) {
	if r.root == nil {
		var t T
		return t, false
	}
	node := r.root
	list := make([]TrieNode[T], 0, 8)
	var ok bool
	for _, k := range key {
		node, ok = node.children.Get(k)
		if !ok {
			break
		}
		if node.has {
			if matchType == ShortestMatchType {
				return node.data, true
			}
			list = append(list, *node)
		}
	}
	if len(list) == 0 {
		var t T
		return t, false
	}
	//最长匹配
	return list[len(list)-1].data, true
}
