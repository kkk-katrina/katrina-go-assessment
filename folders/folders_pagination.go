// Please go to annotation.md to view design explanation for Component2
package folders

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

func EncodedToken(orgID uuid.UUID) string {
	return base64.StdEncoding.EncodeToString([]byte(orgID.String()))
}

/*
	I assume a user story:
	When users login this website, I show all folders which is after pagination.
	Users should put new test file into folders container and name it as 'sample.json'.
*/

// work load is heavy, heap memory is not enough
func Pagination(limit int, inputfile string) error {
	folders, _, err := GetSampleData(inputfile)
	if err != nil {
		return nil
	}
	size := len(folders)

	//1. Create files
	ChunkFile, fileErr := os.Create("chunk.json")
	if fileErr != nil {
		fmt.Print("created failed")
		return errors.New("created failed")
	}
	HashFile, fileErr := os.Create("chunkHash.json")
	if fileErr != nil {
		return errors.New("created failed")
	}

	var ChunkData []PaginationResult
	Dict := []Pair{{Token: "", Index: 0}}
	// all folders -> {{folders,string}}
	count := 1
	for i := 0; i < size; i += limit {
		if i+limit < size {
			writeInJsonData := PaginationResult{
				Folders: folders[i : i+limit],
				Token:   EncodedToken(folders[i+limit-1].Id),
			}
			ChunkData = append(ChunkData, writeInJsonData)
			HashData := Pair{
				Token: EncodedToken(folders[i+limit-1].Id),
				Index: count,
			}
			count++
			Dict = append(Dict, HashData)
		} else {
			writeInJsonData := PaginationResult{
				Folders: folders[i:],
				Token:   EncodedToken(folders[size-1].Id),
			}
			ChunkData = append(ChunkData, writeInJsonData)
			HashData := Pair{
				Token: EncodedToken(folders[size-1].Id),
				Index: count,
			}
			count++
			Dict = append(Dict, HashData)
		}
	}

	ChunkData[len(ChunkData)-1].Token = ""
	Dict = Dict[:len(Dict)-1]

	jsonData, err := json.MarshalIndent(ChunkData, "", "    ")
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return errors.New("json is not correct")
	}
	_, err = ChunkFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return errors.New("json is not correct")
	}

	PairData, err1 := json.MarshalIndent(Dict, "", "    ")
	if err1 != nil {
		fmt.Println("Error encoding JSON:", err)
		return errors.New("pair is not correct")
	}
	_, err = HashFile.Write(PairData)
	if err != nil {
		fmt.Println("Error writing pair to file:", err)
		return errors.New("pair is not correct")
	}
	return err
}

func (r RealGetData) GetChunk(files string) ([]*PaginationResult, error) {
	filePath := filepath.Join(filepath.Dir(basePath), files)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print("Open file failed\n")
		return nil, fmt.Errorf("Open chunk file failed")
	}
	defer file.Close()

	jsonByte, err := io.ReadAll(file)
	// Read file error
	if err != nil {
		return nil, fmt.Errorf("Failed to read chunk file")
	}
	var folders []*PaginationResult
	err = json.Unmarshal(jsonByte, &folders)
	// Parsing JSON failed
	if err != nil {
		fmt.Print("Parsing JSON failed\n")
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return folders, nil
}

func (r RealGetData) GetPair(files string) ([]*Pair, error) {
	//_, filename, _, _ := runtime.Caller(0)

	filePath := filepath.Join(filepath.Dir(basePath), files)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print("Open file failed\n")
		return nil, errors.New("open file failed")
	}
	defer file.Close()

	jsonByte, err := io.ReadAll(file)
	// Read file error
	if err != nil {
		fmt.Print("Read file failed\n")
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var folders []*Pair
	err = json.Unmarshal(jsonByte, &folders)
	// Parsing JSON failed
	if err != nil {
		fmt.Print("Parsing JSON failed\n")
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return folders, nil
}

func Request(token string, gt GetData) (*PaginationResult, error) {
	// 1. open chunk file
	// 2. open hash file
	// 3. use token to find the index stored in hash file => get index
	// 4. use index to find the chunk data stored in chunk file

	ChunkData, err := gt.GetChunk("chunk.json")
	if err != nil {
		return nil, fmt.Errorf("Get chunk data failed\n")
	}

	PairData, err := gt.GetPair("chunkHash.json")
	if err != nil {
		return nil, fmt.Errorf("Get pair data failed\n")
	}

	if token == "" {
		return ChunkData[0], nil
	}

	index := -1
	PairSize := len(PairData)
	for i := 0; i < PairSize; i++ {
		if (*PairData[i]).Token == token {
			index = (*PairData[i]).Index
		}
	}
	if index == -1 {
		return nil, fmt.Errorf("Token does not exist\n")
	}
	res := ChunkData[index]
	return res, nil
}
