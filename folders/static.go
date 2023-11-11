package folders

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"github.com/gofrs/uuid"
	"github.com/lucasepe/codename"
)

// These are all helper methods and fixed types.
// There's no real need for you to be editting these, but feel free to tweak it to suit your needs.
// If you do make changes here, be ready to discuss why these changes were made.

const dataSetSize = 1000
const DefaultOrgID = "c1556e17-b7c0-45a3-a6ae-9546248fb17a"

var _, basePath, _ = GetSampleData("sample.json")

func GenerateData() []*Folder {
	rng, _ := codename.DefaultRNG()
	sampleData := []*Folder{}

	for i := 1; i < dataSetSize; i++ {
		orgId := uuid.FromStringOrNil(DefaultOrgID)

		if i%3 == 0 {
			orgId = uuid.Must(uuid.NewV4())
		}

		deleted := rand.Int() % 2

		sampleData = append(sampleData, &Folder{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    codename.Generate(rng, 0),
			OrgId:   orgId,
			Deleted: deleted != 0,
		})
	}

	return sampleData
}

func PrettyPrint(b interface{}) {
	s, _ := json.MarshalIndent(b, "", "\t")
	fmt.Print(string(s))
	fmt.Print("\n")
}

func ChunkPrint(chunk *PaginationResult) {
	var res PrintFormat
	tmp := chunk.Folders
	for i := 0; i < len(tmp); i++ {
		res.Folders = append(res.Folders, *tmp[i])
	}

	fmt.Print("{data : \n")
	for i := 0; i < len(tmp); i++ {
		fmt.Println(res.Folders[i])
	}
	fmt.Print( ",\ntoken : ", (*chunk).Token, "\n}\n")
}

// Add error handling and filename
func GetSampleData(files string) ([]*Folder, string, error) {
	_, filename, _, _ := runtime.Caller(100)
	basePath := filepath.Dir(filename)
	filePath := filepath.Join(basePath, files)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, filename, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	jsonByte, err := io.ReadAll(file)
	// Read file error
	if err != nil {
		return nil, filename, fmt.Errorf("failed to read file: %w", err)
	}
	folders := []*Folder{}
	err = json.Unmarshal(jsonByte, &folders)
	// Parsing JSON failed
	if err != nil {
		return nil, filename, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return folders, filename, nil
}
