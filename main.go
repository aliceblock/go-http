package main

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

var staticHTTP *http.ServeMux
var routes Routes
var tpl *template.Template

type Route struct {
	Pattern string
	Fn      func(http.ResponseWriter, *http.Request)
}

type Routes struct {
	RouteList []Route
}

func (r *Routes) setRoutes(routes ...interface{}) {
	for _, route := range routes {
		router := route.(Route)
		http.HandleFunc(router.Pattern, router.Fn)
		r.RouteList = append(r.RouteList, route.(Route))
	}
}

func (r *Routes) getPatterns() []string {
	var patterns []string
	for _, route := range r.RouteList {
		patternString := strings.Replace(route.Pattern, "/", "", -1)
		if patternString != "" {
			patterns = append(patterns, patternString)
		}
	}
	return patterns
}

func (r *Routes) isValidPath(path string) []string {
	validPath := regexp.MustCompile(("^/(" + strings.Join(r.getPatterns(), "|") + ")*$"))
	return validPath.FindStringSubmatch(path)
}

func serveSingleFile(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if routes.isValidPath(req.URL.Path) == nil {
			staticHTTP.ServeHTTP(rw, req)
			return
		}
		fn(rw, req)
	}
}

func homeHandler(rw http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(rw, "index.html", nil)
}

func init() {
	staticHTTP = http.NewServeMux()

	// Pre load .html template file
	tpl = template.Must(template.ParseGlob("./*.html"))

	// Setting URL path
	routes.setRoutes(
		Route{"/", makeHandler(homeHandler)},
	)
	// Setting static files directory
	staticHTTP.Handle("/", http.FileServer(http.Dir("./public")))
}

func main() {
	http.ListenAndServe(":8080", nil)
}
