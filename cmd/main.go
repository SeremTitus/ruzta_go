package main

import (
	"ruzta/pkg/tokenizer"
)

func main() {
	newTokenizer := tokenizer.NewTokenizer(`
mod Demo {
    // Single-line comment
    # Alternate line comment

    /* Block comment */
    class Bar {
        var r = 0
    }

    class Foo extends Bar {
        var x = 1
        fn add(a, b) Int {
            return a + b
        }
    }

    fn main() {
        var y = 1
        var z = 2
        if (y < z) {
            y = y + 1
        }
        match (y) {
            1, 2 when z < 3 { return }
            _: { return }
        }
    }
}
`)
	for {
		newToken := newTokenizer.Scan()
		println("Token: " + newToken.GetDebugName())
		if newToken.Type == tokenizer.EOF {
			break
		}
	}
}
