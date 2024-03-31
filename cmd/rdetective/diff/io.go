package diff

import (
	"io"
	"os"
)

type FileSource struct {
	fileName string
}

func (f FileSource) GetReader() (io.Reader, error) {
	file, err := os.Open(f.fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}
