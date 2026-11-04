package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"strings"

	"homework6/internal/inventory"
	"homework6/internal/items"
)

type Storage interface {
	Load() (*inventory.Inventory, error)
	Save(inv *inventory.Inventory) error
}

type Runner struct {
	Store  Storage
	Stdout io.Writer
	Stderr io.Writer
}

func Run(args []string, stdout, stderr io.Writer, store Storage) int {
	r := &Runner{Store: store, Stdout: stdout, Stderr: stderr}

	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: inventory <list|add|use|delete> [flags]")
		return 2
	}

	cmd, rest := args[0], args[1:]
	switch cmd {
	case "list":
		return r.list(rest)
	case "add":
		return r.add(rest)
	case "use":
		return r.use(rest)
	case "delete":
		return r.delete(rest)
	default:
		fmt.Fprintf(stderr, "error: unknown command %q\n", cmd)
		return 2
	}
}

func (r *Runner) loadOrFail() (*inventory.Inventory, bool) {
	inv, err := r.Store.Load()
	if err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return nil, false
	}
	if inv == nil {
		fmt.Fprintln(r.Stderr, "error: storage returned nil inventory")
		return nil, false
	}
	return inv, true
}

func (r *Runner) list(args []string) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(r.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(r.Stderr, "error: unexpected arguments: %v\n", fs.Args())
		return 2
	}

	inv, ok := r.loadOrFail()
	if !ok {
		return 1
	}
	if inv.Len() == 0 {
		fmt.Fprintln(r.Stdout, "инвентарь пуст")
		return 0
	}
	for _, e := range inv.Entries() {
		fmt.Fprintf(r.Stdout, "%d. %s (вес: %g)\n", e.ID, e.Item.Name(), e.Item.Weight())
	}
	return 0
}

func (r *Runner) add(args []string) int {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(r.Stderr)
	name := fs.String("name", "", "item name")
	weight := fs.Float64("weight", 0, "item weight")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(r.Stderr, "error: unexpected arguments: %v\n", fs.Args())
		return 2
	}
	if err := validateName(*name); err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return 2
	}
	if err := validateWeight(*weight); err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return 2
	}

	inv, ok := r.loadOrFail()
	if !ok {
		return 1
	}
	item := &items.Weapon{Title: *name, Mass: *weight, Damage: 1, Durability: 1}
	id := inv.Add(item)
	if err := r.Store.Save(inv); err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return 1
	}
	fmt.Fprintf(r.Stdout, "added id=%d\n", id)
	return 0
}

func (r *Runner) use(args []string) int {
	fs := flag.NewFlagSet("use", flag.ContinueOnError)
	fs.SetOutput(r.Stderr)
	id := fs.Int("id", 0, "item id")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(r.Stderr, "error: unexpected arguments: %v\n", fs.Args())
		return 2
	}

	inv, ok := r.loadOrFail()
	if !ok {
		return 1
	}
	msg, err := inv.Use(*id)
	if err != nil {
		printItemError(r.Stderr, err)
		return 1
	}
	if err := r.Store.Save(inv); err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return 1
	}
	fmt.Fprintln(r.Stdout, msg)
	return 0
}

func (r *Runner) delete(args []string) int {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	fs.SetOutput(r.Stderr)
	id := fs.Int("id", 0, "item id")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(r.Stderr, "error: unexpected arguments: %v\n", fs.Args())
		return 2
	}

	inv, ok := r.loadOrFail()
	if !ok {
		return 1
	}
	if err := inv.Delete(*id); err != nil {
		printItemError(r.Stderr, err)
		return 1
	}
	if err := r.Store.Save(inv); err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return 1
	}
	fmt.Fprintf(r.Stdout, "deleted id=%d\n", *id)
	return 0
}

func validateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(name) != name {
		return errors.New("name must not have leading/trailing spaces")
	}
	if strings.ContainsAny(name, "|\n\r") {
		return errors.New(`name must not contain '|', '\n' or '\r'`)
	}
	return nil
}

func validateWeight(w float64) error {
	if math.IsNaN(w) || math.IsInf(w, 0) {
		return errors.New("weight must be a finite number")
	}
	if w <= 0 {
		return errors.New("weight must be > 0")
	}
	return nil
}

func printItemError(w io.Writer, err error) {
	var pie *items.ParseItemError
	switch {
	case errors.As(err, &pie):
		fmt.Fprintf(w, "error: parse at line %d (%q): %v\n", pie.Line, pie.Value, pie.Err)
	case errors.Is(err, items.ErrBroken):
		fmt.Fprintln(w, "error: предмет сломан")
	case errors.Is(err, items.ErrExhausted):
		fmt.Fprintln(w, "error: заряды предмета закончились")
	case errors.Is(err, inventory.ErrNotFound):
		fmt.Fprintln(w, "error: предмет не найден")
	default:
		fmt.Fprintln(w, "error:", err)
	}
}
