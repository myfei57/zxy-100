package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Store struct {
	root string
	mu   sync.RWMutex
}

func New(root string) *Store {
	return &Store{root: root}
}

func (s *Store) Path(key string) string {
	name := strings.NewReplacer("/", "_", ".", "_", "\\", "_").Replace(key)
	return filepath.Join(s.root, name+".json")
}

func (s *Store) Put(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return writeFileAtomic(s.Path(key), data)
}

func (s *Store) Get(key string, out any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := os.ReadFile(s.Path(key))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, err := os.Stat(s.Path(key))
	return err == nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return os.Remove(s.Path(key))
}

func (s *Store) List() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, err := os.ReadDir(s.root)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".json") {
			names = append(names, strings.TrimSuffix(name, ".json"))
		}
	}
	sort.Strings(names)
	return names, nil
}
