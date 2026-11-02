package items

import (
	"errors"
	"fmt"
	"strconv"
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
		return nil, fmt.Errorf("too few fields")
	}
	switch parts[0] {
	case "Weapon":
		return parseWeapon(parts)
	case "Armor":
		return parseArmor(parts)
	case "Potion":
		return parsePotion(parts)
	}
	return nil, fmt.Errorf("unknown type %q", parts[0])
}

func parseWeapon(parts []string) (Item, error) {
	if len(parts) != 5 {
		return nil, fmt.Errorf("weapon needs 5 fields, got %d", len(parts))
	}
	dmg, e1 := strconv.Atoi(parts[2])
	dur, e2 := strconv.Atoi(parts[3])
	wt, e3 := strconv.ParseFloat(parts[4], 64)
	if err := errors.Join(e1, e2, e3); err != nil {
		return nil, err
	}
	return &Weapon{Title: parts[1], Damage: dmg, Durability: dur, Mass: wt}, nil
}

func parseArmor(parts []string) (Item, error) {
	if len(parts) != 5 {
		return nil, fmt.Errorf("armor needs 5 fields, got %d", len(parts))
	}
	def, e1 := strconv.Atoi(parts[2])
	broken, e2 := strconv.Atoi(parts[3])
	wt, e3 := strconv.ParseFloat(parts[4], 64)
	if err := errors.Join(e1, e2, e3); err != nil {
		return nil, err
	}
	return &Armor{Title: parts[1], Defense: def, Broken: broken != 0, Mass: wt}, nil
}

func parsePotion(parts []string) (Item, error) {
	if len(parts) != 5 {
		return nil, fmt.Errorf("potion needs 5 fields, got %d", len(parts))
	}
	ch, e1 := strconv.Atoi(parts[3])
	wt, e2 := strconv.ParseFloat(parts[4], 64)
	if err := errors.Join(e1, e2); err != nil {
		return nil, err
	}
	return &Potion{Title: parts[1], Effect: parts[2], Charges: ch, Mass: wt}, nil
}
