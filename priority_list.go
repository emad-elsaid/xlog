package xlog

type plist[T any] struct {
	priority float32
	value    T
	next     *plist[T]
}

func (p *plist[T]) insert(priority float32, value T) *plist[T] {
	node := plist[T]{
		priority: priority,
		value:    value,
	}

	if p == nil {
		return &node
	} else if p.priority > priority {
		node.next = p
		return &node
	}

	var prev *plist[T]
	for prev = p; prev.next != nil && prev.next.priority <= priority; prev = prev.next {
	}

	node.next = prev.next
	prev.next = &node

	return p
}

func (p *plist[T]) each(f func(t T)) {
	for i := p; i != nil; i = i.next {
		f(i.value)
	}
}
