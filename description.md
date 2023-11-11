# Project Description
This is the explanation document for this project
## Component 1
```go
// This function retrieves all folders associated with a specific OrgID
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	var (
		err error   //unused variables
		f1  Folder
		fs  []*Folder
	)
    // It's redundant that folders are converted from pointers to values and then back to pointers.
	f := []Folder{}
	r, _ := FetchAllFoldersByOrgID(req.OrgID) //It ignores error from FetchAllFoldersByOrgID
	for k, v := range r {
		f = append(f, *v)
	}
	var fp []*Folder
	for k1, v1 := range f {
		fp = append(fp, &v1)
	}
	var ffr *FetchFolderResponse
	ffr = &FetchFolderResponse{Folders: fp}
	return ffr, nil
}
// This function retrieves folders that belongs to a given OrgID from a sample file
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	folders := GetSampleData()

	resFolder := []*Folder{}
	for _, folder := range folders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}
```
The original implementation of these two functions overlooked potential errors and did not adequately address certain edge cases. To enhance their robustness and efficiency, I've refined the code by streamlining its structure and incorporating checks for edge cases. These improvements are also reflected in the test file, ensuring that the functions are now better equipped to handle a variety of scenarios with greater reliability.

## Component 2
"To implement a cursor based pagination, I followed the steps below:

1. Designed a pagination function to break the dataset into smaller chunks.
2. Used Base64 encoding on the ID of the last data item of each chunk to generate a token.
3. Saved these chunks along with their tokens in a new JSON file.
4. Stored the token with its index in a new hash file.
5. Implemented a request function that uses the token as a query parameter to return the corresponding data.

The rationale behind choosing to save chunks in a JSON file is that we only need to perform pagination once. Then, to retrieve a specific chunk, we can simply use the request function with the existing JSON file.
Additionally, I made some changes to static.go to accommodate different dataset files."