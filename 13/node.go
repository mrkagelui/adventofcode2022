package main

import (
	"fmt"
	"strconv"
	"strings"
)

type node struct {
	value int
	subs  []node
}

func (n node) isValue() bool {
	return n.subs == nil
}

func parseNode(str string) (node, error) {
	n, i, err := parse(str)
	if err != nil {
		return node{}, err
	}
	if i != len(str) {
		return node{}, fmt.Errorf("wrong bytes scanned for [%v]", str)
	}
	return n, nil
}

func parse(str string) (node, int, error) {
	currIndex := 1 // because the first byte is always '[', can skip
	currentNodes := make([]node, 0)
	var buffer strings.Builder

	for currIndex < len(str) {
		increment := 1
		switch b := str[currIndex]; b {
		case '[':
			n, inc, err := parse(str[currIndex:])
			if err != nil {
				return node{}, 0, err
			}
			increment = inc

			currentNodes = append(currentNodes, n)
		case ',':
			if buffer.Len() != 0 {
				num, err := strconv.Atoi(buffer.String())
				if err != nil {
					return node{}, 0, err
				}
				currentNodes = append(currentNodes, node{value: num})

				buffer.Reset()
			}
		case ']':
			if buffer.Len() != 0 {
				num, err := strconv.Atoi(buffer.String())
				if err != nil {
					return node{}, 0, err
				}
				currentNodes = append(currentNodes, node{value: num})

				buffer.Reset()
			}

			return node{subs: currentNodes}, currIndex + 1, nil
		default:
			buffer.WriteByte(b)
		}
		currIndex += increment
	}

	return node{}, 0, fmt.Errorf("malformed part: {%v}", str)
}

func smallerThan(a, b node) *bool {
	switch {
	case a.isValue() && b.isValue():
		switch {
		case a.value == b.value:
			return nil
		case a.value > b.value:
			return ptrOf(false)
		default:
			return ptrOf(true)
		}
	case a.isValue() && !b.isValue():
		return smallerThan(node{subs: []node{a}}, b)
	case !a.isValue() && b.isValue():
		return smallerThan(a, node{subs: []node{b}})
	default:
		for i := 0; i < len(a.subs) || i < len(b.subs); i++ {
			if len(a.subs) < len(b.subs) && i == len(a.subs) {
				return ptrOf(true)
			}
			if len(b.subs) < len(a.subs) && i == len(b.subs) {
				return ptrOf(false)
			}
			if s := smallerThan(a.subs[i], b.subs[i]); s != nil {
				return s
			}
		}

	}
	return nil
}

func ptrOf(b bool) *bool {
	return &b
}
