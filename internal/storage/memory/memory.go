package memory

import (
	"sync"
	"url-service/internal/storage"
)

type Memory struct {
	mu         sync.RWMutex
	aliasToURL map[string]string
	urlSet     map[string]struct{}
}

func NewMemory() (*Memory, error) {
	return &Memory{
		aliasToURL: make(map[string]string),
		urlSet:     make(map[string]struct{}),
	}, nil
}

func (memory *Memory) SaveURL(url string, alias string) error {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	if _, exists := memory.aliasToURL[alias]; exists {
		return storage.ErrAliasExists
	}

	if _, exists := memory.urlSet[url]; exists {
		return storage.ErrURLExists
	}

	memory.aliasToURL[alias] = url
	memory.urlSet[url] = struct{}{}
	return nil
}
func (memory *Memory) GetURL(alias string) (string, error) {
	memory.mu.RLock()
	defer memory.mu.RUnlock()

	url, exists := memory.aliasToURL[alias]
	if !exists {
		return "", storage.ErrURLNotFound
	}
	return url, nil
}

var _ storage.Storage = (*Memory)(nil)
