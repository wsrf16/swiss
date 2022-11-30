package stack

import (
	"errors"
	"fmt"
)

type Stack struct {
	Element []interface{}
}

func NewStack() *Stack {
	return &Stack{}
}

func (stack *Stack) Push(value ...interface{}) {
	stack.Element = append(stack.Element, value...)
}

func (stack *Stack) Top() (value interface{}) {
	if stack.Size() > 0 {
		return stack.Element[stack.Size()-1]
	}
	return nil
}

func (stack *Stack) Pop() (interface{}, error) {
	if stack.Size() > 0 {
		value := stack.Element[stack.Size()-1]
		stack.Element = stack.Element[:stack.Size()-1]
		return value, nil
	}
	return nil, errors.New("stack is empty")
}

func (stack *Stack) Swap(other *Stack) {
	switch {
	case stack.Size() == 0 && other.Size() == 0:
		return
	case other.Size() == 0:
		other.Element = stack.Element[:stack.Size()]
		stack.Element = nil
	case stack.Size() == 0:
		stack.Element = other.Element
		other.Element = nil
	default:
		stack.Element, other.Element = other.Element, stack.Element
	}
	return
}

// 修改指定索引的元素
func (stack *Stack) Set(idx int, value interface{}) (err error) {
	if idx >= 0 && stack.Size() > 0 && stack.Size() > idx {
		stack.Element[idx] = value
		return nil
	}
	return errors.New("Set失败!")
}

// 返回指定索引的元素
func (stack *Stack) Get(idx int) (value interface{}) {
	if idx >= 0 && stack.Size() > 0 && stack.Size() > idx {
		return stack.Element[idx]
	}
	return nil
}

// Stack的size
func (stack *Stack) Size() int {
	return len(stack.Element)
}

// 是否为空
func (stack *Stack) Empty() bool {
	if stack.Element == nil || stack.Size() == 0 {
		return true
	}
	return false
}

// 打印
func (stack *Stack) Print() {
	for i := len(stack.Element) - 1; i >= 0; i-- {
		fmt.Println(i, "=>", stack.Element[i])
	}
}
