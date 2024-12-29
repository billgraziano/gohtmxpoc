package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"pochtmx/static"
	"pochtmx/web"
	"time"
)

func init() {
	// Windows 10 messes up the MIME types for these two ... sometimes.
	// var err error
	// log.Println("setting mime types...")
	// err = mime.AddExtensionType(".js", "text/javascript")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = mime.AddExtensionType(".css", "text/css")
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

var useLocal bool

func Run() {
	var err error

	flag.BoolVar(&useLocal, "local", false, "use local templates")
	flag.Parse()
	log.Printf("use local: %v", useLocal)
	web.Local = useLocal

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(static.HttpFS(useLocal))))
	mux.HandleFunc("GET /{$}", indexPage)
	mux.HandleFunc("POST /search", searchPartial)
	mux.HandleFunc("GET /api/employees", searchAPI)

	srv := &http.Server{
		Addr:              "localhost:8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       300 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("url: http://localhost:8080")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	err := web.Execute(w, nil, "templates/base.gohtml", "templates/index.gohtml")
	if err != nil {
		log.Println(err)
		http.Error(w, "error rendering page", http.StatusInternalServerError)
	}
}

func searchAPI(w http.ResponseWriter, r *http.Request) {
	var Result struct {
		Employees []Employee `json:"employees"`
	}
	employees, _, err := searchHelper(r)
	if err != nil {
		log.Println(err)
		err = encode(w, http.StatusInternalServerError, Result)
		if err != nil {
			log.Println(err)
		}
		return
	}

	Result.Employees = employees
	err = encode(w, http.StatusInternalServerError, Result)
	if err != nil {
		log.Println(err)
	}
}

func searchPartial(w http.ResponseWriter, r *http.Request) {
	employees, term, err := searchHelper(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var Result struct {
		Employees  []Employee
		SearchTerm string
	}
	Result.Employees = employees
	Result.SearchTerm = term
	err = web.Execute(w, Result, "partials/search.gohtml")
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func searchHelper(r *http.Request) ([]Employee, string, error) {
	var err error
	err = r.ParseForm()
	if err != nil {
		return []Employee{}, "", fmt.Errorf("parseform: %w", err)
	}
	txt := r.FormValue("search")
	employees, err := GetEmployees(txt)
	if err != nil {
		return []Employee{}, txt, fmt.Errorf("getemployees: %w", err)
	}
	return employees, txt, nil
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("json.encode: %w", err)
	}
	return nil
}
