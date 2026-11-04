package items

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Weapon struct {
	Title      string
	Mass       float64
	Damage     int
	Durability int
}

func (w *Weapon) Name() string    { return w.Title }
func (w *Weapon) Weight() float64 { return w.Mass }

func (w *Weapon) Use() (string, error) {
	if w.Durability <= 0 {
		return "", fmt.Errorf("weapon %q: %w", w.Title, ErrBroken)
	}
	w.Durability--
	return fmt.Sprintf("Используется %s (урон: %d, прочность: %d)",
		w.Title, w.Damage, w.Durability), nil
}

func (w *Weapon) DescribeItem() string {
	return fmt.Sprintf("Weapon %q: урон %d, прочность %d, вес %g",
		w.Title, w.Damage, w.Durability, w.Mass)
}

func (w *Weapon) Serialize() string {
	return fmt.Sprintf("Weapon|%s|%d|%d|%g", w.Title, w.Damage, w.Durability, w.Mass)
}

func (w *Weapon) Deserialize(s string) error {
	parts := strings.Split(s, "|")
	if len(parts) != 5 || parts[0] != "Weapon" {
		return errors.New("invalid weapon format")
	}
	var err error
	w.Title = parts[1]
	if w.Damage, err = strconv.Atoi(parts[2]); err != nil {
		return err
	}
	if w.Durability, err = strconv.Atoi(parts[3]); err != nil {
		return err
	}
	if w.Mass, err = strconv.ParseFloat(parts[4], 64); err != nil {
		return err
	}
	return nil
}
