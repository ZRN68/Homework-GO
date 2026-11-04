package items

import (
	"errors"
	"fmt"
	"strings"
)

// ParseItemError — типизированная ошибка разбора строки.
type ParseItemError struct {
	Line  int
	Value string
	Err   error
}

func (e *ParseItemError) Error() string {
	return fmt.Sprintf("parse error at line %d (%q): %v", e.Line, e.Value, e.Err)
}

func (e *ParseItemError) Unwrap() error { return e.Err }

// ParseItem разбирает строку формата «Type|field|field|...».
// ID предмета сюда не входит — его ведёт инвентарь.
func ParseItem(line string) (Item, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 2 {
		return nil, errors.New("too few fields")
	}

	switch parts[0] {
	case "Weapon":
		w := &Weapon{}
		if err := w.Deserialize(line); err != nil {
			return nil, err
		}
		return w, nil

	case "Armor":
		a := &Armor{}
		if err := a.Deserialize(line); err != nil {
			return nil, err
		}
		return a, nil

	case "Potion":
		p := &Potion{}
		if err := p.Deserialize(line); err != nil {
			return nil, err
		}
		return p, nil
	}

	return nil, fmt.Errorf("unknown type %q", parts[0])
}
