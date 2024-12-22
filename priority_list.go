package xlog

import (
	"iter"
	"sort"
)

type priorityItem[T any] struct {
	Item     T
	Priority float32
}

type priorityList[T any] struct {
	items []priorityItem[T]
}

func (pl *priorityList[T]) Add(item T, priority float32) {
	pl.items = append(pl.items, priorityItem[T]{Item: item, Priority: priority})
	pl.sortByPriority()
}

func (pl *priorityList[T]) sortByPriority() {
	sort.Slice(pl.items, func(i, j int) bool {
		return pl.items[i].Priority < pl.items[j].Priority
	})
}

// An iterator over all items
func (pl *priorityList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range pl.items {
			if !yield(v.Item) {
				return
			}
		}
	}
}
