package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	// "github.com/gofrs/uuid"
)

func main() {
	/*
		component 1:
	*/
	inputFileName := "./folders/sample.json"
	fmt.Print("component 1:\n")
	req := &folders.FetchFolderRequest{
		OrgID:     uuid.FromStringOrNil(folders.DefaultOrgID),
		InputFile: inputFileName,
	}
	res, err := folders.GetAllFolders(req)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	folders.PrettyPrint(res)
	fmt.Println("Component 1 is done!")
	fmt.Println()
	/*
		component 2:
		pagination, output to new json, request
		token -> next chunk
		1. pagination(): all folders-> chunk -> output a new json file
		2. init token variable: varible new json file any element token
		3. request(token): next chunk(which was stored in new json file)
	*/
	fmt.Print("Component 2 Start:\n")

	folders.Pagination(50, inputFileName)
	fmt.Print("Pagination is done!\n")

	gt := folders.RealGetData{}
	a, err := folders.Request("", gt)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	folders.ChunkPrint(a)
	fmt.Println("Component 2 is done!")
}
