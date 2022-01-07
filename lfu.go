package cachelfu

import (
	"container/list"
	"errors"
	"fmt"
)

type Cache interface {
	Add(key string, el interface{}) error
	Get(key string) (interface{}, error)
	Remove(key string) error
}

type frequencyNode struct {
	values *list.List
	freq   int
}

type valueNode struct {
	key       string
	data      interface{}
	frequency *list.Element
}

type cache struct {
	values   map[string]*list.Element
	freqs    *list.List
	capacity int
}

func (c *cache) Add(key string, el interface{}) error {
	// if key already exists return error
	if _, ok := c.values[key]; ok {
		return errors.New(fmt.Sprintf(`Key "%s" already exists`, key))
	}

	c.evict()

	// create new value node
	newValueNode := &valueNode{
		data: el,
		key:  key,
	}

	// check if frequency node with freq 1 exists, if not create it
	firstFreqNodeEl := c.freqs.Front()
	if firstFreqNodeEl == nil || firstFreqNodeEl.Value.(*frequencyNode).freq != 1 {
		firstFreqNode := &frequencyNode{
			values: list.New(),
			freq:   1,
		}

		firstFreqNodeEl = c.freqs.PushFront(firstFreqNode)
	}

	newValueNode.frequency = firstFreqNodeEl

	newValueNodeEl := firstFreqNodeEl.Value.(*frequencyNode).values.PushFront(newValueNode)

	c.values[key] = newValueNodeEl

	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
	var valNodeEl *list.Element
	if v, ok := c.values[key]; !ok {
		return nil, errors.New(fmt.Sprintf(`Key "%s" is not found`, key))
	} else {
		valNodeEl = v
	}

	freqNodeEl := valNodeEl.Value.(*valueNode).frequency
	nextFreqNodeEl := freqNodeEl.Next()

	freqNode := freqNodeEl.Value.(*frequencyNode)
	var nextFreqNode *frequencyNode

	if nextFreqNodeEl != nil {
		nextFreqNode = nextFreqNodeEl.Value.(*frequencyNode)
	}

	// create next frequency node if it not exists
	if nextFreqNode == nil || nextFreqNode.freq != freqNode.freq+1 {
		newFreqNode := &frequencyNode{
			values: list.New(),
			freq:   freqNode.freq + 1,
		}

		nextFreqNode = newFreqNode
		nextFreqNodeEl = c.freqs.InsertAfter(newFreqNode, freqNodeEl)
	}

	nextFreqNode.values.PushBack(valNodeEl)
	valNodeEl.Value.(*valueNode).frequency = nextFreqNodeEl

	return valNodeEl.Value.(*valueNode).data, nil
}

func (c *cache) Remove(key string) error {
	var removingNodeEl *list.Element
	if v, ok := c.values[key]; !ok {
		return errors.New(fmt.Sprintf(`Key "%s" is not found`, key))
	} else {
		removingNodeEl = v
	}

	frequencyNodeEl := removingNodeEl.Value.(*valueNode).frequency

	frequencyNodeEl.Value.(*frequencyNode).values.Remove(removingNodeEl)

	if frequencyNodeEl.Value.(*frequencyNode).values.Front() == nil {
		c.freqs.Remove(frequencyNodeEl)
	}

	// remove element from map
	delete(c.values, key)

	return nil
}

func (c *cache) evict() {
	for len(c.values) >= c.capacity {
		if c.freqs.Front() != nil && c.freqs.Front().Value.(*frequencyNode).values.Front() != nil {
			removingNodeEl := c.freqs.Front().Value.(*frequencyNode).values.Front()

			c.freqs.Front().Value.(*frequencyNode).values.Remove(removingNodeEl)

			// remove element from map
			delete(c.values, removingNodeEl.Value.(*valueNode).key)
		}

		if c.freqs.Front() != nil && c.freqs.Front().Value.(*frequencyNode).values.Front() == nil {
			c.freqs.Remove(c.freqs.Front())
		}
	}
}

func New(cap int) Cache {
	c := new(cache)
	c.values = make(map[string]*list.Element)
	c.freqs = list.New()
	c.capacity = cap

	return c
}
