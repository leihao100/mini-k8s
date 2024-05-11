package kubectl

import (
	"fmt"

	"github.com/spf13/cobra"
)

func create(cmd *cobra.Command, args []string) {
	resource := args[0]
	objType, err := parseResourceType(resource)
	if err != nil {
		fmt.Printf("Error:%s", err)
		return
	}
	filename, err := getFileName()
	if err != nil {
		fmt.Printf("Error:%s", err)
		return
	}
	//TODO
	fmt.Println(err)
	fmt.Println(objType)
	fmt.Println(filename)
}

func delete(cmd *cobra.Command, args []string) {
	resource := args[0]
	name := args[1]

	objType, err := parseResourceType(resource)
	if err != nil {
		fmt.Printf("Error:%s", err)
		return
	}
	//TODO
	fmt.Println(name)
	fmt.Println(objType)
}

func get(cmd *cobra.Command, args []string) {

}

func describe(cmd *cobra.Command, args []string) {

}

func clear(cmd *cobra.Command, args []string) {

}
