package persistence

import (
	"io"
	"os"
)

type Backend interface {
	PutBlob(string, *BlobData) error
	GetBlob(string) (*BlobData, error)

	// PutEnvironment(*Environment) error
	// GetEnvironment(int)  (*Environment, error)
}

// BlobData is the metadata and data associated with stored binary
// data. The fields and their usage varies depending on the operations,
// so please read the documentation for each field carefully.
type BlobData struct {
	// Key is the key for the blob data. This is populated on read and ignored
	// on any other operation.
	Key string

	// Data is the data for the blob data. When writing, this should be
	// the data to write. When reading, this is the data that is read.
	Data io.Reader

	closer io.Closer
}

func (d *BlobData) Close() error {
	if d.closer != nil {
		return d.closer.Close()
	}

	return nil
}

// WriteToFile is a helper to write BlobData to a file. While this is
// a very easy thing to do, it is so common that we provide a function
// for doing so.
func (d *BlobData) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, d.Data)
	return err
}
