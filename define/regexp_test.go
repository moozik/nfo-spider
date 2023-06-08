package define

import (
	"fmt"
	"testing"
)

func TestRegexp(t *testing.T) {
	ret := PattrenAvFile.FindStringSubmatch("ABP-531.mp4")
	fmt.Println(ret)
}
