package lines_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"gonih.org/lines"
)

func Example_unmarshal() {
	type X struct {
		Foo int
		Bar bool
	}

	r := strings.NewReader(`{
		"foo": 42,
		"bar": false,
}`)
	lr := lines.NewReader(r)
	buf, err := io.ReadAll(lr)
	if err != nil {
		log.Fatal(err)
	}
	var x X
	err = json.Unmarshal(buf, &x)

	var se *json.SyntaxError
	if errors.As(err, &se) {
		l, c := lr.Position(se.Offset)
		fmt.Printf("Syntax error at line %d column %d: %v\n", l, c, err)
	}
	// Output: Syntax error at line 4 column 2: invalid character '}' looking for beginning of object key string
}

func Example_decoder() {
	type X struct {
		Foo int
		Bar bool
	}

	r := strings.NewReader(`{
		"foo": 42,
		"bar": false
}`)
	lr := lines.NewReader(r)
	d := json.NewDecoder(lr)
	for {
		t, err := d.Token()
		if err != nil {
			break
		}
		l, c := lr.Position(d.InputOffset())
		fmt.Printf("%d:%d: %T(%v)\n", l, c, t, t)
	}
	// Output:
	// 1:2: json.Delim({)
	// 2:8: string(foo)
	// 2:12: float64(42)
	// 3:8: string(bar)
	// 3:15: bool(false)
	// 4:2: json.Delim(})
}
