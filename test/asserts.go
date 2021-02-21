package test

import (
	"math"
	"reflect"
	"testing"
)

func AssertDeepEqual(x, y interface{}, t *testing.T) {
	if !reflect.DeepEqual(x, y) {
		t.Errorf("\ngot %s\nexp %s", x, y)
	}
}

func AssertDeepEqualWithMessage(x, y interface{}, message string, t *testing.T) {
	if !reflect.DeepEqual(x, y) {
		t.Errorf("%s\ngot %s\nexp %s", message, x, y)
	}
}

func AssertDeepNotEqual(x, y interface{}, t *testing.T) {
	if reflect.DeepEqual(x, y) {
		t.Errorf("\ngot %s\nexp %s", x, y)
	}
}

func AssertAlmostEqual(a float64, b float64, places int, t *testing.T) {
	if !(math.Abs(a-b) <= math.Pow10(-places)) {
		t.Errorf("\ngot %f\nexp %f", a, b)
	}
}
