package trie_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xnslong/guess-stack/core/utils/trie"
)

func TestNewTrie(t *testing.T) {
	var lists = [][]interface{}{
		{1, 2, 3, 4},
		{1, 2, 3, 4, 5},
		{2, 3, 5},
	}

	pt := trie.NewTrie()
	for i, list := range lists {
		pt.AddPath(list, i)
	}

	t.Run("visit_all", func(t *testing.T) {
		vs := make([][]interface{}, 0)
		pt.VisitAllPath(func(path []trie.Node, target interface{}) {
			require.IsType(t, 1, target)
			idx := target.(int)
			if len(vs) <= idx {
				t := make([][]interface{}, idx+1)
				copy(t, vs)
				vs = t
			}

			v := []interface{}{}
			for _, node := range path {
				v = append(v, node.Element())
			}
			vs[idx] = v
		})
		assert.Equal(t, lists, vs)
	})

	t.Run("common_path", func(t *testing.T) {
		path := pt.MaxCommonPath([]interface{}{1, 2, 4, 5})
		assert.Equal(t, []interface{}{1, 2}, elements(path))

		path2 := pt.MaxCommonPath([]interface{}{1, 2, 3, 5})
		assert.Equal(t, []interface{}{1, 2, 3}, elements(path2))
	})

	t.Run("attach", func(t *testing.T) {
		path2 := pt.MaxCommonPath([]interface{}{1, 2, 3, 4})
		for _, node := range path2 {
			assert.Nil(t, node.Attachment())
		}

		path := pt.MaxCommonPath([]interface{}{1, 2, 3, 4, 5})
		for i, node := range path {
			node.Attach(i)
		}

		path3 := pt.MaxCommonPath([]interface{}{1, 2, 3, 4})
		for i, node := range path3 {
			assert.Equal(t, i, node.Attachment())
		}
	})
}

func elements(path []trie.Node) []interface{} {
	list := make([]interface{}, 0, len(path))
	for _, node := range path {
		list = append(list, node.Element())
	}
	return list
}

func assertSameSet(t *testing.T, s1, s2 [][]interface{}) {
	assert.Equal(t, len(s1), len(s2), "length: s1=%d, sw=%d", len(s1), len(s2))
	m1 := toMap(s1)
	m2 := toMap(s2)
	assert.Equal(t, m1, m2, "should have should elements")
}

func toMap(s [][]interface{}) map[string]struct{} {
	m := make(map[string]struct{}, len(s))
	for _, e := range s {
		se := fmt.Sprint(e)
		m[se] = struct{}{}
	}
	return m
}
