package items

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Armor struct {
	Title   string
	Mass    float64
	Defense int
	Broken  bool
}

func (a *Armor) Name() string    { return a.Title }
func (a *Armor) Weight() float64 { return a.Mass }

func (a *Armor) Use() (string, error) {
	if a.Broken {
		return "", fmt.Errorf("armor %q: %w", a.Title, ErrBroken)
	}
	a.Defense--
	if a.Defense <= 0 {
		a.Defense = 0
		a.Broken = true
	}
	return fmt.Sprintf("Надета %s (защита: %d)", a.Title, a.Defense), nil
}

func (a *Armor) DescribeItem() string {
	return fmt.Sprintf("Armor %q: защита %d, вес %g", a.Title, a.Defense, a.Mass)
}

func (a *Armor) Serialize() string {
	broken := 0
	if a.Broken {
		broken = 1
	}
	return fmt.Sprintf("Armor|%s|%d|%d|%g", a.Title, a.Defense, broken, a.Mass)
}

func (a *Armor) Deserialize(s string) error {
	parts := strings.Split(s, "|")
	if len(parts) != 5 || parts[0] != "Armor" {
		return errors.New("invalid armor format")
	}
	var err error
	a.Title = parts[1]
	if a.Defense, err = strconv.Atoi(parts[2]); err != nil {
		return err
	}
	brokenInt, err := strconv.Atoi(parts[3])
	if err != nil {
		return err
	}
	a.Broken = brokenInt != 0
	if a.Mass, err = strconv.ParseFloat(parts[4], 64); err != nil {
		return err
	}
	return nil
}
