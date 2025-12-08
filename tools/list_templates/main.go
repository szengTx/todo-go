package main

import (
	"fmt"
	"html/template"
	"os"
)

func main() {
	t := template.Must(template.ParseGlob("templates/*.html"))
	fmt.Println("Defined templates:")
	for _, tmpl := range t.Templates() {
		fmt.Println(" -", tmpl.Name())
	}
	os.Exit(0)
}
