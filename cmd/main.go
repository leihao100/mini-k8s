package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	fmt.Println(cl.ImageList(context.Background(), types.ImageListOptions{}))

}
