package completeduploads

import (
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"testing"
)

func TestService_PutAndGetWithLevelDb(t *testing.T) {
	item := CompletedUploadedFileItem{
		path:       "foo/bar/baz",
		hash:       235235,
		modifyTime: 34853984868,
	}

	// FIXME: what's the idiomatic go way of creating temporary files in tests?
	f := "/tmp/leveldb.bin"
	defer os.RemoveAll(f)

	ft, err := leveldb.OpenFile(f, nil)
	if err != nil {
		t.Errorf("Failed to open level db")
	}
	repo := NewLevelDBRepository(ft)

	repo.Put(item)

	retrieved, _ := repo.Get(item.path)

	if retrieved != item {
		t.Errorf("got %#v, want %#v", retrieved, item)
	}
}
