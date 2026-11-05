package items

import (
	"errors"
	"strings"
	"testing"
)

func TestWeaponUse(t *testing.T) {
	tests := []struct {
		name    string
		weapon  *Weapon
		wantErr error
		wantDur int
	}{
		{"success", &Weapon{Title: "Меч", Durability: 1}, nil, 0},
		{"broken", &Weapon{Title: "Меч", Durability: 0}, ErrBroken, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.weapon.Use()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.weapon.Durability != tt.wantDur {
				t.Errorf("Durability = %d, want %d", tt.weapon.Durability, tt.wantDur)
			}
		})
	}
}

func TestPotionUse(t *testing.T) {
	tests := []struct {
		name        string
		potion      *Potion
		wantErr     error
		wantCharges int
	}{
		{"success", &Potion{Title: "Зелье", Charges: 1}, nil, 0},
		{"exhausted", &Potion{Title: "Зелье", Charges: 0}, ErrExhausted, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.potion.Use()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.potion.Charges != tt.wantCharges {
				t.Errorf("Charges = %d, want %d", tt.potion.Charges, tt.wantCharges)
			}
		})
	}
}

func TestArmorUse(t *testing.T) {
	a := &Armor{Title: "Щит", Defense: 1}
	_, err := a.Use()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if a.Defense != 0 || !a.Broken {
		t.Errorf("Armor state not changed: defense=%d, broken=%v", a.Defense, a.Broken)
	}

	_, err = a.Use()
	if !errors.Is(err, ErrBroken) {
		t.Fatalf("want ErrBroken, got %v", err)
	}
}

func TestDeserialize(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		obj     Storable
		wantErr bool
		check   func(Storable) bool
	}{
		{
			name: "valid weapon",
			line: "Weapon|Меч|10|5|3.5",
			obj:  &Weapon{},
			check: func(s Storable) bool {
				w := s.(*Weapon)
				return w.Title == "Меч" && w.Damage == 10 && w.Durability == 5 && w.Mass == 3.5
			},
		},
		{
			name:    "invalid weapon fields",
			line:    "Weapon|Меч|10",
			obj:     &Weapon{},
			wantErr: true,
		},
		{
			name:    "invalid type",
			line:    "Sword|Меч|10|5|3.5",
			obj:     &Weapon{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.obj.Deserialize(tt.line)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !tt.check(tt.obj) {
				t.Errorf("object not populated correctly")
			}
		})
	}
}

func TestDescribe(t *testing.T) {
	var typedNil *Weapon = nil
	var nilItem Item = nil

	tests := []struct {
		name string
		in   Item
		want string
	}{
		{"nil interface", nilItem, "<nil item>"},
		{"typed nil", typedNil, "<nil item>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Describe(tt.in); got != tt.want {
				t.Errorf("Describe() = %q, want %q", got, tt.want)
			}
		})
	}

	w := &Weapon{Title: "Меч", Damage: 10}
	if !strings.Contains(Describe(w), "Меч") {
		t.Errorf("Describe() missing name")
	}
}

func TestItemGettersAndSerialize(t *testing.T) {
	w := &Weapon{Title: "Меч", Mass: 3.5, Damage: 10, Durability: 5}
	if w.Name() != "Меч" || w.Weight() != 3.5 {
		t.Errorf("Weapon Name/Weight failed")
	}
	if w.Serialize() != "Weapon|Меч|10|5|3.5" {
		t.Errorf("Weapon Serialize failed: %s", w.Serialize())
	}

	a := &Armor{Title: "Щит", Mass: 8.0, Defense: 15, Broken: true}
	if a.Name() != "Щит" || a.Weight() != 8.0 {
		t.Errorf("Armor Name/Weight failed")
	}
	if a.Serialize() != "Armor|Щит|15|1|8" {
		t.Errorf("Armor Serialize failed: %s", a.Serialize())
	}

	p := &Potion{Title: "Зелье", Mass: 0.5, Effect: "HP", Charges: 2}
	if p.Name() != "Зелье" || p.Weight() != 0.5 {
		t.Errorf("Potion Name/Weight failed")
	}
	if p.Serialize() != "Potion|Зелье|HP|2|0.5" {
		t.Errorf("Potion Serialize failed: %s", p.Serialize())
	}
}

func TestParseItem(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{"valid weapon", "Weapon|Меч|10|5|3.5", false},
		{"valid armor", "Armor|Щит|15|0|8", false},
		{"valid potion", "Potion|Зелье|HP|2|0.5", false},
		{"unknown type", "Sword|Меч|10|5|3.5", true},
		{"too few fields", "Weapon", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseItem(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseItem() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
