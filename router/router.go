package router

import (
	"github.com/KateGF/Http-Server-Project-SO/core"
	"sort"
	"strings"
)

// Handle es el tipo de función que atiende una petición.
type Handle = func(*core.HttpRequest) (*core.HttpResponse, error)

// Route almacena un método, ruta y su handler.
type Route struct {
	Method string
	Path   string
	Handle Handle
}

// Router mantiene la lista de rutas.
type Router struct {
	routes []Route
}

// New crea un Router vacío.
func New() *Router {
	return &Router{routes: make([]Route, 0)}
}

// Get registra una ruta GET.
func (r *Router) Get(path string, h Handle) {
	r.routes = append(r.routes, Route{"GET", path, h})
}

// Post, Delete… (idéntico a Get, cambiando Method)
func (r *Router) Post(path string, h Handle) {
	r.routes = append(r.routes, Route{"POST", path, h})
}
func (r *Router) Delete(path string, h Handle) {
	r.routes = append(r.routes, Route{"DELETE", path, h})
}

// match comprueba coincidencia exacta o prefijo (si path termina en '/').
func match(reqPath, routePath string) bool {
	if reqPath == routePath {
		return true
	}
	if !strings.HasSuffix(routePath, "/") {
		routePath += "/"
	}
	return strings.HasPrefix(reqPath, routePath)
}

// Handle despacha la HttpRequest al handler adecuado o retorna 404.
func (r *Router) Handle(req *core.HttpRequest) (*core.HttpResponse, error) {
	for _, rt := range r.routes {
		if req.Method == rt.Method && match(req.Target.Path, rt.Path) {
			return rt.Handle(req)
		}
	}
	return core.NotFound().Text("no route"), nil
}

// SortHandlers ordena las rutas por número de segmentos (‘/’) y luego por longitud.
func (r *Router) SortHandlers() {
	sort.Slice(r.routes, func(i, j int) bool {
		ic := strings.Count(r.routes[i].Path, "/")
		jc := strings.Count(r.routes[j].Path, "/")
		if ic != jc {
			return ic > jc
		}
		return len(r.routes[i].Path) > len(r.routes[j].Path)
	})
}
