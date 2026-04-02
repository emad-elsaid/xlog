package autolink_pages

import (
	"sync"

	. "github.com/emad-elsaid/xlog"
)

// Simple trie node for efficient prefix matching
type trieNode struct {
	children map[rune]*trieNode
	page     Page // nil for intermediate nodes
	isEnd    bool // true if this node represents end of a page name
}

type trie struct {
	root *trieNode
	mu   sync.RWMutex // Protect concurrent access during building
}

func newTrie() *trie {
	return &trie{
		root: &trieNode{
			children: make(map[rune]*trieNode),
		},
	}
}

// insert adds a page name to the trie (normalized/lowercase)
// Thread-safe for concurrent inserts
func (t *trie) insert(name string, page Page) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	runes := []rune(name)

	for _, r := range runes {
		if node.children == nil {
			node.children = make(map[rune]*trieNode)
		}

		if _, exists := node.children[r]; !exists {
			node.children[r] = &trieNode{
				children: make(map[rune]*trieNode),
			}
		}
		node = node.children[r]
	}

	node.isEnd = true
	node.page = page
}

// search finds the longest matching page name prefix in the input
// Returns (page, matchLength in bytes) or (nil, 0) if no match
func (t *trie) search(input string) (Page, int) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	runes := []rune(input)

	var lastMatch Page
	var lastMatchByteLen int

	bytePos := 0
	for _, r := range runes {
		if node.children == nil {
			break
		}

		next, exists := node.children[r]
		if !exists {
			break
		}

		node = next
		bytePos += len(string(r)) // Track byte position

		// If this is end of a page name, record it
		if node.isEnd {
			lastMatch = node.page
			lastMatchByteLen = bytePos // Store byte length, not rune count
		}
	}

	return lastMatch, lastMatchByteLen
}
