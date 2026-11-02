package storage

import (
	"bytes"
	"fmt"

	"homework6/internal/inventory"
)

// MemoryStorage — хранилище в памяти. Соблюдает те же правила, что FileStorage.
type MemoryStorage struct {
	data []byte
}

func NewMemory() *MemoryStorage { return &MemoryStorage{} }

func NewMemoryFromData(data []byte) *MemoryStorage {
	return &MemoryStorage{data: append([]byte(nil), data...)}
}

func (s *MemoryStorage) Load() (*inventory.Inventory, error) {
	inv := inventory.New()
	if len(s.data) == 0 {
		return inv, nil
	}
	if err := loadFrom(inv, bytes.NewReader(s.data)); err != nil {
		return inv, fmt.Errorf("load memory: %w", err)
	}
	return inv, nil
}

func (s *MemoryStorage) Save(inv *inventory.Inventory) error {
	var buf bytes.Buffer
	if err := saveTo(inv, &buf); err != nil {
		return fmt.Errorf("save memory: %w", err)
	}
	s.data = buf.Bytes()
	return nil
}
