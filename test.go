package main

import (
	"fmt"

	"github.com/rivo/uniseg"
)

func main() {
	str := "🤖⛏️🏰"
	str1 := "abc🤖⛏️🏰"
	fmt.Println(uniseg.GraphemeClusterCount(str))
	fmt.Println(uniseg.GraphemeClusterCount(str1))
}
