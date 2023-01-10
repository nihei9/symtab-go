package symtab

import (
	"errors"
	"strconv"
	"testing"
)

func TestWriter(t *testing.T) {
	const (
		FLAG_LOW  uint16 = 0x00ff
		FLAG_HIGH uint16 = 0xff00
	)

	t.Run("A writer can write a string", func(t *testing.T) {
		tab := NewSymbolTable()
		w := tab.Writer()
		sym, err := w.Add("aaa", FLAG_HIGH|FLAG_LOW)
		if err != nil {
			t.Fatal(err)
		}
		if sym.Num() != 1 || sym.Flags() != FLAG_HIGH|FLAG_LOW {
			t.Fatalf("unexpected symbol: %v", sym)
		}
	})

	t.Run("A writer can write another string", func(t *testing.T) {
		tab := NewSymbolTable()
		w := tab.Writer()
		_, err := w.Add("aaa", FLAG_HIGH|FLAG_LOW)
		if err != nil {
			t.Fatal(err)
		}
		sym, err := w.Add("bbb", FLAG_HIGH|FLAG_LOW)
		if err != nil {
			t.Fatal(err)
		}
		if sym.Num() != 2 || sym.Flags() != FLAG_HIGH|FLAG_LOW {
			t.Fatalf("unexpected symbol: %v", sym)
		}
	})

	t.Run("A writer cannot write an empty string", func(t *testing.T) {
		tab := NewSymbolTable()
		w := tab.Writer()
		sym, err := w.Add("", FLAG_HIGH|FLAG_LOW)
		if !errors.Is(err, errEmptyString) {
			t.Fatalf("unexpected error occurred. want: %v, got: %v", errEmptyString, err)
		}
		if sym != ZeroSymbol {
			t.Fatalf("symbol should be zero value: %v", sym)
		}
	})

	t.Run("A writer cannot write duplicated strings", func(t *testing.T) {
		tab := NewSymbolTable()
		w := tab.Writer()
		_, err := w.Add("aaa", FLAG_HIGH|FLAG_LOW)
		if err != nil {
			t.Fatal(err)
		}
		sym, err := w.Add("aaa", FLAG_HIGH|FLAG_LOW)
		if !errors.Is(err, errStringAlreadyAdded) {
			t.Fatalf("unexpected error occurred. want: %v, got: %v", errStringAlreadyAdded, err)
		}
		if sym != ZeroSymbol {
			t.Fatalf("symbol should be zero value: %v", sym)
		}
	})

	t.Run("A writer raises an error when the maximum number of symbols has been reached", func(t *testing.T) {
		tab := NewSymbolTable()
		maxNum = 100
		w := tab.Writer()
		for n := 0; n < 100; n++ {
			_, err := w.Add(strconv.Itoa(n), FLAG_HIGH|FLAG_LOW)
			if err != nil {
				t.Fatal(err)
			}
		}
		sym, err := w.Add("aaa", FLAG_HIGH|FLAG_LOW)
		if !errors.Is(err, errSymbolCountLimitReached) {
			t.Fatalf("unexpected error occurred. want: %v, got: %v", errSymbolCountLimitReached, err)
		}
		if sym != ZeroSymbol {
			t.Fatalf("symbol should be zero value: %v", sym)
		}
	})
}

func TestReader(t *testing.T) {
	const (
		FLAG_LOW  uint16 = 0x00ff
		FLAG_HIGH uint16 = 0xff00
	)

	tab := NewSymbolTable()
	w := tab.Writer()
	symA, err := w.Add("aaa", FLAG_HIGH|FLAG_LOW)
	if err != nil {
		t.Fatal(err)
	}
	symB, err := w.Add("bbb", FLAG_HIGH|FLAG_LOW)
	if err != nil {
		t.Fatal(err)
	}

	r := tab.Reader()

	t.Run("A reader can read registered symbols", func(t *testing.T) {
		str, ok := r.ToString(symA)
		if !ok || str != "aaa" {
			t.Fatalf("unexpected result. want: (\"aaa\", true), got: (%#v, %v)", str, ok)
		}
		str, ok = r.ToString(symB)
		if !ok || str != "bbb" {
			t.Fatalf("unexpected result. want: (\"bbb\", true), got: (%#v, %v)", str, ok)
		}
	})

	t.Run("A reader returns `false` when a symbol is the zero value", func(t *testing.T) {
		str, ok := r.ToString(ZeroSymbol)
		if ok || str != "" {
			t.Fatalf("unexpected result. want: (\"\", false), got: (%#v, %v)", str, ok)
		}
	})

	t.Run("A reader can read registered strings", func(t *testing.T) {
		sym, ok := r.ToSymbol("aaa")
		if !ok || sym != symA {
			t.Fatalf("unexpected result. want: (%v, true), got: (%#v, %v)", symA, sym, ok)
		}
		sym, ok = r.ToSymbol("bbb")
		if !ok || sym != symB {
			t.Fatalf("unexpected result. want: (%v, true), got: (%#v, %v)", symB, sym, ok)
		}
	})

	t.Run("A reader returns `false` when a string is not registered", func(t *testing.T) {
		sym, ok := r.ToSymbol("zzz")
		if ok || sym != ZeroSymbol {
			t.Fatalf("unexpected result. want: (%v, false), got: (%#v, %v)", ZeroSymbol, sym, ok)
		}

		sym, ok = r.ToSymbol("")
		if ok || sym != ZeroSymbol {
			t.Fatalf("unexpected result. want: (%v, false), got: (%#v, %v)", ZeroSymbol, sym, ok)
		}
	})
}
