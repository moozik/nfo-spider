package utils

import (
	"fmt"
	"testing"
)

func TestCache_Set(t *testing.T) {
	c := NewCache("sensen", "av")
	fmt.Println(c.Set("wang", []byte("aoyun2")))
}
