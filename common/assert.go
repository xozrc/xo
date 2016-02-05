package common

import (
	"testing"
)

func Assert(flag bool, desc string) {
	if !flag {
		panic(desc)
	}
}

func AssertTest(t *testing.T, flag bool, desc string) {
	if !flag {
		t.Error(desc)
	}
}
