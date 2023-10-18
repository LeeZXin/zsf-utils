package trieutil

// TNode 前缀树

type TrieNode[T any] struct {
	label    string
	children map[rune]*TrieNode[T]
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
			children: make(map[rune]*TrieNode[T], 8),
		}
	}
	if key == "" {
		return
	}
	node := r.root
	for i, k := range key {
		if c, ok := node.children[k]; !ok {
			c = &TrieNode[T]{
				label:    key[:i+1],
				children: make(map[rune]*TrieNode[T], 8),
			}
			node.children[k] = c
		}
		node = node.children[k]
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
		node, ok = node.children[k]
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
		node, ok = node.children[k]
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
