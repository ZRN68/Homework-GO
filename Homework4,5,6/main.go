package main

import (
	"fmt"
	"os"

	"homework6/internal/cli"
	"homework6/internal/inventory"
	"homework6/internal/items"
	"homework6/internal/storage"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-demo" {
		runDemo()
		return
	}

	store := storage.NewFile("inventory.txt")
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, store))
}

type panicky struct{}

func (p *panicky) Name() string         { return "Паникёр" }
func (p *panicky) Weight() float64      { return 1 }
func (p *panicky) Use() (string, error) { panic("boom") }
func (p *panicky) DescribeItem() string { return "Паникёр" }
func (p *panicky) Serialize() string    { return "panicky||1" }

func runDemo() {
	inv := inventory.New()

	wID := inv.Add(&items.Weapon{Title: "Меч", Mass: 3.5, Damage: 10, Durability: 2})
	pID := inv.Add(&panicky{})
	inv.Add(&items.Armor{Title: "Щит", Mass: 8, Defense: 15})
	inv.Add(&items.Potion{Title: "Зелье", Mass: 0.5, Effect: "HP", Charges: 1})

	fmt.Println("=== 1. Use предмета ===")
	if msg, err := inv.UseDangerousItem(wID); err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("ok:", msg)
	}

	fmt.Println("\n=== 2. UseDangerousItem: паника (recover) ===")
	if msg, err := inv.UseDangerousItem(pID); err != nil {
		fmt.Println("err:", err) // panic during use id 2: boom
	} else {
		fmt.Println("ok:", msg)
	}

	fmt.Println("\n=== 3. UseDangerousItem: несуществующий id ===")
	if _, err := inv.UseDangerousItem(999); err != nil {
		fmt.Println("err:", err) // id 999: item not found
	}

	fmt.Println("\n=== 4. Weapons ===")
	for _, it := range inv.Weapons() {
		fmt.Println(" -", it.Name())
	}

	fmt.Println("\n=== 5. BrokenItems ===")
	for i := 0; i < 5; i++ {
		_, _ = inv.UseDangerousItem(wID)
	}
	for _, it := range inv.BrokenItems() {
		fmt.Println(" -", it.Name())
	}

	fmt.Println("\n=== 6. ItemNames ===")
	fmt.Println(inv.ItemNames())

	fmt.Println("\n=== 7. FindItemByName(\"Щит\") ===")
	if it, ok := inv.FindItemByName("Щит"); ok {
		fmt.Println("найден:", it.Name())
	} else {
		fmt.Println("не найден")
	}
}
