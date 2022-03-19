package core

import (
	"sort"
	"sync"
)

type CommonRootFixer struct {
	MinOverlaps int
}

const NonExistIndex = -1

type Joint struct {
	CurrentIdx  int
	JoinPathIdx int
	JoinNodeIdx int
	Overlaps    int
	JoinGroup   int
}

type JointSlice []*Joint

func (js JointSlice) Len() int {
	return len(js)
}

func (js JointSlice) Less(i, j int) bool {
	return js[i].Overlaps > js[j].Overlaps
}

func (js JointSlice) Swap(i, j int) {
	js[i], js[j] = js[j], js[i]
}

type ComputePath struct {
	Stack
	*Joint
}

func (c *ComputePath) AddRoot(root []StackNode) {
	path := c.Stack.Path()
	newPath := make([]StackNode, len(root)+len(path))
	copy(newPath, root)
	copy(newPath[len(root):], path)
	c.Stack.SetPath(newPath)
}

func (c *CommonRootFixer) Fix(paths []Stack) {
	computeStacks := InitComputePath(paths)

	ComputeJoints(computeStacks)

	joints := ExtractJoints(computeStacks)
	sort.Stable(JointSlice(joints))

	joinedJoints := joints[:0]

	for len(joints) > 0 {
		j0 := joints[0]

		if j0.JoinPathIdx == NonExistIndex {
			joints = joints[1:]
			joinedJoints = append(joinedJoints, j0)
			continue
		}

		if j0.Overlaps < c.MinOverlaps {
			joints = joints[1:]
			joinedJoints = append(joinedJoints, j0)
			continue
		}

		joints = joints[1:]
		joinedJoints = append(joinedJoints, j0)

		targetStack := computeStacks[j0.JoinPathIdx]
		root := GetRoot(targetStack, j0.JoinNodeIdx)

		// update stacks in the old group:
		// 1) add root
		// 2) update group to the new group
		UpdateRootForGroup(computeStacks, root, j0.JoinGroup, targetStack.JoinGroup)

		// if target joins back and makes a loop join, reset it and re-compute later.
		ReComputeGroupJointWhenHasLoop(computeStacks, joints, j0.JoinGroup)

	}
}

func ReComputeGroupJointWhenHasLoop(computePaths []*ComputePath, joints []*Joint, toTest int) {
	groupJoint := computePaths[toTest].Joint
	if groupJoint.JoinPathIdx != NonExistIndex {
		if computePaths[groupJoint.JoinPathIdx].JoinGroup == groupJoint.CurrentIdx {
			ResetJoint(groupJoint)
			ComputeJoint(computePaths[groupJoint.CurrentIdx], computePaths)
			sort.Stable(JointSlice(joints))
		}
	}
}

func ResetJoint(joint *Joint) {
	joint.JoinGroup = joint.CurrentIdx
	joint.JoinPathIdx = NonExistIndex
	joint.JoinNodeIdx = 0
	joint.Overlaps = 0
}

func UpdateRootForGroup(computePaths []*ComputePath, root []StackNode, oldGroup int, newGroup int) {
	for _, stack := range computePaths {
		if stack.Joint.JoinGroup == oldGroup {
			stack.AddRoot(root)
			stack.JoinGroup = newGroup
		}
	}
	return
}

type Pool struct {
	Concurrency int
	sync.WaitGroup
	taskChan chan func()
}

func NewPool(concurrency int) *Pool {
	taskChan := make(chan func(), 100)

	p := &Pool{
		Concurrency: concurrency,
		WaitGroup:   sync.WaitGroup{},
		taskChan:    taskChan,
	}
	for i := 0; i < concurrency; i++ {
		go p.work()
	}

	return p
}

func (p *Pool) work() {
	for {
		t, ok := <-p.taskChan
		if !ok {
			return
		}
		p.doTask(t)
	}
}

func (p *Pool) doTask(t func()) {
	p.WaitGroup.Done()
	t()
}

func (p *Pool) Submit(task func()) {
	p.WaitGroup.Add(1)
	p.taskChan <- task
}

func (p *Pool) StopSubmit() {
	close(p.taskChan)
}

func (p *Pool) WaitAll() {
	p.WaitGroup.Wait()
}

func ComputeJoints(paths []*ComputePath) {
	p := NewPool(100)

	for _, path := range paths {
		p0 := path
		p.Submit(func() {
			ComputeJoint(p0, paths)
		})
	}
	p.StopSubmit()
	p.WaitAll()
}

func ComputeJoint(path *ComputePath, stacks []*ComputePath) {
	if !path.NeedFix() {
		return
	}

	currentStack := path.Path()

	for _, stack := range stacks {
		if stack.JoinGroup == path.CurrentIdx {
			continue
		}
		begin, length := maxOverlappingMiddleRange(currentStack, stack.Path())
		if length > path.Joint.Overlaps {
			path.JoinPathIdx = stack.CurrentIdx
			path.JoinNodeIdx = begin
			path.Overlaps = length
		}
	}
	return
}

func maxOverlappingMiddleRange(p1, p2 []StackNode) (begin, length int) {
	if len(p2) <= 1 {
		return
	}

	for i, pl := 1, len(p2); i < pl; i++ {
		l := commonPrefixLen(p1, p2, i)
		if l > length {
			length = l
			begin = pl - i
		}
	}
	return
}

func GetRoot(node *ComputePath, leafCount int) []StackNode {
	stack := node.Stack.Path()
	i := len(stack) - leafCount
	return stack[:i]
}

func ExtractJoints(paths []*ComputePath) []*Joint {
	joints := make([]*Joint, len(paths))
	for i := range joints {
		joints[i] = paths[i].Joint
	}
	return joints
}

func InitComputePath(paths []Stack) []*ComputePath {
	result := make([]*ComputePath, 0, len(paths))

	for i, path := range paths {
		result = append(result, &ComputePath{
			Stack: path,
			Joint: initialJointFor(i),
		})
	}

	return result
}

func initialJointFor(i int) *Joint {
	return &Joint{
		CurrentIdx:  i,
		JoinPathIdx: NonExistIndex,
		JoinGroup:   i,
	}
}

func commonPrefixLen(a1, a2 []StackNode, s int) int {
	m := len(a1)
	l2 := len(a2) - s
	if m > l2 {
		m = l2
	}

	for j := 0; j < m; j++ {
		if !a1[j].EqualsTo(a2[j+s]) {
			return j
		}
	}

	return m
}
