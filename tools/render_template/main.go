package main

import (
	"html/template"
	"os"
	"todo-go/database"
)

func main() {
	database.InitDB()
	t := template.Must(template.ParseGlob("templates/*.html")) // load templates
	// render register.html
	t.ExecuteTemplate(os.Stdout, "register.html", nil)
}
