package timeutil_test

import (
	"fmt"
	"github.com/xo/common"
	"github.com/xo/timeutil"
	"testing"
	"time"
)

var print = fmt.Print

func TestNow(t *testing.T) {
	common.AssertTest(t, timeutil.Now() > 0, "time should more than 0 ")
}
