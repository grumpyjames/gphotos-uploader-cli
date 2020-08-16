package upload

import (
	"context"
	"fmt"
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

type EnqueuedJob struct {
	Context       context.Context
	PhotosService gPhotosService
	FileTracker   FileTracker
	Logger        log.Logger

	Path            string
	AlbumName       string
	DeleteOnSuccess bool
}

func uploadFile(ctx context.Context, path string, photos gPhotosService) (string, error) {
	panic("Not implemented; should be a method on PhotosService")
}

func batchAddMediaItems(ctx context.Context, uploadTokens []string, photos gPhotosService, albumId string) (string, error) {
	panic("Not implemented; should be a method on PhotosService")
}

func (job *EnqueuedJob) Process() error {
	// Get or create the album
	albumId, err := job.albumID()
	if err != nil {
		return err
	}

	token, err := uploadFile(job.Context, job.Path, job.PhotosService)
	if err != nil {
		return err
	}

	err = job.FileTracker.CacheAsAlreadyUploaded(job.Path, token, job.AlbumName)
	if err != nil {
		job.Logger.Warnf("Tracking file as uploaded failed: file=%s, error=%v", job.Path, err)
	}

	// Upload the file and add it to PhotosService.
	_, err = batchAddMediaItems(job.Context, []string{token}, job.PhotosService, albumId)
	if err != nil {
		return err
	}

	// Mark the file as uploaded in the FileTracker.
	err = job.FileTracker.MediaItemCreated(job.Path)
	if err != nil {
		job.Logger.Warnf("Tracking file as uploaded failed: file=%s, error=%v", job.Path, err)
	}

	// If was requested, remove the file after being uploaded.
	if job.DeleteOnSuccess {
		if err := os.Remove(job.Path); err != nil {
			job.Logger.Errorf("Deletion request failed: file=%s, err=%v", job.Path, err)
		}
	}
	return nil
}

func (job *EnqueuedJob) ID() string {
	return job.Path
}

// albumID returns the album ID of the created (or existent) album in PhotosService.
func (job *EnqueuedJob) albumID() (string, error) {
	// Return if empty to avoid a PhotosService call.
	if job.AlbumName == "" {
		return "", nil
	}

	album, err := job.PhotosService.GetOrCreateAlbumByName(job.AlbumName)
	if err != nil {
		return "", fmt.Errorf("Album creation failed: name=%s, error=%s", job.AlbumName, err)
	}
	return album.Id, nil
}
