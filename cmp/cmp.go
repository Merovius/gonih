// Package cmp implements helpers for sorting composite types.
//
// The package is a drop-in replacement for the standard library [cmp] package
// (added in Go 1.21) with some extra helpers to support structs.
//
// The main focus of these helpers is convenience, not speed. They can be used
// with [slices.SortFunc] and [slices.BinarySearchFunc] (both added in Go 1.21),
// but are slower than hand-written versions by a factor of around 4 in casual
// benchmarks.
//
// The main helper is [Chain], which combines comparator functions to allow
// sorting structs by multiple fields. The rest are helpers to make it easy to
// extract fields and sort composite types like pointers and slices.
//
// TODO: Add benchmark.
package cmp

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
// If future releases of Go add new ordered types,
// this constraint will be modified to include them.
//
// Note that floating-point types may contain NaN ("not-a-number") values.
// An operator such as == or < will always report false when
// comparing a NaN value with any other value, NaN or not.
// See the [Compare] function for a consistent way to compare NaN values.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

// Less reports whether x is less than y.
// For floating-point types, a NaN is considered less than any non-NaN,
// and -0.0 is not less than (is equal to) 0.0.
func Less[T Ordered](x, y T) bool {
	return (isNaN(x) && !isNaN(y)) || x < y
}

// Compare returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
//
// For floating-point types, a NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
func Compare[T Ordered](x, y T) int {
	xNaN := isNaN(x)
	yNaN := isNaN(y)
	if xNaN && yNaN {
		return 0
	}
	if xNaN || x < y {
		return -1
	}
	if yNaN || x > y {
		return +1
	}
	return 0
}

// Bool compares two booleans, sorting false before true.
func Bool[T ~bool](x, y T) int {
	switch {
	case x == y:
		return 0
	case x == true:
		return 1
	default:
		return -1
	}
}

func id[T any](v T) T { return v }

// Cmp is a comparator function for T. It must return -1 if a < b, 0 if a == b
// and 1 if a > b. It is assignable to func(a, b T), so it can be used with
// [slices.Sort].
type Cmp[T any] func(a, b T) int

// Less returns if a is less than b. It can be used to plug helpers from this
// package into APIs that do not use three-way comparison operators.
func (cmp Cmp[T]) Less(a, b T) bool {
	return cmp(a, b) < 0
}

// Chain a sequence of comparison functions. The comparators are tried in
// sequence and the first non-zero result is returned. If all comparators
// consider the inputs equal, then so does the chain.
func Chain[T any](cmp ...Cmp[T]) Cmp[T] {
	return func(l, r T) int {
		for _, cmp := range cmp {
			if c := cmp(l, r); c < 0 {
				return -1
			} else if c > 0 {
				return 1
			}
		}
		return 0
	}
}

// isNaN reports whether x is a NaN without requiring the math package.
// This will always return false if T is not floating-point.
func isNaN[T Ordered](x T) bool {
	return x != x
}

// Slice compares two slices in lexicographical order. It is assignable to
// Cmp[S].
func Slice[S ~[]T, T Ordered](l, r S) int {
	n := len(l)
	if len(r) < n {
		n = len(r)
	}
	for i := 0; i < n; i++ {
		if c := Compare(l[i], r[i]); c != 0 {
			return c
		}
	}
	return Compare(len(l), len(r))
}

// SliceFunc compares two slices in lexicographical order, using a custom
// comparator.
func SliceFunc[S ~[]T, T any](cmp Cmp[T]) Cmp[S] {
	return func(l, r S) int {
		for i := 0; i < len(l) && i < len(r); i++ {
			if c := cmp(l[i], r[i]); c != 0 {
				return c
			}
		}
		return Compare(len(l), len(r))
	}
}

// Deref transforms a comparator of values into a comparator of pointers by
// derefencing the arguments.
//
// This combines well with ByPointerFunc, which guarantees that the comparator
// is only called with non-nil pointers. It is also useful if you have slices
// of non-nil pointers and a comparison function operating on values.
func Deref[T any](cmp Cmp[T]) Cmp[*T] {
	return func(l, r *T) int { return cmp(*l, *r) }
}

// PointerFunc is like Pointer but uses a custom comparator function.
//
// All nil pointers are considered equal and less than any non-nil pointer. Two
// non-nil pointers are compared using cmp.
//
// This is useful to sort pointers if you do not know whether or not they can be nil.
func PointerFunc[T any](cmp Cmp[*T]) Cmp[*T] {
	return ByPointerFunc(id[*T], cmp)
}

// Pointer compares two pointers to an ordered type.
//
// All nil pointers are considered equal and less than any non-nil pointer. Two
// non-nil pointers are compared using Compare.
//
// This is useful to sort pointers if you do not know whether or not they can be nil.
func Pointer[T Ordered](l, r *T) int {
	return PointerFunc(Deref(Compare[T]))(l, r)
}

// Reverse a comparator function.
func Reverse[T any](cmp Cmp[T]) Cmp[T] {
	return func(l, r T) int { return cmp(r, l) }
}

// ByFunc is like By, but uses a custom comparator for the field type.
func ByFunc[C, F any](by func(C) F, cmp Cmp[F]) Cmp[C] {
	return func(l, r C) int { return cmp(by(l), by(r)) }
}

// By compares a type by mapping it to an Ordered type and using Compare. The
// most common use is for C to be a struct and F to be a field type.
func By[C any, F Ordered](by func(C) F) Cmp[C] {
	return ByFunc(by, Compare[F])
}

// ByPointerFunc is like ByPointer, but uses a custom comparator function.
//
// All nil pointers are considered equal and less than any non-nil pointer. Two
// non-nil pointers are compared using cmp.
//
// It combines well with protocol buffers and other code generators, which
// create getters for message fields.
func ByPointerFunc[C, F any](by func(C) *F, cmp Cmp[*F]) Cmp[C] {
	return func(l, r C) int {
		pl, pr := by(l), by(r)
		switch {
		case pl == nil && pr == nil:
			return 0
		case pl == nil:
			return -1
		case pr == nil:
			return 1
		default:
			return cmp(pl, pr)
		}
	}
}

// ByPointer is like By but operates on pointer fields.
//
// All nil pointers are considered equal and less than any non-nil pointer. Two
// non-nil pointers are compared using Compare.
//
// It combines well with protocol buffers and other code generators, which
// create getters for message fields.
func ByPointer[C any, F Ordered](by func(C) *F) Cmp[C] {
	return ByPointerFunc(by, Deref(Compare[F]))
}
