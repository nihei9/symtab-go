package symtab

import "fmt"

func ExampleSymbolTable() {
	tab := NewSymbolTable()
	w := tab.Writer()
	sym, err := w.Add("foo", 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := tab.Reader()
	if str, ok := r.ToString(sym); ok {
		fmt.Println(str)
	}

	// Output:
	// foo
}
