package router

import (
	"net/url"
	"testing"
	"github.com/KateGF/Http-Server-Project-SO/core"
)

// makeReq construye un core.HttpRequest para pruebas.
func makeReq(method, path string) *core.HttpRequest {
	u, _ := url.Parse(path)
	return core.NewHttpRequest(method, u, map[string]string{}, "")
}

func TestHandleExactPath(t *testing.T) {
	r := New()
	r.Get("/foo", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok().Text("FOO"), nil
	})
	res, err := r.Handle(makeReq("GET", "/foo"))
	if err != nil {
		t.Fatalf("Handle devolvió error: %v", err)
	}
	if res.StatusCode != 200 || res.Body != "FOO" {
		t.Errorf("Esperaba 200 FOO; obtuve %d %q", res.StatusCode, res.Body)
	}
}

func TestHandleNotFound(t *testing.T) {
	r := New()
	res, err := r.Handle(makeReq("GET", "/nope"))
	if err != nil {
		t.Fatalf("Handle devolvió error: %v", err)
	}
	if res.StatusCode != 404 {
		t.Errorf("Esperaba 404 para ruta no registrada; obtuve %d", res.StatusCode)
	}
}

func TestMethodMismatch(t *testing.T) {
	r := New()
	r.Post("/foo", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok(), nil
	})
	res, err := r.Handle(makeReq("GET", "/foo"))
	if err != nil {
		t.Fatalf("Handle devolvió error: %v", err)
	}
	if res.StatusCode != 404 {
		t.Errorf("Esperaba 404 por método incorrecto; obtuve %d", res.StatusCode)
	}
}

func TestPrefixMatchWithSlash(t *testing.T) {
	r := New()
	r.Get("/api/", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok().Text("API"), nil
	})
	res, err := r.Handle(makeReq("GET", "/api/books"))
	if err != nil {
		t.Fatalf("Handle devolvió error: %v", err)
	}
	if res.Body != "API" {
		t.Errorf("Esperaba 'API' por prefijo con '/'; obtuve %q", res.Body)
	}
}

func TestPrefixMatchWithoutSlash(t *testing.T) {
	r := New()
	r.Get("/foo", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok().Text("PREFIX"), nil
	})
	res, err := r.Handle(makeReq("GET", "/foo/bar"))
	if err != nil {
		t.Fatalf("Handle devolvió error: %v", err)
	}
	if res.Body != "PREFIX" {
		t.Errorf("Esperaba 'PREFIX' por prefijo implícito; obtuve %q", res.Body)
	}
}

func TestMatchFunction(t *testing.T) {
	cases := []struct {
		reqPath, handlerPath string
		want                 bool
	}{
		{"/a", "/a", true},
		{"/a/b", "/a", true},
		{"/a", "/a/", false},
		{"/abc", "/a", false},
		{"/api/v1", "/api/", true},
	}
	for _, tc := range cases {
		got := match(tc.reqPath, tc.handlerPath)
		if got != tc.want {
			t.Errorf("match(%q, %q) = %v; want %v", tc.reqPath, tc.handlerPath, got, tc.want)
		}
	}
}

func TestSpecificityOrdering(t *testing.T) {
	r := New()
	called := ""
	r.Get("/foo", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		called = "/foo"
		return core.Ok(), nil
	})
	r.Get("/foo/bar", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		called = "/foo/bar"
		return core.Ok(), nil
	})
	// Sin ordenar, el registro primero ganaría; con SortHandlers se invierte
	r.SortHandlers()
	_, err := r.Handle(makeReq("GET", "/foo/bar"))
	if err != nil {
		t.Fatalf("Handle devolvió error: %v", err)
	}
	if called != "/foo/bar" {
		t.Errorf("Esperaba handler '/foo/bar', llamado %q", called)
	}
}

func TestSortHandlersOrder(t *testing.T) {
	r := New()
	h := func(req *core.HttpRequest) (*core.HttpResponse, error) { return core.Ok(), nil }
	r.Get("/a", h)
	r.Get("/a/b/c", h)
	r.Get("/a/b", h)
	r.SortHandlers()
	paths := []string{r.routes[0].Path, r.routes[1].Path, r.routes[2].Path}
	want := []string{"/a/b/c", "/a/b", "/a"}
	for i := range want {
		if paths[i] != want[i] {
			t.Errorf("SortHandlers order: en %d, got %q; want %q", i, paths[i], want[i])
		}
	}
}

func TestDeleteAndPostRoutes(t *testing.T) {
	r := New()
	r.Delete("/item", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok().Text("DEL"), nil
	})
	r.Post("/item", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok().Text("POST"), nil
	})
	resDel, _ := r.Handle(makeReq("DELETE", "/item"))
	if resDel.Body != "DEL" {
		t.Errorf("Esperaba 'DEL'; obtuve %q", resDel.Body)
	}
	resPost, _ := r.Handle(makeReq("POST", "/item"))
	if resPost.Body != "POST" {
		t.Errorf("Esperaba 'POST'; obtuve %q", resPost.Body)
	}
}
