package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5

	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"both positive", 5, 2, 3},
		{"positive minus zero", 7, 0, 7},
		{"negative minus positive", -2, 3, -5},
		{"both negative", -2, -3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := Divide(10, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 5 {
			t.Errorf("Divide(10, 2) = %d; want %d", got, 5)
		}
	})

	t.Run("divide by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
