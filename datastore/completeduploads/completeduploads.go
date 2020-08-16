package completeduploads

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/utils/filesystem"
)

// Service represents the repository where uploaded objects are tracked
type Service struct {
	repo Repository
}

// NewService created a Service to track uploaded objects
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// Close closes the service.
//
// No operation could be done after that.
func (s *Service) Close() error {
	return s.repo.Close()
}

// IsAlreadyUploaded checks if the file was already uploaded
func (s *Service) IsAlreadyUploaded(filePath string) (bool, error) {
	// find a previous upload in the repository
	item, err := s.repo.Get(filePath)
	if err != nil {
		// this file was not uploaded before
		return false, nil
	}

	// check stored last modified time with the current one
	// to see if the file has been modified
	fileMtime, err := filesystem.GetMTime(filePath)
	if err != nil {
		return false, err
	}
	if fileMtime.Unix() == item.modifyTime {
		return true, nil
	}

	// file was not uploaded before or modified time has changed after being
	// uploaded
	fileHash, err := Hash(filePath)
	if err != nil {
		return false, err
	}

	// checks if the file is the same (equal hash)
	if item.hash == fileHash {
		// update last modified time on the cache
		err = s.CacheAsAlreadyUploaded(filePath)
		if err != nil {
			return true, err
		}
	}

	return false, nil
}

// CacheAsAlreadyUploaded marks a file as already uploaded to prevent re-uploads
func (s *Service) CacheAsAlreadyUploaded(filePath string) error {
	item, err := NewCompletedUploadedFileItem(filePath)
	if err != nil {
		return err
	}
	return s.repo.Put(item)
}

// RemoveAsAlreadyUploaded removes a file previously marked as uploaded
func (s *Service) RemoveAsAlreadyUploaded(filePath string) error {
	return s.repo.Delete(filePath)
}
