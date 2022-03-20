package guess

import (
	"sort"
	"sync"

	"github.com/xnslong/guess-stack/utils"
)

type CommonRootFixer struct {
	MinOverlaps int
}

const nonExistIndex = -1

type joint struct {
	CurrentIdx  int
	JoinPathIdx int
	JoinNodeIdx int
	Overlaps    int
	JoinGroup   int
}

type jointSlice []*joint

func (js jointSlice) Len() int {
	return len(js)
}

func (js jointSlice) Less(i, j int) bool {
	return js[i].Overlaps > js[j].Overlaps
}

func (js jointSlice) Swap(i, j int) {
	js[i], js[j] = js[j], js[i]
}

type computePath struct {
	Stack
	*joint
}

func (c *computePath) AddRoot(root []StackNode) {
	path := c.Stack.Path()
	newPath := make([]StackNode, len(root)+len(path))
	copy(newPath, root)
	copy(newPath[len(root):], path)
	c.Stack.SetPath(newPath)
}

func (c *CommonRootFixer) Fix(paths []Stack) {
	computeStacks := initComputePath(paths)

	computeJoints(computeStacks)

	joints := extractJoints(computeStacks)
	sort.Stable(jointSlice(joints))

	joinedJoints := joints[:0]

	for len(joints) > 0 {
		j0 := joints[0]

		if j0.JoinPathIdx == nonExistIndex {
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
		root := getRoot(targetStack, j0.JoinNodeIdx)

		// update stacks in the old group:
		// 1) add root
		// 2) update group to the new group
		updateRootForGroup(computeStacks, root, j0.JoinGroup, targetStack.JoinGroup)

		// if target joins back and makes a loop join, reset it and re-compute later.
		reComputeGroupJointWhenHasLoop(computeStacks, joints, j0.JoinGroup)

	}
}

func reComputeGroupJointWhenHasLoop(computePaths []*computePath, joints []*joint, toTest int) {
	groupJoint := computePaths[toTest].joint
	if groupJoint.JoinPathIdx != nonExistIndex {
		if computePaths[groupJoint.JoinPathIdx].JoinGroup == groupJoint.CurrentIdx {
			resetJoint(groupJoint)
			computeJoint(computePaths[groupJoint.CurrentIdx], computePaths)
			sort.Stable(jointSlice(joints))
		}
	}
}

func resetJoint(joint *joint) {
	joint.JoinGroup = joint.CurrentIdx
	joint.JoinPathIdx = nonExistIndex
	joint.JoinNodeIdx = 0
	joint.Overlaps = 0
}

func updateRootForGroup(computePaths []*computePath, root []StackNode, oldGroup int, newGroup int) {
	for _, stack := range computePaths {
		if stack.joint.JoinGroup == oldGroup {
			stack.AddRoot(root)
			stack.JoinGroup = newGroup
		}
	}
	return
}

func computeJoints(paths []*computePath) {
	p := utils.NewPool(100)

	wg := &sync.WaitGroup{}
	for _, path := range paths {
		wg.Add(1)
		p0 := path

		p.Submit(func() {
			computeJoint(p0, paths)
			wg.Done()
		})
	}
	p.Close()

	wg.Wait()
}

func computeJoint(path *computePath, stacks []*computePath) {
	if !path.NeedFix() {
		return
	}

	currentStack := path.Path()

	for _, stack := range stacks {
		if stack.JoinGroup == path.CurrentIdx {
			continue
		}
		if stack.Group() != path.Group() {
			continue
		}
		begin, length := maxOverlappingMiddleRange(currentStack, stack.Path())
		if length > path.joint.Overlaps {
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

func getRoot(node *computePath, leafCount int) []StackNode {
	stack := node.Stack.Path()
	i := len(stack) - leafCount
	return stack[:i]
}

func extractJoints(paths []*computePath) []*joint {
	joints := make([]*joint, len(paths))
	for i := range joints {
		joints[i] = paths[i].joint
	}
	return joints
}

func initComputePath(paths []Stack) []*computePath {
	result := make([]*computePath, 0, len(paths))

	for i, path := range paths {
		result = append(result, &computePath{
			Stack: path,
			joint: initialJointFor(i),
		})
	}

	return result
}

func initialJointFor(i int) *joint {
	return &joint{
		CurrentIdx:  i,
		JoinPathIdx: nonExistIndex,
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
