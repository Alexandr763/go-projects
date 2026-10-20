package httpapi

import "net/http"

import "gamification/internal/transport/http/apigen"

func registerOtherMethods(mux *http.ServeMux) {
	routes := []struct {
		path    string
		allowed map[string]struct{}
	}{
		{path: "/api/v1/auth/register", allowed: set("POST")},
		{path: "/api/v1/auth/login", allowed: set("POST")},
		{path: "/api/v1/auth/refresh", allowed: set("POST")},
		{path: "/api/v1/auth/logout", allowed: set("POST")},
		{path: "/api/v1/tasks", allowed: set("GET", "POST")},
		{path: "/api/v1/tasks/{id}", allowed: set("GET", "PATCH", "DELETE")},
		{path: "/api/v1/tasks/{id}/complete", allowed: set("POST")},
		{path: "/api/v1/profile", allowed: set("GET")},
		{path: "/api/v1/profile/transactions", allowed: set("GET")},
		{path: "/health/live", allowed: set("GET")},
		{path: "/health/ready", allowed: set("GET")},
	}
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	for _, route := range routes {
		for _, method := range methods {
			if _, ok := route.allowed[method]; ok {
				continue
			}
			mux.HandleFunc(method+" "+route.path, methodNotAllowed)
		}
	}
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusMethodNotAllowed, apigen.ErrorBodyErrorCodeMethodNotAllowed, "method not allowed")
}

func set(methods ...string) map[string]struct{} {
	out := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		out[method] = struct{}{}
	}
	return out
}
