package folders_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllFolders(t *testing.T) {
	// use mockFetchAllFoldersByOrgID for testing
	originFunc := folders.FetchAllFoldersByOrgIDFunc
	defer func() { folders.FetchAllFoldersByOrgIDFunc = originFunc }()
	// mock data
	validOrgID := uuid.FromStringOrNil(folders.DefaultOrgID)
	emptyOrgID := uuid.Nil
	defaultPath := ".folders/sample.json"
	mockFolders := []*folders.Folder{
		{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    "Active Folder",
			OrgId:   validOrgID,
			Deleted: false,
		},
		{
			Id:      uuid.Must(uuid.NewV4()),
			Name:    "Deleted Folder",
			OrgId:   validOrgID,
			Deleted: true,
		},
	}
	folders.FetchAllFoldersByOrgIDFunc = func(filename string, orgID uuid.UUID) ([]*folders.Folder, error) {
		var resFolders []*folders.Folder
		for _, folder := range mockFolders {
			if folder.OrgId == orgID && !folder.Deleted {
				resFolders = append(resFolders, folder)
			}
		}
		return resFolders, nil
	}

	// Check if valid orgID returns correct folders
	t.Run("Get folders with valid OrgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID:     validOrgID,
			InputFile: defaultPath,
		}
		res, err := folders.GetAllFolders(req)
		assert.NoError(t, err)
		assert.Equal(t, len(res.Folders), 1)
	})

	// Check if empty orgID returns error
	t.Run("Get folders with empty OrgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID:     emptyOrgID,
			InputFile: defaultPath,
		}
		res, err := folders.GetAllFolders(req)
		assert.Error(t, err) // Expect error
		assert.Nil(t, res)
	})

	// Check if invalid orgID returns empty folders
	t.Run("Get folders with invalid OrgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID:     uuid.Must(uuid.NewV4()),
			InputFile: defaultPath,
		}
		res, err := folders.GetAllFolders(req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res.Folders)
	})

	// Test if deleted folders are returned
	t.Run("Handle deleted folders", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID:     validOrgID,
			InputFile: defaultPath,
		}
		res, err := folders.GetAllFolders(req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		for _, folder := range res.Folders {
			assert.False(t, folder.Deleted)
		}
	})

	// Test error from FetchAllFoldersByOrgID
	t.Run("Error from FetchAllFoldersByOrgID", func(t *testing.T) {
		folders.FetchAllFoldersByOrgIDFunc = func(filename string, orgID uuid.UUID) ([]*folders.Folder, error) {
			return nil, errors.New("Internal error")
		}
		req := &folders.FetchFolderRequest{
			OrgID:     validOrgID,
			InputFile: defaultPath,
		}
		res, err := folders.GetAllFolders(req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
