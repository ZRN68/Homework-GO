package inventory

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"homework6/internal/items"
)

type panicItem struct{}

func (p *panicItem) Name() string         { return "Panic" }
func (p *panicItem) Weight() float64      { return 1 }
func (p *panicItem) Use() (string, error) { panic("boom") }

func TestInventoryBasic(t *testing.T) {
	inv := New()

	if id := inv.Add(&items.Weapon{Title: "Меч"}); id != 1 {
		t.Errorf("first ID should be 1, got %d", id)
	}
	if id := inv.Add(&items.Armor{Title: "Щит"}); id != 2 {
		t.Errorf("second ID should be 2, got %d", id)
	}

	if inv.Len() != 2 {
		t.Errorf("Len() = %d, want 2", inv.Len())
	}

	if err := inv.Delete(1); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if inv.Len() != 1 {
		t.Errorf("Len() after delete = %d, want 1", inv.Len())
	}

	if err := inv.Delete(99); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestInventoryFilters(t *testing.T) {
	inv := New()
	inv.Add(&items.Weapon{Title: "Меч", Durability: 10})
	inv.Add(&items.Armor{Title: "Щит", Broken: false})
	inv.Add(&items.Weapon{Title: "Сломанный меч", Durability: 0})

	weapons := inv.Weapons()
	if len(weapons) != 2 {
		t.Errorf("Weapons() = %d, want 2", len(weapons))
	}

	broken := inv.BrokenItems()
	if len(broken) != 1 || broken[0].Name() != "Сломанный меч" {
		t.Errorf("BrokenItems() failed, got %v", broken)
	}

	names := inv.ItemNames()
	if len(names) != 3 {
		t.Errorf("ItemNames() = %d, want 3", len(names))
	}

	if it, ok := inv.FindItemByName("Щит"); !ok || it.Name() != "Щит" {
		t.Errorf("FindItemByName failed")
	}
}

func TestUseDangerousItem(t *testing.T) {
	inv := New()
	okID := inv.Add(&items.Weapon{Title: "Меч", Durability: 1})
	panicID := inv.Add(&panicItem{})

	tests := []struct {
		name    string
		id      int
		wantErr bool
		wantMsg string
	}{
		{"ok", okID, false, "Используется Меч"},
		{"panic recovered", panicID, true, "panic during use"},
		{"not found", 999, true, "item not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := inv.UseDangerousItem(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantMsg != "" {
				if err != nil && !strings.Contains(err.Error(), tt.wantMsg) {
					t.Errorf("error %q missing %q", err.Error(), tt.wantMsg)
				}
				if err == nil && !strings.Contains(msg, tt.wantMsg) {
					t.Errorf("msg %q missing %q", msg, tt.wantMsg)
				}
			}
		})
	}
}

func TestInventoryMore(t *testing.T) {
	inv := New()
	inv.AddWithID(10, &items.Weapon{Title: "Меч", Durability: 5})
	inv.AddWithID(5, &items.Armor{Title: "Щит"})

	if len(inv.Items()) != 2 {
		t.Errorf("Items() len = %d, want 2", len(inv.Items()))
	}

	mapped := inv.Map(func(it items.Item) items.Item {
		return &items.Weapon{Title: "Новый " + it.Name()}
	})
	if mapped[0].Name() != "Новый Меч" {
		t.Errorf("Map() failed")
	}

	if it, ok := inv.Find(func(it items.Item) bool { return it.Name() == "Щит" }); !ok || it.Name() != "Щит" {
		t.Errorf("Find() failed")
	}

	if _, err := inv.Use(10); err != nil {
		t.Errorf("Use() failed: %v", err)
	}
	if _, err := inv.Use(99); !errors.Is(err, ErrNotFound) {
		t.Errorf("Use() want ErrNotFound, got %v", err)
	}
}

func prepareInventory() *Inventory {
	inv := New()
	for i := 0; i < 10000; i++ {
		inv.Add(&items.Weapon{Title: fmt.Sprintf("Меч-%d", i)})
	}
	return inv
}

func BenchmarkManualFind(b *testing.B) {
	inv := prepareInventory()
	for b.Loop() {
		for _, e := range inv.entries {
			if e.Item.Name() == "Несуществующий" {
				break
			}
		}
	}
}

func BenchmarkGenericFind(b *testing.B) {
	inv := prepareInventory()
	for b.Loop() {
		_, _ = inv.Find(func(it items.Item) bool {
			return it.Name() == "Несуществующий"
		})
	}
}
