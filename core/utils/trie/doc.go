// Package trie provides the Trie (prefix tree) ability,
// We can add a series of paths to the Trie, where paths with the same prefix will be shared.
//
// Suppose we define the following method to operate trie.Node.
//     type extractorFunc func (n trie.Node) interface{}
//     func forElement(n trie.Node) interface{} {  return n.Element()  }
//     func forAttachment(n trie.Node) interface{} {  return n.Attachment()  }
//     func extract(n []trie.Node, extractor extractorFunc) interface{}{
//         result := make([]interface{}, len(n))
//         for i, e := range n {
//             result[i] = extractor(e)
//         }
//         return result
//     }
//
// Then for the following example
//     t := trie.NewTrie()
//     // build a trie by adding path to it.
//     t.AddPath([]interface{}{1, 2, 3, 4}, 1)
//     t.AddPath([]interface{}{1, 2, 6, 7}, 2)
//     t.AddPath([]interface{}{3, 4, 5},    3)
//     t.AddPath([]interface{}{3, 4, 5, 6}, 4)
//
//     // return a prefix path in the trie where the given path is supposed to stay.
//     r := t.PrefixFor([]interface{}{1, 2, 3, 5})
//     assertEquals([]interface{}{1, 2, 3}, extract(r, forElement))
//
//     // set attachment for nodes in the trie, it will reflect to all those paths sharing the prefix nodes.
//     assertEquals([]interface{}{nil, nil, nil}, extract(r, forAttachment))
//     r[0].Attach("a")
//     r[1].Attach("b")
//     r[2].Attach("c")
//     assertEquals([]interface{}{"a", "b", "c"}, extract(r, forAttachment))
//
//     r1 := t.PrefixFor([]interface{}{1, 2, 6, 7})
//     assertEquals([]interface{}{1, 2, 6, 7}, extract(r, forElement))
//     assertEquals([]interface{}{"a", "b", nil, nil}, extract(r, forAttachment))
//
//     t.VisitAllPath(func(path []trie.Node, target interface{}) {
//         fmt.Println(extract(path, forElement), extract(path, forAttachment), target)
//         // will get the following results in arbitrary order
//         // [1, 2, 3, 4] [a b c <nil>] 1
//         // [1, 2, 6, 7] [a b <nil> <nil>] 2
//         // [3, 4, 5] [<nil> <nil> <nil>] 3
//         // [3, 4, 5, 6] [<nil> <nil> <nil> <nil>] 4
//     })
package trie
