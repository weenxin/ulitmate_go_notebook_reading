package section10_hash

type hashFunction[K comparable] func(key K, buckets int) int

type KeyValuePair [K comparable, V any] struct{
	key K
	value V
}

type Table[K comparable, V any] struct {
	hashFunc hashFunction[K]
	buckets int
	data [][]KeyValuePair[K,V]
}

func New[K comparable, V any](buckets int, function hashFunction[K]) *Table[K,V]{
	return &Table[K,V]{
		hashFunc: function,
		buckets:  buckets,
		data:     make([][]KeyValuePair[K,V], buckets),
	}
}

func (t * Table[K,V]) Insert(key K, value V){
	bucket := t.hashFunc(key,t.buckets)
	for index, theKey := range t.data[bucket] {
		if key == theKey.key {
			t.data[bucket][index].value = value
			return
		}
	}
	pair := KeyValuePair[K,V] {
		key: key,
		value: value,
	}
	t.data[bucket] = append(t.data[bucket], pair)
}

func (t * Table[K,V]) Get(key K) (V, bool) {
	bucket := t.hashFunc(key,t.buckets)
	for index, theKey := range t.data[bucket] {
		if key == theKey.key {
			return t.data[bucket][index].value, true
		}
	}
	var zero  V
	return zero, false
}
