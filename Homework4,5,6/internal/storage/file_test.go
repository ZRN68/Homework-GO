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
	inv.Add(&items.Weapon{Title: "Меч", Mass: 3.5, Damage: 10, Durability: 5})
	inv.Add(&items.Armor{Title: "Щит", Mass: 8, Defense: 15})
	inv.Add(&items.Potion{Title: "Зелье", Mass: 0.5, Effect: "HP", Charges: 2})

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
	if w, ok := entries[0].Item.(*items.Weapon); !ok || w.Title != "Меч" {
		t.Errorf("weapon loaded incorrectly")
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
		t.Errorf("expected ParseItemError, got %T", err)
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
