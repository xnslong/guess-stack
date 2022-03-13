package fix

import "sort"

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
	NeedFix     bool
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

func maxOverlapsForPathNodeArray(p1, p2 []StackNode) (begin, length int) {
	if len(p2) == 0 {
		return
	}

	p2 = p2[1:]
	for len(p2) > 0 {
		l := commonPrefixLen(p1, p2)
		if l > length {
			length = l
			begin = len(p2)
		}
		p2 = p2[1:]
	}
	return
}

func (c *CommonRootFixer) Fix(paths []Stack, needFix []bool) {
	computeStacks := InitComputePath(paths, needFix)

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
		joinedStacks := UpdateRootForGroup(computeStacks, root, j0.JoinGroup, targetStack.JoinGroup)

		// if target joins back and makes a loop join, reset it and re-compute later.
		ReComputeGroupJointWhenHasLoop(computeStacks, joinedStacks, joints, j0.JoinGroup)

	}
}

func ReComputeGroupJointWhenHasLoop(computePaths []*ComputePath, stacks []*ComputePath, joints []*Joint, toTest int) {
	groupJoint := computePaths[toTest].Joint
	if groupJoint.JoinPathIdx != NonExistIndex {
		if computePaths[groupJoint.JoinPathIdx].JoinGroup == groupJoint.CurrentIdx {
			ResetJoint(groupJoint)
			ReComputeJoint(computePaths[groupJoint.CurrentIdx], computePaths)
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

func UpdateRootForGroup(computePaths []*ComputePath, root []StackNode, oldGroup int, newGroup int) []*ComputePath {
	var joinedStacks []*ComputePath
	for _, stack := range computePaths {
		if stack.Joint.JoinGroup == oldGroup {
			stack.AddRoot(root)
			joinedStacks = append(joinedStacks, stack)
			stack.JoinGroup = newGroup
		}
	}
	return joinedStacks
}

func ReComputeJoint(path *ComputePath, stacks []*ComputePath) bool {
	if !path.NeedFix {
		return false
	}

	currentStack := path.Path()

	changed := false
	for _, stack := range stacks {
		if stack.JoinGroup == path.CurrentIdx {
			continue
		}
		begin, length := maxOverlapsForPathNodeArray(currentStack, stack.Path())
		if length > path.Joint.Overlaps {
			path.JoinPathIdx = stack.CurrentIdx
			path.JoinNodeIdx = begin
			path.Overlaps = length
			changed = true
		}
	}
	return changed
}

func ComputeJoints(paths []*ComputePath) {
	for i, path := range paths {
		if !path.NeedFix {
			continue
		}
		for j, another := range paths {
			if i != j {
				begin, length := maxOverlapsForPathNodeArray(path.Stack.Path(), another.Stack.Path())
				if length > path.Joint.Overlaps {
					path.Joint.JoinPathIdx = j
					path.Joint.JoinNodeIdx = begin
					path.Joint.Overlaps = length
				}
			}
		}
	}
}

func GetRoot(node *ComputePath, leafCount int) []StackNode {
	stack := node.Stack.Path()
	i := len(stack) - leafCount
	return stack[:i]
}

func InitComputePath(paths []Stack, needFix []bool) []*ComputePath {
	result := make([]*ComputePath, 0, len(paths))

	for i, path := range paths {
		result = append(result, &ComputePath{
			Stack: path,
			Joint: initialJointFor(i, needFix[i]),
		})
	}

	return result
}

func ExtractJoints(paths []*ComputePath) []*Joint {
	joints := make([]*Joint, len(paths))
	for i := range joints {
		joints[i] = paths[i].Joint
	}
	return joints
}

func initialJointFor(i int, needFix bool) *Joint {
	return &Joint{
		CurrentIdx:  i,
		JoinPathIdx: NonExistIndex,
		JoinGroup:   i,
		NeedFix:     needFix,
	}
}

func commonPrefixLen(a1, a2 []StackNode) int {
	m := len(a1)
	if m > len(a2) {
		m = len(a2)
	}

	for i := 0; i < m; i++ {
		if !a1[i].EqualsTo(a2[i]) {
			return i
		}
	}

	return m
}
