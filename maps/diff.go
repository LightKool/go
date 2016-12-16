package maps

import (
	"fmt"
	"reflect"
)

type MapDifference interface {
	AreEqual() bool
	OnlyInLeft() map[interface{}]interface{}
	OnlyInRight() map[interface{}]interface{}
	InCommon() map[interface{}]interface{}
	Differences() map[interface{}]ValueDifference
}

type mapDiff struct {
	onlyInLeft  map[interface{}]interface{}
	onlyInRight map[interface{}]interface{}
	inCommon    map[interface{}]interface{}
	diff        map[interface{}]ValueDifference
}

func (d *mapDiff) AreEqual() bool {
	return len(d.onlyInLeft) == 0 && len(d.onlyInRight) == 0 && len(d.diff) == 0
}

func (d *mapDiff) OnlyInLeft() map[interface{}]interface{} {
	return d.onlyInLeft
}

func (d *mapDiff) OnlyInRight() map[interface{}]interface{} {
	return d.onlyInRight
}

func (d *mapDiff) InCommon() map[interface{}]interface{} {
	return d.inCommon
}

func (d *mapDiff) Differences() map[interface{}]ValueDifference {
	return d.diff
}

type ValueDifference interface {
	LeftValue() interface{}
	RightValue() interface{}
}

type valueDiff struct {
	leftValue  interface{}
	rightValue interface{}
}

func (d *valueDiff) LeftValue() interface{} {
	return d.leftValue
}

func (d *valueDiff) RightValue() interface{} {
	return d.leftValue
}

func (d *valueDiff) String() string {
	return "(" + fmt.Sprint(d.leftValue) + ", " + fmt.Sprint(d.rightValue) + ")"
}

func Difference(left map[interface{}]interface{}, right map[interface{}]interface{}) MapDifference {
	onlyInLeft := make(map[interface{}]interface{})
	onlyInRight := Copy(right)
	inCommon := make(map[interface{}]interface{})
	diff := make(map[interface{}]ValueDifference)

	for leftKey, leftValue := range left {
		if rightValue, ok := onlyInRight[leftKey]; ok {
			if reflect.DeepEqual(leftValue, rightValue) {
				inCommon[leftKey] = leftValue
			} else {
				diff[leftKey] = &valueDiff{leftValue, rightValue}
			}
			delete(onlyInRight, leftKey)
		} else {
			onlyInLeft[leftKey] = leftValue
		}
	}
	return &mapDiff{onlyInLeft, onlyInRight, inCommon, diff}
}

func Copy(src map[interface{}]interface{}) map[interface{}]interface{} {
	dst := make(map[interface{}]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
