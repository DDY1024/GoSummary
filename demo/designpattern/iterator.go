package main

type Iterator interface {
	HasNext() bool
	Next()
	CurrentItem() interface{}
}

type ArrayInt []int

func (a ArrayInt) Iterator() Iterator {
	return &ArrayIntIterator{
		arrayInt: a,
		index:    0,
	}
}

type ArrayIntIterator struct {
	arrayInt ArrayInt
	index    int
}

func (iter *ArrayIntIterator) HasNext() bool {
	return iter.index < len(iter.arrayInt)-1
}

func (iter *ArrayIntIterator) Next() {
	iter.index++
}

func (iter *ArrayIntIterator) CurrentItem() interface{} {
	return iter.arrayInt[iter.index]
}
