package fix

import "sort"

type CommonRootFixer struct {
	CommonCount int
}

const NonExistIndex = -1

type Joint struct {
	CurrentIdx  int
	JoinPathIdx int
	JoinNodeIdx int
	CommonCount int
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

func (c *CommonRootFixer) Fix(paths []Path) {
	joints := InitJoints(paths)

	ComputeJoints(paths, joints)
	nodes := PathToNodes(paths)

	natureJoints := make([]*Joint, len(paths))
	copy(natureJoints, joints)

	sort.Stable(JointSlice(joints))

	for len(joints) > 0 {
		j0 := joints[0]
		joints = joints[1:]

		currentNode := nodes[j0.CurrentIdx]

		if j0.JoinPathIdx == NonExistIndex {
			continue
		}

		if j0.CommonCount < c.CommonCount {
			continue
		}

		jointNode := nodes[j0.JoinPathIdx]
		currentRootNode := lastNode(currentNode)

		if IsPointsTo(currentRootNode, jointNode) {
			// TODO
			*j0 = *initialJointFor(j0.CurrentIdx)
			continue
		}

		toSkip := j0.JoinNodeIdx
		jointNode = SkipNodes(jointNode, toSkip)
		currentRootNode.Parent = jointNode
		paths[j0.CurrentIdx].SetPath(Nodes2Path(currentNode))

		c := 0
		for _, joint := range joints {
			if joint.JoinPathIdx == j0.JoinPathIdx && between(joint.JoinNodeIdx, j0.JoinNodeIdx, j0.JoinNodeIdx+j0.CommonCount) {
				compare(joint, paths[joint.CurrentIdx], j0.CurrentIdx, paths[j0.CurrentIdx])
				if joint.CurrentIdx == j0.CurrentIdx {
					c++
				}
			}
		}

		if c > 0 {
			sort.Stable(JointSlice(joints))
		}
	}

	for i, node := range nodes {
		paths[i].SetPath(Nodes2Path(node))
	}
}

func SkipNodes(node *Node, toSkip int) *Node {
	for j := toSkip; j > 0; j-- {
		node = node.Parent
	}
	return node
}

func PathToNodes(paths []Path) []*Node {
	nodes := make([]*Node, len(paths))
	for i, path := range paths {
		nodes[i] = Path2Nodes(path.Path())
	}
	return nodes
}

func ComputeJoints(paths []Path, joints []*Joint) {
	for i, path := range paths {
		for j, another := range paths {
			if i != j {
				compare(joints[i], path, j, another)
			}
		}
	}
}

func InitJoints(paths []Path) []*Joint {
	joints := make([]*Joint, len(paths))
	for i := range joints {
		joints[i] = initialJointFor(i)
	}
	return joints
}

func initialJointFor(i int) *Joint {
	return &Joint{
		CurrentIdx:  i,
		JoinPathIdx: NonExistIndex,
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

func compare(j *Joint, path Path, i int, another Path) {
	current := path.Path()
	nodes := another.Path()

	c := 0
	for len(nodes) > 0 {
		l := commonPrefixLen(current, nodes)
		if l > j.CommonCount {
			j.JoinPathIdx = i
			j.JoinNodeIdx = len(nodes)
			j.CommonCount = l
		}
		nodes = nodes[1:]
		c++
	}
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
