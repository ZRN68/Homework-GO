package cli

import (
	"flag"
	"fmt"
	"io"

	"homework6/internal/domain"
	"homework6/internal/storage"
)

type Runner struct {
	Store  storage.Storage
	Stdout io.Writer
	Stderr io.Writer
}

func Run(args []string, stdout, stderr io.Writer, store storage.Storage) int {
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
		fmt.Fprintf(stderr, "unknown command %q\n", cmd)
		return 2
	}
}

func (r *Runner) list(args []string) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(r.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	inv, err := r.Store.Load()
	if err != nil {
		fmt.Fprintln(r.Stderr, "warning:", err)
	}
	if len(inv.Items) == 0 {
		fmt.Fprintln(r.Stdout, "инвентарь пуст")
		return 0
	}
	for i, it := range inv.Items {
		fmt.Fprintf(r.Stdout, "%d. %s (вес: %g)\n", i+1, it.Name(), it.Weight())
	}
	return 0
}

func (r *Runner) add(args []string) int {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(r.Stderr)
	name := fs.String("name", "", "имя предмета")
	weight := fs.Float64("weight", 0, "вес")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *name == "" {
		fmt.Fprintln(r.Stderr, "error: -name is required")
		return 2
	}
	if *weight <= 0 {
		fmt.Fprintln(r.Stderr, "error: -weight must be > 0")
		return 2
	}

	inv, err := r.Store.Load()
	if err != nil {
		fmt.Fprintln(r.Stderr, "warning:", err)
	}
	id := inv.Add(&domain.Weapon{Title: *name, Mass: *weight, Damage: 1, Durability: 1})
	if err := r.Store.Save(inv); err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return 1
	}
	fmt.Fprintf(r.Stdout, "added id=%d %s\n", id, *name)
	return 0
}

func (r *Runner) use(args []string) int {
	fs := flag.NewFlagSet("use", flag.ContinueOnError)
	fs.SetOutput(r.Stderr)
	id := fs.Int("id", 0, "id предмета")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *id <= 0 {
		fmt.Fprintln(r.Stderr, "error: -id is required")
		return 2
	}

	inv, err := r.Store.Load()
	if err != nil {
		fmt.Fprintln(r.Stderr, "warning:", err)
	}
	msg, err := inv.Use(*id)
	if err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
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
	id := fs.Int("id", 0, "id предмета")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *id <= 0 {
		fmt.Fprintln(r.Stderr, "error: -id is required")
		return 2
	}

	inv, err := r.Store.Load()
	if err != nil {
		fmt.Fprintln(r.Stderr, "warning:", err)
	}
	if err := inv.Delete(*id); err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return 1
	}
	if err := r.Store.Save(inv); err != nil {
		fmt.Fprintln(r.Stderr, "error:", err)
		return 1
	}
	fmt.Fprintf(r.Stdout, "deleted id=%d\n", *id)
	return 0
}
