package fix

import "sort"

type CommonRootFixer struct {
	CommonCount int
}

const NonExistIndex = -1

type Joint struct {
	CurrentIdx      int
	JoinPathIdx     int
	JoinNodeIdx     int
	CommonCount     int
	RootJoinPathIdx int
}

type JointSlice []*Joint

func (js JointSlice) Len() int {
	return len(js)
}

func (js JointSlice) Less(i, j int) bool {
	return js[i].CommonCount > js[j].CommonCount
}

func (js JointSlice) Swap(i, j int) {
	js[i], js[j] = js[j], js[i]
}

type Node struct {
	PathNode
	Parent *Node
}

func IsPointsTo(target, from *Node) bool {

	for from != nil {
		if from == target {
			return true
		}
		from = from.Parent
	}

	return false
}

func Path2Nodes(pn []PathNode) *Node {
	n := (*Node)(nil)

	for _, node := range pn {
		n = &Node{
			PathNode: node,
			Parent:   n,
		}
	}

	return n
}

type ComputePath struct {
	Path  Path
	Joint *Joint
	Node  *Node
}

func (c *ComputePath) ResetPath() {
	c.Path.SetPath(Nodes2Path(c.Node))
}

func maxOverlaps(n1, n2 *Node) (begin, length int) {
	p1 := Nodes2Path(n1)
	p2 := Nodes2Path(n2)

	begin, length = maxOverlapsForPathNodeArray(p1, p2)

	return
}

func maxOverlapsForPathNodeArray(p1, p2 []PathNode) (begin, length int) {
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

func computeJoint(p0 *ComputePath, all []*ComputePath, accept func(path *ComputePath) bool) {
	for _, path := range all {
		if p0.Joint.CurrentIdx != path.Joint.CurrentIdx && accept(path) {
			begin, length := maxOverlaps(p0.Node, path.Node)
			if length > p0.Joint.CommonCount {
				p0.Joint.CommonCount = length
				p0.Joint.JoinPathIdx = path.Joint.CurrentIdx
				p0.Joint.JoinNodeIdx = begin
			}
		}
	}

}

func (c *CommonRootFixer) Fix(paths []Path) {
	computePaths := InitComputePath(paths)

	ComputeJoints(computePaths)

	joints := ExtractJoints(computePaths)
	sort.Stable(JointSlice(joints))

	joinedJoints := joints[:0]

	for len(joints) > 0 {
		j0 := joints[0]

		if j0.JoinPathIdx == NonExistIndex {
			joints = joints[1:]
			joinedJoints = append(joinedJoints, j0)
			continue
		}

		if j0.CommonCount < c.CommonCount {
			joints = joints[1:]
			joinedJoints = append(joinedJoints, j0)
			continue
		}

		currentNode := computePaths[j0.CurrentIdx]
		currentRootNode := lastNode(currentNode.Node)
		jointNode := computePaths[j0.JoinPathIdx]

		if IsPointsTo(currentRootNode, jointNode.Node) {
			*j0 = *initialJointFor(j0.CurrentIdx)
			computeJoint(currentNode, computePaths, func(path *ComputePath) bool {
				return path.Joint.RootJoinPathIdx != jointNode.Joint.RootJoinPathIdx
			})
			sort.Stable(JointSlice(joints))
			continue
		}


		toJointNode := SkipNodes(jointNode.Node, j0.JoinNodeIdx)
		currentRootNode.Parent = toJointNode

		joints = joints[1:]
		joinedJoints = append(joinedJoints, j0)

		c := 0
		for _, joint := range joints {
			if joint.JoinPathIdx == j0.JoinPathIdx &&
				between(joint.JoinNodeIdx, j0.JoinNodeIdx, j0.JoinNodeIdx-j0.CommonCount) {

				lastNode := lastNode(computePaths[joint.CurrentIdx].Node)

				if IsPointsTo(currentRootNode, lastNode) {
					*joint = *initialJointFor(joint.CurrentIdx)
					unacceptable := make(map[int]bool)
					computeJoint(computePaths[joint.CurrentIdx], computePaths, func(path *ComputePath) bool {
						if unacceptable[path.Joint.RootJoinPathIdx] {
							return false
						}
						if IsPointsTo(path.Node, lastNode) {
							unacceptable[path.Joint.RootJoinPathIdx] = true
							return false
						}
						return true
					})
				} else {
					computeJoint(computePaths[joint.CurrentIdx], computePaths, func(path *ComputePath) bool {
						return path.Joint.RootJoinPathIdx == j0.CurrentIdx
					})
				}

				if joint.CurrentIdx == j0.CurrentIdx {
					c++
				}
			}
		}
		if c > 0 {
			sort.Stable(JointSlice(joints))
		}

		updateRootJoinPathIndex(j0, jointNode, joinedJoints)
	}

	for _, node := range computePaths {
		node.Path.SetPath(Nodes2Path(node.Node))
	}
}

func updateRootJoinPathIndex(original *Joint, target *ComputePath, toUpdate []*Joint) {
	for _, joint := range toUpdate {
		if joint.RootJoinPathIdx == original.CurrentIdx {
			joint.RootJoinPathIdx = target.Joint.RootJoinPathIdx
		}
	}
}

func InitComputePath(paths []Path) []*ComputePath {
	result := make([]*ComputePath, 0, len(paths))

	for i, path := range paths {
		result = append(result, &ComputePath{
			Path:  path,
			Joint: initialJointFor(i),
			Node:  Path2Nodes(path.Path()),
		})
	}

	return result
}

func SkipNodes(node *Node, toSkip int) *Node {
	for j := toSkip; j > 0; j-- {
		node = node.Parent
	}
	return node
}

func ComputeJoints(paths []*ComputePath) {
	for i, path := range paths {
		for j, another := range paths {
			if i != j {
				begin, length := maxOverlapsForPathNodeArray(path.Path.Path(), another.Path.Path())
				if length > path.Joint.CommonCount {
					path.Joint.JoinPathIdx = j
					path.Joint.JoinNodeIdx = begin
					path.Joint.CommonCount = length
				}
			}
		}
	}
}

func ExtractJoints(paths []*ComputePath) []*Joint {
	joints := make([]*Joint, len(paths))
	for i := range joints {
		joints[i] = paths[i].Joint
	}
	return joints
}

func initialJointFor(i int) *Joint {
	return &Joint{
		CurrentIdx:      i,
		JoinPathIdx:     NonExistIndex,
		RootJoinPathIdx: i,
	}
}

func lastNode(currentNode *Node) *Node {
	for currentNode != nil && currentNode.Parent != nil {
		currentNode = currentNode.Parent
	}
	return currentNode
}

func Nodes2Path(node *Node) []PathNode {
	dep := count(node)
	m := make([]PathNode, dep)
	for i := dep - 1; i >= 0; i-- {
		m[i] = node.PathNode
		node = node.Parent
	}
	return m
}

func count(node *Node) int {
	c := 0

	for node != nil {
		c++
		node = node.Parent
	}

	return c
}

func between(i, l, u int) bool {
	if i < l {
		return false
	}
	if i >= u {
		return false
	}
	return true
}

func commonPrefixLen(a1, a2 []PathNode) int {
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
