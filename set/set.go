package set

type Set[E comparable] map[E]struct{}

func Make[E comparable](v ...E) Set[E] {
	s := make(Set[E], len(v))
	for _, v := range v {
		s.Add(v)
	}
	return s
}

func Keys[K comparable, V any](m map[K]V) Set[K] {
	s := make(Set[K], len(m))
	for k := range m {
		s.Add(k)
	}
	return s
}

func Values[K, V comparable](m map[K]V) Set[V] {
	s := make(Set[V], len(m))
	for _, v := range m {
		s.Add(v)
	}
	return s
}

func Slurp[E comparable](ch <-chan E) Set[E] {
	s := make(Set[E])
	for v := range ch {
		s.Add(v)
	}
	return s
}

func (s Set[E]) Add(v E) bool {
	_, ok := s[v]
	s[v] = struct{}{}
	return !ok
}

func (s Set[E]) AddAll(t Set[E]) {
	for v := range t {
		s.Add(v)
	}
}

func (s Set[E]) Delete(v E) bool {
	_, ok := s[v]
	delete(s, v)
	return ok
}

func (s Set[E]) DeleteFunc(f func(E) bool) {
	for v := range s {
		if f(v) {
			s.Delete(v)
		}
	}
}

func (s Set[E]) Contains(v E) bool {
	_, ok := s[v]
	return ok
}

func (s Set[E]) ContainsFunc(f func(E) bool) bool {
	for v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

func (s Set[E]) Union(t Set[E]) Set[E] {
	out := make(Set[E])
	for v := range s {
		out.Add(v)
	}
	for v := range t {
		out.Add(v)
	}
	return out
}

func (s Set[E]) Intersect(t Set[E]) Set[E] {
	if len(s) > len(t) {
		s, t = t, s
	}
	out := make(Set[E])
	for v := range s {
		if t.Contains(v) {
			out.Add(v)
		}
	}
	return out
}

func (s Set[E]) Difference(t Set[E]) Set[E] {
	out := make(Set[E])
	for v := range s {
		if !t.Contains(v) {
			out.Add(v)
		}
	}
	return out
}

func (s Set[E]) SymmetricDifference(t Set[E]) Set[E] {
	out := make(Set[E])
	for v := range s {
		if !t.Contains(v) {
			out.Add(v)
		}
	}
	for v := range t {
		if !s.Contains(v) {
			out.Add(v)
		}
	}
	return out
}
