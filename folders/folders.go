// Please go to annotation.md to view detailed annotation for each function
package folders

import (
	"errors"

	"github.com/gofrs/uuid"
)

var FetchAllFoldersByOrgIDFunc = FetchAllFoldersByOrgID

// GetAllFolders fetches all folders associated with a given orgID
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	// Check if orgID is valid
	if req.OrgID == uuid.Nil {
		return nil, errors.New("organization ID cannot be empty")
	}
	folders, err := FetchAllFoldersByOrgIDFunc(req.InputFile, req.OrgID)
	if err != nil {
		return nil, err
	}
	response := &FetchFolderResponse{Folders: folders}
	return response, nil
}

func FetchAllFoldersByOrgID(filename string, orgID uuid.UUID) ([]*Folder, error) {
	// Check if orgID is valid
	if orgID == uuid.Nil {

		return nil, errors.New("organization ID is invalid")
	}

	folders, _, err := GetSampleData(filename)
	// Check if error from GetSampleData
	if err != nil {
		return nil, err
	}

	resFolder := []*Folder{}
	for _, folder := range folders {
		// Check if folder is associated with orgID and not deleted
		if folder.OrgId == orgID && !folder.Deleted {
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}
