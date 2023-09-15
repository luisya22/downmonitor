package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/luisya22/downmonitor/monitor"
)

func Start() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}

		if r.URL.Path != "/" {
			http.NotFound(w, r)
		}

		file, err := os.Open("./downmonitor.log")
		if err != nil {
			fmt.Printf("error opening file: %v", err.Error())
			return
		}
		defer file.Close()

		data, err := monitor.QueryData(file, &monitor.RealTime{})
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			http.Error(w, "error querying data", http.StatusInternalServerError)
			return
		}

		fmt.Printf("data: %v", data)

		tmplFilePath := "./templates/html/template.gohtml"

		tmpl, err := template.ParseFiles(tmplFilePath)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			http.Error(w, "error with template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			http.Error(w, "error executing template", http.StatusInternalServerError)
			return
		}

	})

	http.ListenAndServe(":8720", nil)
}
