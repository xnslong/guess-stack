package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/xnslong/guess-stack/core"
)

type stackElement string

func (s stackElement) EqualsTo(another core.StackNode) bool {
	se, ok := another.(stackElement)
	if !ok {
		return false
	}

	return s == se
}

type foldedStack struct {
	Stack []stackElement
	Count int
	need  bool
}

func (f *foldedStack) NeedFix() bool {
	return f.need
}

func (f *foldedStack) SetNeedFix(need bool) {
	f.need = need
}

func (f *foldedStack) Path() []core.StackNode {
	nodes := make([]core.StackNode, len(f.Stack))

	for i := 0; i < len(f.Stack); i++ {
		nodes[i] = f.Stack[i]
	}

	return nodes
}

func (f *foldedStack) SetPath(path []core.StackNode) {
	result := make([]stackElement, len(path))

	for i := 0; i < len(path); i++ {
		result[i] = path[i].(stackElement)
	}

	f.Stack = result
}

type Profile struct {
	stacks []core.Stack
}

func (p *Profile) Stacks() []core.Stack {
	return p.stacks
}

func (p *Profile) WriteTo(writer io.Writer) error {
	for _, stack := range p.stacks {
		fs := stack.(*foldedStack)
		err := FormatStack(fs, writer)
		if err != nil {
			return err
		}

		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) ReadFrom(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	result := make([]core.Stack, 0)
	for scanner.Scan() {
		stack, err := ParseStack(scanner.Text())
		if err != nil {
			return err
		}

		result = append(result, stack)
	}

	p.stacks = result
	return nil
}

func ParseStack(line string) (*foldedStack, error) {
	index := strings.LastIndex(line, " ")

	if index < 0 {
		return nil, errors.New("invalid folded format")
	}

	val := line[index+1:]
	valInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, errors.New("invalid folded format")
	}

	stack := line[:index]

	stackElementStrList := strings.Split(stack, ";")
	stackElementList := make([]stackElement, len(stackElementStrList))
	for i, v := range stackElementStrList {
		stackElementList[i] = stackElement(v)
	}

	return &foldedStack{
		Stack: stackElementList,
		Count: int(valInt),
		need:  true,
	}, nil
}

func FormatStack(stack *foldedStack, writer io.Writer) error {
	for i, element := range stack.Stack {
		if i > 0 {
			_, err := writer.Write([]byte(";"))
			if err != nil {
				return err
			}
		}
		_, err := writer.Write([]byte(element))
		if err != nil {
			return err
		}
	}

	_, err := writer.Write([]byte(" "))
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(strconv.Itoa(stack.Count)))
	if err != nil {
		return err
	}
	return nil
}
