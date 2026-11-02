package items

import "errors"

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

// Storable — предмет умеет себя сериализовать в строку файлового формата.
type Storable interface {
	Serialize() string
}

// Describable — предмет умеет себя описывать (без изменения состояния).
type Describable interface {
	DescribeItem() string
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
