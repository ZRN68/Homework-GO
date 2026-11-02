package items

import "fmt"

type Armor struct {
	Title   string
	Mass    float64
	Defense int
	Broken  bool
}

func (a *Armor) Name() string    { return a.Title }
func (a *Armor) Weight() float64 { return a.Mass }

// Use меняет состояние: уменьшает защиту. Когда защита кончилась — броня сломана.
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
