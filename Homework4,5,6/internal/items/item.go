package items

import (
	"errors"
	"reflect"
)

var (
	ErrBroken    = errors.New("item is broken")
	ErrExhausted = errors.New("item is exhausted")
)

// Item — контракт предмета.
type Item interface {
	Name() string
	Weight() float64
	Use() (string, error)
}

// Storable — предмет умеет себя сериализовать и десериализовать.
type Storable interface {
	Serialize() string
	Deserialize(s string) error
}

// Describable — предмет умеет себя описывать (без изменения состояния).
type Describable interface {
	DescribeItem() string
}

// Describe — общий дескрайб с проверкой nil и typed nil.
func Describe(it Item) string {
	if it == nil {
		return "<nil item>"
	}
	v := reflect.ValueOf(it)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return "<nil item>"
	}
	if d, ok := it.(Describable); ok {
		return d.DescribeItem()
	}
	return it.Name()
}

// compile-time проверки
var (
	_ Item        = (*Weapon)(nil)
	_ Item        = (*Armor)(nil)
	_ Item        = (*Potion)(nil)
	_ Storable    = (*Weapon)(nil)
	_ Storable    = (*Armor)(nil)
	_ Storable    = (*Potion)(nil)
	_ Describable = (*Weapon)(nil)
	_ Describable = (*Armor)(nil)
	_ Describable = (*Potion)(nil)
)
