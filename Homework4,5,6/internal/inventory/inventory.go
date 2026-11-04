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

// ---------- Filter / Map / Find ----------

// Filter возвращает предметы, удовлетворяющие предикату.
func (inv *Inventory) Filter(pred func(items.Item) bool) []items.Item {
	var res []items.Item
	for _, e := range inv.entries {
		if pred(e.Item) {
			res = append(res, e.Item)
		}
	}
	return res
}

// Map применяет функцию к каждому предмету и возвращает результаты.
func (inv *Inventory) Map(fn func(items.Item) items.Item) []items.Item {
	res := make([]items.Item, 0, len(inv.entries))
	for _, e := range inv.entries {
		res = append(res, fn(e.Item))
	}
	return res
}

// Find возвращает первый предмет, удовлетворяющий предикату.
func (inv *Inventory) Find(pred func(items.Item) bool) (items.Item, bool) {
	for _, e := range inv.entries {
		if pred(e.Item) {
			return e.Item, true
		}
	}
	return nil, false
}

// ---------- методы из ДЗ №4 ----------

// Weapons — только оружие.
func (inv *Inventory) Weapons() []items.Item {
	return inv.Filter(func(it items.Item) bool {
		_, ok := it.(*items.Weapon)
		return ok
	})
}

// BrokenItems — сломанные/исчерпанные предметы.
func (inv *Inventory) BrokenItems() []items.Item {
	return inv.Filter(func(it items.Item) bool {
		switch v := it.(type) {
		case *items.Weapon:
			return v.Durability <= 0
		case *items.Armor:
			return v.Broken
		case *items.Potion:
			return v.Charges <= 0
		}
		return false
	})
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
	return inv.Find(func(it items.Item) bool { return it.Name() == name })
}

// ---------- методы из ДЗ №5 ----------

// UseDangerousItem вызывает Use и восстанавливается после паники,
// превращая её в ошибку. Нужен для предметов, чей Use может паниковать.
// Заодно распаковывает *items.ParseItemError через errors.As.
func (inv *Inventory) UseDangerousItem(id int) (msg string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during use id %d: %v", id, r)
		}
	}()

	i, ferr := inv.findIndex(id)
	if ferr != nil {
		return "", ferr
	}

	msg, uerr := inv.entries[i].Item.Use()
	if uerr != nil {
		var pie *items.ParseItemError
		if errors.As(uerr, &pie) {
			return "", fmt.Errorf("parse error at line %d (%q): %w", pie.Line, pie.Value, uerr)
		}
		return "", uerr
	}
	return msg, nil
}