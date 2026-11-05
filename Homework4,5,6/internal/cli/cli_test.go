package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"homework6/internal/inventory"
	"homework6/internal/items"
)

type mockStorage struct {
	inv       *inventory.Inventory
	loadErr   error
	saveErr   error
	loadCount int // Счётчик вызовов Load
	saveCount int // Счётчик вызовов Save
}

func (m *mockStorage) Load() (*inventory.Inventory, error) {
	m.loadCount++
	return m.inv, m.loadErr
}

func (m *mockStorage) Save(inv *inventory.Inventory) error {
	m.saveCount++
	m.inv = inv
	return m.saveErr
}

func TestRun_CLI(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		mockStore  *mockStorage
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{
			name:       "no args",
			args:       []string{},
			wantCode:   2,
			wantStderr: "usage:",
		},
		{
			name:       "unknown command",
			args:       []string{"foo"},
			wantCode:   2,
			wantStderr: "unknown command",
		},
		{
			name:       "list empty",
			args:       []string{"list"},
			mockStore:  &mockStorage{inv: inventory.New()},
			wantCode:   0,
			wantStdout: "инвентарь пуст",
		},
		{
			name:       "list with items",
			args:       []string{"list"},
			mockStore:  &mockStorage{inv: inventory.New()},
			wantCode:   0,
			wantStdout: "1. Меч",
		},
		{
			name:       "add missing flags",
			args:       []string{"add"},
			wantCode:   2,
			wantStderr: "name is required",
		},
		{
			name:       "use without id",
			args:       []string{"use"},
			wantCode:   2,
			wantStderr: "-id must be > 0",
		},
		{
			name:       "use not found",
			args:       []string{"use", "-id", "99"},
			mockStore:  &mockStorage{inv: inventory.New()},
			wantCode:   1,
			wantStderr: "предмет не найден",
		},
		{
			name:       "load error",
			args:       []string{"list"},
			mockStore:  &mockStorage{loadErr: errors.New("disk fail")},
			wantCode:   1,
			wantStderr: "disk fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "list with items" {
				tt.mockStore.inv.Add(&items.Weapon{Title: "Меч", Mass: 1})
			}

			var stdout, stderr bytes.Buffer
			code := Run(tt.args, &stdout, &stderr, tt.mockStore)

			if code != tt.wantCode {
				t.Errorf("code = %d, want %d", code, tt.wantCode)
			}
			if tt.wantStdout != "" && !strings.Contains(stdout.String(), tt.wantStdout) {
				t.Errorf("stdout = %q, want contains %q", stdout.String(), tt.wantStdout)
			}
			if tt.wantStderr != "" && !strings.Contains(stderr.String(), tt.wantStderr) {
				t.Errorf("stderr = %q, want contains %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestRun_CLI_SuccessOperations(t *testing.T) {
	mock := &mockStorage{inv: inventory.New()}
	var stdout, stderr bytes.Buffer

	// 1. Тест ADD
	code := Run([]string{"add", "-name", "Меч", "-weight", "3.5"}, &stdout, &stderr, mock)
	if code != 0 {
		t.Fatalf("add failed with code %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "added id=1") {
		t.Errorf("add output wrong: %s", stdout.String())
	}
	if mock.inv.Len() != 1 {
		t.Errorf("mock inventory not updated after add")
	}
	if mock.saveCount != 1 {
		t.Errorf("Save() was not called after add, count = %d", mock.saveCount)
	}

	stdout.Reset()
	stderr.Reset()
	mock.saveCount = 0 // сбрасываем счётчик перед следующим вызовом

	code = Run([]string{"use", "-id", "1"}, &stdout, &stderr, mock)
	if code != 0 {
		t.Fatalf("use failed with code %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Используется Меч") {
		t.Errorf("use output wrong: %s", stdout.String())
	}
	if mock.saveCount != 1 {
		t.Errorf("Save() was not called after use, count = %d", mock.saveCount)
	}

	stdout.Reset()
	stderr.Reset()
	mock.saveCount = 0

	code = Run([]string{"delete", "-id", "1"}, &stdout, &stderr, mock)
	if code != 0 {
		t.Fatalf("delete failed with code %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "deleted id=1") {
		t.Errorf("delete output wrong: %s", stdout.String())
	}
	if mock.inv.Len() != 0 {
		t.Errorf("mock inventory not updated after delete")
	}
	if mock.saveCount != 1 {
		t.Errorf("Save() was not called after delete, count = %d", mock.saveCount)
	}
}

func TestRun_CLI_Validation(t *testing.T) {
	mock := &mockStorage{inv: inventory.New()}
	var stdout, stderr bytes.Buffer

	code := Run([]string{"add", "-name", "Меч", "-weight", "-1"}, &stdout, &stderr, mock)
	if code != 2 {
		t.Errorf("expected code 2 for negative weight, got %d", code)
	}
}
