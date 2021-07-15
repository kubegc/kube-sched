/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package util

/**
 *   authors: wuheng@iscas.ac.cn
 *
 **/

type Elem struct {
	value interface{}
	prev *Elem
	next *Elem
}

type LinkedQueue struct {
	head *Elem
	tail *Elem
	size int
}

func NewLinkedQueue() *LinkedQueue {
	return &LinkedQueue{nil, nil, 0}
}

func (queue *LinkedQueue) Len() int {
	return queue.size
}

func (queue *LinkedQueue) Get() interface{} {
	if queue.head == nil {
		panic("Empty queue.")
	}
	return queue.head.value
}

func (queue *LinkedQueue) Add(value interface{}) {
	elem := &Elem{value, queue.tail, nil}
	if queue.tail == nil {
		queue.head = elem
		queue.tail = elem
	} else {
		queue.tail.next = elem
		queue.tail = elem
	}
	queue.size++
	elem = nil
}

func (queue *LinkedQueue) Remove() {
	if queue.head == nil {
		panic("Empty queue.")
	}
	elem := queue.head
	queue.head = elem.next
	elem.next = nil
	elem.value = nil
	queue.size--
	elem = nil
}
