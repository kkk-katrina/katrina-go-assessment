package folders_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

type MockData struct{}

type MockDataWithError struct{}

func (m MockData) GetChunk(file string) ([]*folders.PaginationResult, error) {
	OrgID := uuid.FromStringOrNil(folders.DefaultOrgID)
	testFolder1 := []*folders.Folder{
		{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    "Active Folder",
			OrgId:   OrgID,
			Deleted: false,
		},
		{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    "Deleted Folder",
			OrgId:   OrgID,
			Deleted: true,
		},
	}
	testFolder2 := []*folders.Folder{
		{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    "Active Folder Two",
			OrgId:   uuid.Must(uuid.NewV4()),
			Deleted: false,
		},
		{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    "Deleted Folder Two",
			OrgId:   uuid.Must(uuid.NewV4()),
			Deleted: true,
		},
	}
	testFolder3 := []*folders.Folder{
		{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    "Active Folder Three",
			OrgId:   OrgID,
			Deleted: false,
		},
	}

	res := []*folders.PaginationResult{
		{
			Folders: testFolder3,
			Token:   "",
		},
		{
			Folders: testFolder1,
			Token:   "token1",
		},
		{
			Folders: testFolder2,
			Token:   "token2",
		},
	}
	return res, nil
}

func (m MockData) GetPair(file string) ([]*folders.Pair, error) {
	pairs := []*folders.Pair{
		{
			Token: "",
			Index: 0,
		},
		{
			Token: "token1",
			Index: 1,
		},
		{
			Token: "token2",
			Index: 2,
		},
	}
	return pairs, nil
}

func (m MockDataWithError) GetChunk(file string) ([]*folders.PaginationResult, error) {
	return nil, errors.New("Open Chunk file error")
}

func (m MockDataWithError) GetPair(file string) ([]*folders.Pair, error) {
	return nil, errors.New("Open Hash file error")
}

func Test_Request(t *testing.T) {
	mockData := MockData{}
	// Check if valid token returns correct result
	t.Run("Request chunk with valid token", func(t *testing.T) {
		res, err := folders.Request("token1", mockData)
		assert.NoError(t, err)
		assert.Equal(t, len(res.Folders), 2)
	})
	// Check if empty token returns first chunk
	t.Run("Request chunk with empty token", func(t *testing.T) {
		res, err := folders.Request("", mockData)
		assert.NoError(t, err)
		assert.Equal(t, len(res.Folders), 1)
	})
	// Check if invalid token returns error
	t.Run("Request chunk with invalid token", func(t *testing.T) {
		res, err := folders.Request("valid", mockData)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	// Check if error from GetChunk or GetPair returns internal error
	t.Run("Request chunk with error from GetChunk or GetPair", func(t *testing.T) {
		mockDataWithError := MockDataWithError{}
		res, err := folders.Request("token", mockDataWithError)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
