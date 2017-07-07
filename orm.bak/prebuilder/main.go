package main

import (
	`flag`
	`fmt`
	`log`
)

var (
	file = flag.String(`m`, ``, `-m \path\to\model.go`)
)

func main() {
	flag.Parse()
	// fmt.Print(`hhhh`)
	ParseFile(*file)
}

func echo(i ...interface{}) {
	fmt.Println(i...)
}
func echoStrings(strs ...[]string) {
	for i, val := range strs {
		fmt.Print(`[`, i, `] `)
		for j, v := range val {
			fmt.Println(j, `: `, v)
		}
		fmt.Print("\n")
	}
}
func echoBytes(args interface{}) {
	switch v := args.(type) {
	case []byte:
		echo(string(v))
	case []rune:
		echo(string(v))
	case [][]byte:
		for i, val := range v {
			fmt.Print(i, `: `, string(val), ` `)
		}
		fmt.Print("\n")
	case [][]rune:
		for i, val := range v {
			fmt.Print(i, `: `, string(val), ` `)
		}
		fmt.Print("\n")
	case [][][]byte:
		for i, val := range v {
			fmt.Print(i, ` `)
			echoBytes(val)
		}
	case [][][]rune:
		for i, val := range v {
			fmt.Print(i, ` `)
			echoBytes(val)
		}
	default:
		echo(v)
	}

}

func logit(data ...interface{}) {
	log.Println(data...)
}
