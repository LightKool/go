package maps

import (
	"testing"
)

func TestDifference(t *testing.T) {
	m1 := map[interface{}]interface{}{"key1": "value1", "onlyIn1": "valueOnlyIn1", "diff": "diff1"}
	m2 := map[interface{}]interface{}{"key1": "value1", "onlyIn2": "valueOnlyIn2", "diff": "diff2"}
	diff := Difference(m1, m2)
	t.Log(diff.AreEqual())
	t.Log(diff.OnlyInLeft())
	t.Log(diff.OnlyInRight())
	t.Log(diff.InCommon())
	t.Log(diff.Differences())
	t.Log(m1)
	t.Log(m2)
}
