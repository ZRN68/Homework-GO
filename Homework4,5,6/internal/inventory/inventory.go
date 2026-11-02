package inventory

import (
	"errors"
	"fmt"

	"homework6/internal/items"
)

var ErrNotFound = errors.New("item not found")

// Entry — запись инвентаря: ID живёт вместе с предметом,
// а не вычисляется из позиции в слайсе.
type Entry struct {
	ID   int
	Item items.Item
}

// Inventory — набор записей с уникальными ID.
type Inventory struct {
	entries []Entry
	nextID  int
}

func New() *Inventory {
	return &Inventory{nextID: 1}
}

// Len — количество предметов.
func (inv *Inventory) Len() int { return len(inv.entries) }

// Entries — копия списка записей (чтобы снаружи не портили инвариант).
func (inv *Inventory) Entries() []Entry {
	out := make([]Entry, len(inv.entries))
	copy(out, inv.entries)
	return out
}

// Add возвращает ID добавленного предмета.
func (inv *Inventory) Add(item items.Item) int {
	id := inv.nextID
	inv.nextID++
	inv.entries = append(inv.entries, Entry{ID: id, Item: item})
	return id
}

// AddWithID нужен хранилищу при загрузке: кладём предмет с уже известным ID.
func (inv *Inventory) AddWithID(id int, item items.Item) {
	inv.entries = append(inv.entries, Entry{ID: id, Item: item})
	if id >= inv.nextID {
		inv.nextID = id + 1
	}
}

func (inv *Inventory) findIndex(id int) (int, error) {
	for i, e := range inv.entries {
		if e.ID == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("id %d: %w", id, ErrNotFound)
}

func (inv *Inventory) Delete(id int) error {
	i, err := inv.findIndex(id)
	if err != nil {
		return err
	}
	inv.entries = append(inv.entries[:i], inv.entries[i+1:]...)
	return nil
}

func (inv *Inventory) Use(id int) (string, error) {
	i, err := inv.findIndex(id)
	if err != nil {
		return "", err
	}
	return inv.entries[i].Item.Use()
}

// ---------- методы из ДЗ №4 ----------

// Weapons — только оружие.
func (inv *Inventory) Weapons() []items.Item {
	var res []items.Item
	for _, e := range inv.entries {
		if _, ok := e.Item.(*items.Weapon); ok {
			res = append(res, e.Item)
		}
	}
	return res
}

// BrokenItems — сломанные/исчерпанные предметы.
func (inv *Inventory) BrokenItems() []items.Item {
	var res []items.Item
	for _, e := range inv.entries {
		switch v := e.Item.(type) {
		case *items.Weapon:
			if v.Durability <= 0 {
				res = append(res, e.Item)
			}
		case *items.Armor:
			if v.Broken {
				res = append(res, e.Item)
			}
		case *items.Potion:
			if v.Charges <= 0 {
				res = append(res, e.Item)
			}
		}
	}
	return res
}

// ItemNames — имена всех предметов.
func (inv *Inventory) ItemNames() []string {
	res := make([]string, 0, len(inv.entries))
	for _, e := range inv.entries {
		res = append(res, e.Item.Name())
	}
	return res
}

// FindItemByName — поиск предмета по точному имени.
func (inv *Inventory) FindItemByName(name string) (items.Item, bool) {
	for _, e := range inv.entries {
		if e.Item.Name() == name {
			return e.Item, true
		}
	}
	return nil, false
}
