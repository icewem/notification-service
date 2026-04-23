package handler

import (
	"net/http"
	_ "net/http/pprof" // регистрирует /debug/pprof endpoints
)

// RegisterDebugHandlers — регистрирует pprof эндпоинты
func RegisterDebugHandlers(mux *http.ServeMux) {
	// /debug/pprof/          — список профилей
	// /debug/pprof/goroutine — все горутины
	// /debug/pprof/heap      — heap профиль
	// /debug/pprof/profile   — CPU профиль (30 сек)
	// /debug/pprof/allocs    — все аллокации
	mux.HandleFunc("/debug/pprof/", http.DefaultServeMux.ServeHTTP)
	mux.HandleFunc("/debug/pprof/cmdline", http.DefaultServeMux.ServeHTTP)
	mux.HandleFunc("/debug/pprof/profile", http.DefaultServeMux.ServeHTTP)
	mux.HandleFunc("/debug/pprof/symbol", http.DefaultServeMux.ServeHTTP)
	mux.HandleFunc("/debug/pprof/trace", http.DefaultServeMux.ServeHTTP)
}
