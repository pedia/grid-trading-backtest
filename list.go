package grid

type list[T any] []T

func List[T any](ts ...T) list[T] {
	return ts
}

func (a *list[T]) Add(ts ...T) *list[T] {
	*a = append(*a, ts...)
	return a
}

func (a *list[T]) Insert(i int, t T) *list[T] {
	// extend array by one
	*a = append(*a, t)

	// shift values
	copy((*a)[i+1:], (*a)[i:])

	// insert value
	(*a)[i] = t
	return a
}

func (a *list[T]) RemoveAt(i int) *list[T] {
	copy((*a)[i:], (*a)[i+1:])
	*a = (*a)[:len((*a))-1]
	return a
}

func (a *list[T]) SubList(start int, ends ...int) list[T] {
	end := a.Len()
	if len(ends) > 0 {
		end = ends[0]
	}
	return (*a)[start:end]
}

func (a *list[T]) IsEmpty() bool {
	return len(*a) == 0
}

func (a *list[T]) Len() int {
	return len(*a)
}

func (a *list[T]) Items() []T {
	return *a
}

func (a *list[T]) Map(iter func(T) any) []any {
	out := make([]any, 0, a.Len())
	for _, i := range *a {
		out = append(out, iter(i))
	}
	return out
}

func (a *list[T]) Clear() *list[T] {
	*a = make([]T, 0)
	return a
}

// add
// addAll
// any
// asMap
// cast
// clear
// contains
// elementAt
// every
// expand
// fillRange
// firstWhere
// fold
// followedBy
// forEach
// getRange
// indexOf
// indexWhere
// insert
// insertAll
// join
// lastIndexOf
// lastIndexWhere
// lastWhere
// map
// noSuchMethod
// reduce
// remove
// removeAt
// removeLast
// removeRange
// removeWhere
// replaceRange
// retainWhere
// setAll
// setRange
// shuffle
// singleWhere
// skip
// skipWhile
// sort
// sublist
// take
// takeWhile
// toList
// toSet
// toString
// where
// whereType

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T any, R any](collection []T, iteratee func(item T, index int) R) list[R] {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}
	return result
}
