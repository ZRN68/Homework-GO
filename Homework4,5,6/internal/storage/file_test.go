package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"homework6/internal/inventory"
	"homework6/internal/items"
)

func TestFileStorage_SaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test_inventory.txt")
	store := NewFile(path)

	inv := inventory.New()
	w := &items.Weapon{Title: "Меч", Mass: 3.5, Damage: 10, Durability: 5}
	a := &items.Armor{Title: "Щит", Mass: 8, Defense: 15, Broken: true}
	p := &items.Potion{Title: "Зелье", Mass: 0.5, Effect: "HP", Charges: 2}

	idW := inv.Add(w)
	idA := inv.Add(a)
	idP := inv.Add(p)

	if err := store.Save(inv); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		t.Fatalf("file empty or unreadable: %v", err)
	}

	loadedInv, err := store.Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if loadedInv.Len() != 3 {
		t.Fatalf("Load returned %d items, want 3", loadedInv.Len())
	}

	entries := loadedInv.Entries()

	loadedW, ok := entries[0].Item.(*items.Weapon)
	if !ok || entries[0].ID != idW {
		t.Fatalf("weapon ID or type loaded incorrectly, got ID %d", entries[0].ID)
	}
	if *loadedW != *w {
		t.Errorf("weapon fields mismatch: got %+v, want %+v", loadedW, w)
	}

	loadedA, ok := entries[1].Item.(*items.Armor)
	if !ok || entries[1].ID != idA {
		t.Fatalf("armor ID or type loaded incorrectly, got ID %d", entries[1].ID)
	}
	if *loadedA != *a {
		t.Errorf("armor fields mismatch: got %+v, want %+v", loadedA, a)
	}

	loadedP, ok := entries[2].Item.(*items.Potion)
	if !ok || entries[2].ID != idP {
		t.Fatalf("potion ID or type loaded incorrectly, got ID %d", entries[2].ID)
	}
	if *loadedP != *p {
		t.Errorf("potion fields mismatch: got %+v, want %+v", loadedP, p)
	}

	newID := loadedInv.Add(&items.Weapon{Title: "Новый"})
	if newID != 4 {
		t.Errorf("next ID after load = %d, want 4", newID)
	}
}

func TestFileStorage_LoadNonExistent(t *testing.T) {
	store := NewFile("does_not_exist_file.txt")
	inv, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if inv.Len() != 0 {
		t.Errorf("expected empty inventory")
	}
}

func TestFileStorage_LoadCorruptedAndDuplicates(t *testing.T) {
	content := strings.Join([]string{
		"1|Weapon|Меч|10|5|3.5",
		"1|Armor|Щит|15|0|8",
		"bad line without pipes",
		"3|Potion|Зелье|HP|1|0.5",
	}, "\n")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "corrupt.txt")
	os.WriteFile(path, []byte(content), 0644)

	store := NewFile(path)
	inv, err := store.Load()

	if err == nil {
		t.Fatal("expected error for corrupted/duplicate lines, got nil")
	}

	if inv.Len() != 2 {
		t.Errorf("expected 2 valid items, got %d", inv.Len())
	}

	var target *items.ParseItemError
	if !errors.As(err, &target) {
		t.Fatalf("expected ParseItemError, got %T", err)
	}

	if target.Line != 2 {
		t.Errorf("ParseItemError Line = %d, want 2", target.Line)
	}
	if !strings.Contains(target.Value, "Armor") {
		t.Errorf("ParseItemError Value = %q, want contains 'Armor'", target.Value)
	}

	if !errors.Is(err, ErrDuplicateID) {
		t.Errorf("expected error chain to contain ErrDuplicateID, got %v", err)
	}
}

func TestMemoryStorage_SaveLoad(t *testing.T) {
	mem := NewMemory()

	inv, err := mem.Load()
	if err != nil || inv.Len() != 0 {
		t.Fatalf("empty memory load failed")
	}

	inv.Add(&items.Weapon{Title: "Меч"})
	if err := mem.Save(inv); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := mem.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.Len() != 1 {
		t.Errorf("load returned %d items, want 1", loaded.Len())
	}
}
