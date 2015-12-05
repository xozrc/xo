package randomutil_test

import (
	"fmt"
	"github.com/xo/common"
	"github.com/xo/randomutil"
	"testing"
)

var print = fmt.Print

type Hello struct {
}

func TestRandomNoRepeat(t *testing.T) {

	totalNum := 10
	tempWeights := make([]int, totalNum)
	for i := 0; i < totalNum; i++ {
		tempWeights[i] = i + 1
	}
	randNums := 10
	randItemList := randomutil.RandomListByWeights(tempWeights, randNums, false)

	for _, randItem := range randItemList {
		fmt.Println(randItem)
	}
	common.AssertTest(t, len(randItemList) == randNums, "random nums error")
}
