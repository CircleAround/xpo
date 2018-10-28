package stdkit

type ImmutableStringSet struct {
	values map[string]struct{}
	src    []string
}

func (l *ImmutableStringSet) Src() []string {
	ret := make([]string, len(l.src))
	copy(ret, l.src)
	return ret
}

func (l *ImmutableStringSet) ToArray() []string {
	ret := make([]string, len(l.values))
	index := 0
	for value, _ := range l.values {
		ret[index] = value
		index++
	}
	return ret
}

func (l *ImmutableStringSet) Contains(v string) bool {
	_, ok := l.values[v]
	return ok
}

func (l *ImmutableStringSet) Size() int {
	return len(l.values)
}

func NewImmutableStringSet(src []string) *ImmutableStringSet {
	values := make(map[string]struct{})
	for _, value := range src {
		values[value] = struct{}{}
	}

	return &ImmutableStringSet{
		values: values,
		src:    src,
	}
}
