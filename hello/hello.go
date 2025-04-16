package hello

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func Hello() {
	fmt.Println(reverse.String("Hello"))
}
