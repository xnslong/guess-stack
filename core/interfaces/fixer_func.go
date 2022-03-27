package interfaces

type FixerFunc func(stacks []Stack)

func (f FixerFunc) Fix(stacks []Stack) {
	f(stacks)
}
