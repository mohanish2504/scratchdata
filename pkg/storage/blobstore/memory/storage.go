package memory

import (
	"fmt"
	"io"
	"sync"

	"github.com/scratchdata/scratchdata/pkg/storage/blobstore/models"
)

type Storage struct {
	mu    sync.RWMutex
	items map[string][]byte
}

func (s *Storage) Upload(path string, r io.ReadSeeker) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[path] = data

	return nil
}

func (s *Storage) Download(path string, w io.WriterAt) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.items[path]
	if !ok {
		return models.ErrNotFound
	}
	if _, err := w.WriteAt(data, 0); err != nil {
		return fmt.Errorf("Storage.Download: %s: %w", path, err)
	}
	return nil
}

func (s *Storage) Delete(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// delete the key from the map, won't throw an error if the key doesn't exist
	delete(s.items, path)
	return nil
}

// NewStorage returns a new initialized Storage
func NewStorage(conf map[string]any) (*Storage, error) {
	rc := &Storage{
		items: map[string][]byte{},
	}
	return rc, nil
}
