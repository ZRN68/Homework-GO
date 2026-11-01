package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"homework6/internal/domain"
	"homework6/internal/inventory"
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

// FileStorage — файловое хранилище.
type FileStorage struct {
	Path string
}

func NewFile(path string) *FileStorage { return &FileStorage{Path: path} }

func (s *FileStorage) Load() (*inventory.Inventory, error) {
	inv := inventory.New()

	f, err := os.Open(s.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return inv, nil
		}
		return nil, fmt.Errorf("open %s: %w", s.Path, err)
	}
	defer f.Close()

	if err := loadFrom(inv, f); err != nil {
		return inv, fmt.Errorf("load %s: %w", s.Path, err)
	}
	return inv, nil
}

func (s *FileStorage) Save(inv *inventory.Inventory) (err error) {
	f, err := os.Create(s.Path)
	if err != nil {
		return fmt.Errorf("create %s: %w", s.Path, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close %s: %w", s.Path, cerr)
		}
	}()
	if err := saveTo(inv, f); err != nil {
		return fmt.Errorf("save %s: %w", s.Path, err)
	}
	return nil
}

// ---------- общий формат ----------

func serializeItem(id int, item domain.Item) (string, bool) {
	switch v := item.(type) {
	case *domain.Weapon:
		return fmt.Sprintf("%d|Weapon|%s|%d|%d|%g\n", id, v.Title, v.Damage, v.Durability, v.Mass), true
	case *domain.Armor:
		return fmt.Sprintf("%d|Armor|%s|%d|%g\n", id, v.Title, v.Defense, v.Mass), true
	case *domain.Potion:
		return fmt.Sprintf("%d|Potion|%s|%s|%d|%g\n", id, v.Title, v.Effect, v.Charges, v.Mass), true
	}
	return "", false
}

func parseItemLine(line string) (domain.Item, int, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 2 {
		return nil, 0, fmt.Errorf("too few fields")
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, 0, fmt.Errorf("bad id: %w", err)
	}
	switch parts[1] {
	case "Weapon":
		if len(parts) != 6 {
			return nil, 0, fmt.Errorf("weapon needs 6 fields, got %d", len(parts))
		}
		dmg, e1 := strconv.Atoi(parts[3])
		dur, e2 := strconv.Atoi(parts[4])
		wt, e3 := strconv.ParseFloat(parts[5], 64)
		if err := errors.Join(e1, e2, e3); err != nil {
			return nil, 0, err
		}
		return &domain.Weapon{Title: parts[2], Damage: dmg, Durability: dur, Mass: wt}, id, nil
	case "Armor":
		if len(parts) != 5 {
			return nil, 0, fmt.Errorf("armor needs 5 fields, got %d", len(parts))
		}
		def, e1 := strconv.Atoi(parts[3])
		wt, e2 := strconv.ParseFloat(parts[4], 64)
		if err := errors.Join(e1, e2); err != nil {
			return nil, 0, err
		}
		return &domain.Armor{Title: parts[2], Defense: def, Mass: wt}, id, nil
	case "Potion":
		if len(parts) != 6 {
			return nil, 0, fmt.Errorf("potion needs 6 fields, got %d", len(parts))
		}
		ch, e1 := strconv.Atoi(parts[4])
		wt, e2 := strconv.ParseFloat(parts[5], 64)
		if err := errors.Join(e1, e2); err != nil {
			return nil, 0, err
		}
		return &domain.Potion{Title: parts[2], Effect: parts[3], Charges: ch, Mass: wt}, id, nil
	}
	return nil, 0, fmt.Errorf("unknown type %q", parts[1])
}

// loadFrom — общая логика загрузки из io.Reader.
func loadFrom(inv *inventory.Inventory, r io.Reader) error {
	sc := bufio.NewScanner(r)
	lineNo := 0
	var errs []error

	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		item, id, err := parseItemLine(line)
		if err != nil {
			errs = append(errs, &ParseItemError{Line: lineNo, Value: line, Err: err})
			continue
		}
		inv.Items = append(inv.Items, item)
		inv.SetNextID(id + 1)
	}
	if err := sc.Err(); err != nil {
		errs = append(errs, fmt.Errorf("scan: %w", err))
	}
	return errors.Join(errs...)
}

// saveTo — общая логика сохранения в io.Writer.
func saveTo(inv *inventory.Inventory, w io.Writer) error {
	bw := bufio.NewWriter(w)
	for i, item := range inv.Items {
		line, ok := serializeItem(i+1, item)
		if !ok {
			continue
		}
		if _, err := bw.WriteString(line); err != nil {
			return fmt.Errorf("write item %d: %w", i+1, err)
		}
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	return nil
}
