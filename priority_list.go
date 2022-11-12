package xlog

type priorityItem[T any] struct {
	priority float32
	value    T
}

type byPriority[T any] []priorityItem[T]

func (a byPriority[T]) Len() int           { return len(a) }
func (a byPriority[T]) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority[T]) Less(i, j int) bool { return a[i].priority < a[j].priority }
