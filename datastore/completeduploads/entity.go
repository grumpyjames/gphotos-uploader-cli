package completeduploads

import (
	"fmt"
	"io"
	"os"

	"github.com/pierrec/xxHash/xxHash32"

	"github.com/gphotosuploader/gphotos-uploader-cli/utils/filesystem"
)

var (
	// ErrNotFound not found
	ErrNotFound = fmt.Errorf("not found")

	// ErrCannotBeDeleted bookmark cannot be deleted
	ErrCannotBeDeleted = fmt.Errorf("cannot be deleted")

	// ErrCannotGetMTime
	ErrCannotGetMTime = fmt.Errorf("failed getting local image mtime")
)

type CompletedUploadedFileItem struct {
	path       string
	hash       uint32
	modifyTime int64
}

// NewCompletedUploadedFileItem creates a new item for the specified file
func NewCompletedUploadedFileItem(filePath string) (CompletedUploadedFileItem, error) {
	item := CompletedUploadedFileItem{
		path: filePath,
	}

	fileHash, err := Hash(filePath)
	if err != nil {
		return item, err
	}
	item.hash = fileHash

	mTime, err := filesystem.GetMTime(filePath)
	if err != nil {
		return item, ErrCannotGetMTime
	}

	item.modifyTime = mTime.Unix()

	return item, nil
}

// Hash return the hash of a file
func Hash(filePath string) (uint32, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer inputFile.Close()

	hasher := xxHash32.New(0xCAFE) // value.Hash32
	defer hasher.Reset()

	_, err = io.Copy(hasher, inputFile)
	if err != nil {
		return 0, err
	}

	return hasher.Sum32(), nil
}
