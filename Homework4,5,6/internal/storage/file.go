package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"homework6/internal/inventory"
	"homework6/internal/items"
)

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
		// ВАЖНО: возвращаем не-nil инвентарь, чтобы cli не упал в panic.
		return inv, fmt.Errorf("open %s: %w", s.Path, err)
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

// loadFrom читает строки вида «ID|Type|field|...».
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

		// отделяем ID от остальной строки
		bar := strings.IndexByte(line, '|')
		if bar < 0 {
			errs = append(errs, &items.ParseItemError{Line: lineNo, Value: line, Err: errors.New("no id separator")})
			continue
		}
		id, err := strconv.Atoi(line[:bar])
		if err != nil {
			errs = append(errs, &items.ParseItemError{Line: lineNo, Value: line, Err: fmt.Errorf("bad id: %w", err)})
			continue
		}

		item, err := items.ParseItem(line[bar+1:])
		if err != nil {
			errs = append(errs, &items.ParseItemError{Line: lineNo, Value: line, Err: err})
			continue
		}
		inv.AddWithID(id, item)
	}
	if err := sc.Err(); err != nil {
		errs = append(errs, fmt.Errorf("scan: %w", err))
	}
	return errors.Join(errs...)
}

// saveTo пишет «ID|SerializedItem».
func saveTo(inv *inventory.Inventory, w io.Writer) error {
	bw := bufio.NewWriter(w)
	for _, e := range inv.Entries() {
		st, ok := e.Item.(items.Storable)
		if !ok {
			continue
		}
		if _, err := fmt.Fprintf(bw, "%d|%s\n", e.ID, st.Serialize()); err != nil {
			return fmt.Errorf("write item %d: %w", e.ID, err)
		}
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}
	return nil
}
