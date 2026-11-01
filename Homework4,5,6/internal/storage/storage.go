package storage

import "homework6/internal/inventory"

// Storage — то, что нужно CLI. Маленький интерфейс, объявлен со стороны
// потребителя (cli будет принимать его, а не конкретный тип).
type Storage interface {
	Load() (*inventory.Inventory, error)
	Save(inv *inventory.Inventory) error
}
