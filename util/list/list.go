package list

import (
	"iter"
	"slices"
)

type List[T any] struct{ slice []T }

func Of[T any](values ...T) *List[T] {
	return &List[T]{slice: values}
}
func New[T any](capacity int) *List[T] {
	return &List[T]{slice: make([]T, 0, capacity)}
}
func (this *List[T]) Add(values ...T) {
	this.slice = append(this.slice, values...)
}
func (this *List[T]) Len() int {
	if this == nil {
		return 0
	}
	return len(this.slice)
}
func (this *List[T]) Clear() {
	clear(this.slice)
	this.slice = this.slice[:0]
}
func (this *List[T]) All() iter.Seq[T] {
	return slices.Values(this.slice)
}
