package completeduploads

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
	"strings"
)

type LevelDBRepository struct {
	db *leveldb.DB
}

// NewLevelDBRepository create a new repository
func NewLevelDBRepository(db *leveldb.DB) *LevelDBRepository {
	return &LevelDBRepository{db: db}
}

// Close closes the DB.
func (r *LevelDBRepository) Close() error {
	return r.db.Close()
}

// Get an item
func (r *LevelDBRepository) Get(path string) (CompletedUploadedFileItem, error) {
	val, err := r.db.Get([]byte(path), nil)
	if err != nil {
		return CompletedUploadedFileItem{}, err
	}

	return fromStoreValue(path, val)
}

// Store an item
func (r *LevelDBRepository) Put(item CompletedUploadedFileItem) error {
	return r.db.Put(
		[]byte(item.path),
		[]byte(toStoreValue(item)),
		nil)
}

// Delete an item
func (r *LevelDBRepository) Delete(path string) error {
	err := r.db.Delete([]byte(path), nil)
	if err != nil {
		return ErrCannotBeDeleted
	}
	return nil
}

// Parse the pipe delimited parts into a fully formed file item
func fromStoreValue(path string, value []byte) (CompletedUploadedFileItem, error) {
	item := CompletedUploadedFileItem{
		path: path,
	}

	parts := strings.Split(string(value), "|")

	if len(parts) < 2 {
		// FIXME: include the serialized form in the error
		return item, fmt.Errorf("Unsupported serialization format")
	}

	hash, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return item, err
	}
	item.hash = uint32(hash)

	cacheMtime, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return item, err
	}

	item.modifyTime = cacheMtime

	return item, nil
}

// We store the items values as a pipe delimited string
func toStoreValue(item CompletedUploadedFileItem) string {
	return strconv.FormatInt(item.modifyTime, 10) + "|" + fmt.Sprint(item.hash)
}
