package inventory

import (
	"errors"
	"fmt"

	"homework6/internal/items"
)

var ErrNotFound = errors.New("item not found")

type Entry struct {
	ID   int
	Item items.Item
}

type Inventory struct {
	entries []Entry
	nextID  int
}

func New() *Inventory { return &Inventory{nextID: 1} }

func (inv *Inventory) Len() int { return len(inv.entries) }

func (inv *Inventory) Entries() []Entry {
	out := make([]Entry, len(inv.entries))
	copy(out, inv.entries)
	return out
}

func (inv *Inventory) Items() []items.Item {
	res := make([]items.Item, 0, len(inv.entries))
	for _, e := range inv.entries {
		res = append(res, e.Item)
	}
	return res
}

func (inv *Inventory) Add(item items.Item) int {
	id := inv.nextID
	inv.nextID++
	inv.entries = append(inv.entries, Entry{ID: id, Item: item})
	return id
}

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

func Filter[T any](src []T, pred func(T) bool) []T {
	var res []T
	for _, v := range src {
		if pred(v) {
			res = append(res, v)
		}
	}
	return res
}

func Map[T, R any](src []T, fn func(T) R) []R {
	res := make([]R, 0, len(src))
	for _, v := range src {
		res = append(res, fn(v))
	}
	return res
}

func Find[T any](src []T, pred func(T) bool) (T, bool) {
	var zero T
	for _, v := range src {
		if pred(v) {
			return v, true
		}
	}
	return zero, false
}

func (inv *Inventory) Filter(pred func(items.Item) bool) []items.Item {
	return Filter(inv.Items(), pred)
}

func (inv *Inventory) Map(fn func(items.Item) items.Item) []items.Item {
	return Map(inv.Items(), fn)
}

func (inv *Inventory) Find(pred func(items.Item) bool) (items.Item, bool) {
	return Find(inv.Items(), pred)
}

func (inv *Inventory) Weapons() []items.Item {
	return inv.Filter(func(it items.Item) bool {
		_, ok := it.(*items.Weapon)
		return ok
	})
}

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

func (inv *Inventory) ItemNames() []string {
	return Map(inv.Items(), func(it items.Item) string { return it.Name() })
}

func (inv *Inventory) FindItemByName(name string) (items.Item, bool) {
	return inv.Find(func(it items.Item) bool { return it.Name() == name })
}

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

// 1. Почему ошибка сломанного предмета не должна быть паникой?
//    Сломанный предмет — это ожидаемая ситуация (штатная ошибка бизнес-логики),
//    а не баг программы. Ожидаемые ошибки возвращаются через error, чтобы их
//    можно было обработать (errors.Is(err, ErrBroken)). Паника — для
//    неисправимых программных сбоев (nil, out of range).
//
// 2. Какая строка после panic не будет выполнена?
//    Все строки функции ниже вызова panic(), включая `return`. Управление
//    сразу перейдет к defer, где сработает recover.
//
// 3. Почему recover должен находиться в отложенной функции?
//    Потому что recover работает только при прямом вызове из defer.
//    В обычном коде он просто вернет nil, и паника продолжит раскрутку стека.
