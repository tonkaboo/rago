package repo

import (
	"fmt"
	"gems/pkg/chunk"
)

type Repo interface {
	CheckIndex() bool
	InsertChunk(*chunk.Chunk) error
	RetrieveRelatedChunks([]float32) ([]string, error)
	CloseConn()
}

func NewRepo(repoType string, fileName string) (Repo, error) {
	if repoType == "MatrixOne" {
		repo, err := NewMatrixOne(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create a new MatrixOne repo: %s", err)
		}
		return repo, nil
	} else if repoType == "FAISS" {
		repo, err := NewFaiss()
		if err != nil {
			return nil, fmt.Errorf("failed to create a new FAISS repo: %s", err)
		}
		return repo, nil
	} else {
		return nil, fmt.Errorf("unknown repo type: %s", repoType)
	}
}
