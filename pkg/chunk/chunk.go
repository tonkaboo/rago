package chunk

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Chunk struct {
	Content   string
	Embedding []float32
	id        int
}

func SplitDoc(docName string, chunkSize int) ([]*Chunk, error) {
	fileDoc, err := os.Open(fmt.Sprintf("../docs/%s", docName))
	if err != nil {
		return nil, fmt.Errorf("an error occured when opening the document: %s", err)
	}
	defer fileDoc.Close()

	byteDoc, err := io.ReadAll(fileDoc)
	if err != nil {
		return nil, fmt.Errorf("an error occured when reading the document: %s", err)
	}

	strDoc := string(byteDoc)
	sentences := strings.Split(strDoc, ".")
	chunks := make([]*Chunk, 0)
	var currentChunk string
	var count int

	chunkId := 0
	for _, sentence := range sentences {
		currentChunk += sentence + "."
		count++
		if count >= chunkSize {
			chunks = append(chunks, &Chunk{Content: currentChunk, id: chunkId})
			currentChunk = ""
			count = 0
			chunkId++
		}
	}

	// If there is any remaining content, append it to the current chunk
	if currentChunk != "" {
		chunks = append(chunks, &Chunk{Content: currentChunk, id: chunkId})
	}

	return chunks, nil
}
