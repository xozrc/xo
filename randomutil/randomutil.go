package randomutil

import (
	"fmt"
	"github.com/xo/common"
	"math/rand"
	"time"
)

var print = fmt.Print

func RandomListByWeights(weights []int, nums int, repeat bool) (randomItemList []int) {
	common.Assert(nums >= 1, "random num  must no less than 1")

	if !repeat {
		common.Assert(nums <= len(weights), "random nums should no more than item list length")
	}

	randomItemList = make([]int, 0)

	tempWeights := make([]int, len(weights))

	totalWeights := 0
	for i, weight := range weights {
		common.Assert(weight >= 0, "weight should no less than 0")
		totalWeights += weight
		tempWeights[i] = weight
	}

	for i := 0; i < nums; i++ {
		common.Assert(totalWeights > 0, "total weights should be more than 0")
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(totalWeights)
		for weightIndex, weight := range tempWeights {
			if randNum >= weight {
				randNum -= weight
				continue
			}
			randomItemList = append(randomItemList, weightIndex)
			if !repeat {
				totalWeights -= weight
				tempWeights[weightIndex] = 0
			}
			break
		}
	}
	return
}

func RandBetween(min int, max int) (num int) {
	common.Assert(min < max, "min should less than max")
	diff := max - min
	num = rand.Intn(diff)
	return
}
