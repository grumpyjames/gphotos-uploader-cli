package completeduploads

import (
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"testing"
)

func TestService_PutAndGetWithLevelDb(t *testing.T) {
	item := CompletedUploadedFileItem{
		path:             "foo/bar/baz",
		hash:             235235,
		modifyTime:       348539868,
		uploadToken:      "423453asasdser2342",
		albumTitle:       "my-photo-album",
		mediaItemCreated: false,
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

func TestService_ParsePreviousRecordFormatLevelDb(t *testing.T) {
	// FIXME: what's the idiomatic go way of creating temporary files in tests?
	f := "/tmp/leveldb.bin"
	defer os.RemoveAll(f)

	ft, err := leveldb.OpenFile(f, nil)
	if err != nil {
		t.Errorf("Failed to open level db")
	}

	ft.Put(
		[]byte("foo/banana/baz"),
		[]byte("47588699|235235"),
		nil,
	)

	repo := NewLevelDBRepository(ft)

	retrieved, err := repo.Get("foo/banana/baz")

	if err != nil {
		t.Errorf("Saw %#v", err)
	}

	expected := CompletedUploadedFileItem{
		path:             "foo/banana/baz",
		hash:             235235,
		modifyTime:       47588699,
		uploadToken:      "",
		albumTitle:       "",
		mediaItemCreated: true,
	}

	if retrieved != expected {
		t.Errorf("got %#v, want %#v", retrieved, expected)
	}
}
