package symtab

import (
	"errors"
	"fmt"
	"math"
)

const (
	flagLen        = 16
	numLen         = 64 - flagLen
	numMask uint64 = math.MaxUint64 >> flagLen
	initNum int64  = 1

	ZeroSymbol Symbol = 0
	MaxNum     int64  = math.MaxUint64 >> flagLen
)

// maxNum is defined as a variable, not a constant, to allow it to be changed during testing.
var maxNum = MaxNum

var (
	errEmptyString             = errors.New("a symbol table cannot contain an empty string")
	errSymbolCountLimitReached = errors.New("no more symbols can be issued because the maximum number of symbols has been reached")
	errStringAlreadyAdded      = errors.New("a string is already added")
)

type Symbol uint64

func newSymbol(num int64, flags uint16) (Symbol, error) {
	if num > maxNum {
		return 0, errSymbolCountLimitReached
	}
	return Symbol(uint64(flags)<<numLen | uint64(num)), nil
}

func (s Symbol) String() string {
	flags := s.Flags()
	return fmt.Sprintf("#%v (%x, %b)", s.Num(), flags, flags)
}

func (s Symbol) Flags() uint16 {
	return uint16(s >> numLen)
}

func (s Symbol) Num() int64 {
	return int64(uint64(s) & numMask)
}

type Reader interface {
	ToString(Symbol) (string, bool)
	ToSymbol(string) (Symbol, bool)
}

var _ Reader = &reader{}

type Writer interface {
	Add(string, uint16) (Symbol, error)
}

var _ Writer = &writer{}

type SymbolTable struct {
	str2Sym map[string]Symbol
	sym2Str map[Symbol]string
	nextNum int64
}

type reader struct {
	tab *SymbolTable
}

func (r *reader) ToString(sym Symbol) (string, bool) {
	return r.tab.toString(sym)
}

func (r *reader) ToSymbol(str string) (Symbol, bool) {
	return r.tab.toSymbol(str)
}

type writer struct {
	tab *SymbolTable
}

func (w *writer) Add(str string, flags uint16) (Symbol, error) {
	return w.tab.add(str, flags)
}

func NewSymbolTable() *SymbolTable {
	const mapCap = 10000
	return &SymbolTable{
		str2Sym: make(map[string]Symbol, mapCap),
		sym2Str: make(map[Symbol]string, mapCap),
		nextNum: initNum,
	}
}

func (t *SymbolTable) Reader() Reader {
	return &reader{
		tab: t,
	}
}

func (t *SymbolTable) Writer() Writer {
	return &writer{
		tab: t,
	}
}

func (t *SymbolTable) add(str string, flags uint16) (Symbol, error) {
	if str == "" {
		return ZeroSymbol, errEmptyString
	}
	if _, added := t.str2Sym[str]; added {
		return ZeroSymbol, fmt.Errorf("failed to add %#v: %w", str, errStringAlreadyAdded)
	}
	sym, err := newSymbol(t.nextNum, flags)
	if err != nil {
		return ZeroSymbol, err
	}
	t.nextNum++
	t.str2Sym[str] = sym
	t.sym2Str[sym] = str
	return sym, nil
}

func (t *SymbolTable) toString(sym Symbol) (string, bool) {
	str, ok := t.sym2Str[sym]
	return str, ok
}

func (t *SymbolTable) toSymbol(str string) (Symbol, bool) {
	sym, ok := t.str2Sym[str]
	return sym, ok
}
