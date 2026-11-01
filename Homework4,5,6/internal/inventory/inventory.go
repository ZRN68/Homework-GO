package inventory

import (
	"errors"
	"fmt"

	"homework6/internal/domain"
)

var ErrNotFound = errors.New("item not found")

// Item — то, что хранит инвентарь. Дублируем контракт здесь,
// чтобы пакет не зависел от domain целиком (маленький интерфейс со стороны потребителя).
type Item = domain.Item

// Inventory — набор предметов с числовыми ID (ID = индекс+1).
type Inventory struct {
	Items  []Item
	nextID int
}

func New() *Inventory {
	return &Inventory{nextID: 1}
}

// Add возвращает ID предмета.
func (inv *Inventory) Add(item Item) int {
	id := inv.nextID
	inv.nextID++
	inv.Items = append(inv.Items, item)
	return id
}

// SetNextID нужен хранилищу при загрузке, чтобы ID не сбрасывались.
func (inv *Inventory) SetNextID(id int) {
	if id > inv.nextID {
		inv.nextID = id
	}
}

func (inv *Inventory) itemByID(id int) (Item, int, error) {
	idx := id - 1
	if idx < 0 || idx >= len(inv.Items) {
		return nil, -1, fmt.Errorf("id %d: %w", id, ErrNotFound)
	}
	return inv.Items[idx], idx, nil
}

func (inv *Inventory) Delete(id int) error {
	_, idx, err := inv.itemByID(id)
	if err != nil {
		return err
	}
	inv.Items = append(inv.Items[:idx], inv.Items[idx+1:]...)
	return nil
}

func (inv *Inventory) Use(id int) (string, error) {
	it, _, err := inv.itemByID(id)
	if err != nil {
		return "", err
	}
	return it.Use()
}
