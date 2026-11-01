package domain

import (
	"errors"
	"fmt"
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

// ---------- Weapon ----------

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

// ---------- Armor ----------

type Armor struct {
	Title   string
	Mass    float64
	Defense int
}

func (a *Armor) Name() string    { return a.Title }
func (a *Armor) Weight() float64 { return a.Mass }

func (a *Armor) Use() (string, error) {
	return fmt.Sprintf("Надета %s (защита: %d)", a.Title, a.Defense), nil
}

// ---------- Potion ----------

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

// compile-time проверки
var (
	_ Item = (*Weapon)(nil)
	_ Item = (*Armor)(nil)
	_ Item = (*Potion)(nil)
)
