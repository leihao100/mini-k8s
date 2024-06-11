package main

import (
	"MiniK8S/utils/net"
	"fmt"
)

func main() {
	//net.AddPrometheus("work", "127.0.0.1")
	//net.AddPrometheus("work2", "127.0.0")
	//net.AddPrometheus("work","127.0.0.1")
	err := net.RemovePrometheus("work2")
	if err != nil {
		fmt.Println(err)
	}

}
