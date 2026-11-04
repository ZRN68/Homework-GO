package items

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Potion struct {
	Title   string
	Mass    float64
	Effect  string
	Charges int
}

func (p *Potion) Name() string    { return p.Title }
func (p *Potion) Weight() float64 { return p.Mass }

func (p *Potion) Use() (string, error) {
	if p.Charges <= 0 {
		return "", fmt.Errorf("potion %q: %w", p.Title, ErrExhausted)
	}
	p.Charges--
	return fmt.Sprintf("Выпито %s (эффект %s, осталось зарядов %d)",
		p.Title, p.Effect, p.Charges), nil
}

func (p *Potion) DescribeItem() string {
	return fmt.Sprintf("Potion %q: эффект %q, зарядов %d, вес %g",
		p.Title, p.Effect, p.Charges, p.Mass)
}

func (p *Potion) Serialize() string {
	return fmt.Sprintf("Potion|%s|%s|%d|%g", p.Title, p.Effect, p.Charges, p.Mass)
}

func (p *Potion) Deserialize(s string) error {
	parts := strings.Split(s, "|")
	if len(parts) != 5 || parts[0] != "Potion" {
		return errors.New("invalid potion format")
	}
	var err error
	p.Title = parts[1]
	p.Effect = parts[2]
	if p.Charges, err = strconv.Atoi(parts[3]); err != nil {
		return err
	}
	if p.Mass, err = strconv.ParseFloat(parts[4], 64); err != nil {
		return err
	}
	return nil
}
