package folders

import (
	"github.com/gofrs/uuid"
)

type Folder struct {
	// An unique identifier for the folder, must be a valid UUID.
	// For example: '00001d65-d336-485a-8331-7b53f37e8f51'
	Id uuid.UUID `json:"id"`
	// Name associated with folder.
	Name string `json:"name"`
	// The organisation that the folder belongs to.
	OrgId uuid.UUID `json:"org_id"`
	// Whether a folder has been marked as deleted or not.
	Deleted bool `json:"deleted"`
}

type PrintFormat struct {
	Folders []Folder `json:"data"`
	Token   string   `json:"token"`
}

type Pair struct {
	Token string `json:"Token"`
	Index int    `json:"Index"`
}

type PaginationResult struct {
	Folders []*Folder `json:"folders"`
	Token   string    `json:"token"`
}

type FetchFolderRequest struct {
	OrgID     uuid.UUID
	InputFile string
}

type FetchFolderResponse struct {
	Folders []*Folder
}

type GetData interface {
	GetChunk(string) ([]*PaginationResult, error)
	GetPair(string) ([]*Pair, error)
}

type RealGetData struct{}
