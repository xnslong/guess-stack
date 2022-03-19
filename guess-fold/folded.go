package main

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/xnslong/guess-stack/core/guess"
)

type stackElement string

func (s stackElement) EqualsTo(another guess.StackNode) bool {
	se, ok := another.(stackElement)
	if !ok {
		return false
	}

	return s == se
}

type foldedStack struct {
	Stack []guess.StackNode
	Value string
	*guess.StackExtraInfo
}

func (f *foldedStack) Path() []guess.StackNode {
	return f.Stack
}

func (f *foldedStack) SetPath(path []guess.StackNode) {
	f.Stack = path
}

type Profile struct {
	stacks []guess.Stack
}

func (p *Profile) Stacks() []guess.Stack {
	return p.stacks
}

func (p *Profile) WriteTo(writer io.Writer) error {
	bw := bufio.NewWriter(writer)
	defer bw.Flush()

	for _, stack := range p.stacks {
		fs := stack.(*foldedStack)
		err := FormatStack(fs, bw)
		if err != nil {
			return err
		}

		_, err = bw.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) ReadFrom(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	result := make([]guess.Stack, 0)
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
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid folded format")
	}

	stack := parts[0]

	stackElementStrList := strings.Split(stack, ";")
	stackElementList := make([]guess.StackNode, len(stackElementStrList))
	for i, v := range stackElementStrList {
		stackElementList[i] = stackElement(v)
	}

	return &foldedStack{
		Stack:          stackElementList,
		Value:          parts[1],
		StackExtraInfo: guess.NewStackExtraInfo(),
	}, nil
}

func FormatStack(stack *foldedStack, writer io.StringWriter) error {
	for i, element := range stack.Stack {
		if i > 0 {
			_, err := writer.WriteString(";")
			if err != nil {
				return err
			}
		}

		se := element.(stackElement)
		_, err := writer.WriteString(string(se))
		if err != nil {
			return err
		}
	}

	_, err := writer.WriteString(" ")
	if err != nil {
		return err
	}

	_, err = writer.WriteString(stack.Value)
	if err != nil {
		return err
	}
	return nil
}
