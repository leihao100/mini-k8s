package main

import (
	"fmt"
	"github.com/docker/docker/testutil"
)

func main() {
	fmt.Println(testutil.GenerateRandomAlphaOnlyString(9))
	//t := time.Now().Format(time.RFC3339Nano)
	//fmt.Println(t)
	////layout := "2024-05-27 11:22:26.2743388"
	//nt, err := time.Parse(time.RFC3339Nano, t)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//time.Sleep(1 * time.Second)
	//fmt.Println(nt)
	//fmt.Println(time.Since(nt))
}
