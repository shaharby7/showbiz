package main

import (
"log/slog"
"net/http"
"os"

"github.com/go-chi/chi/v5"
chimw "github.com/go-chi/chi/v5/middleware"
"github.com/showbiz-io/showbiz/services/fakeprovider/internal/handler"
"github.com/showbiz-io/showbiz/services/fakeprovider/internal/kubevirt"
"github.com/showbiz-io/showbiz/services/fakeprovider/internal/service"
)

func main() {
port := envOrDefault("FAKEPROVIDER_PORT", "8081")
kubeconfig := os.Getenv("FAKEPROVIDER_KUBECONFIG")

kvClient, err := kubevirt.NewClient(kubeconfig)
if err != nil {
slog.Error("failed to create kubevirt client", "error", err)
os.Exit(1)
}

machineSvc := service.NewMachineService(kvClient)
machineHandler := handler.NewMachineHandler(machineSvc)

r := chi.NewRouter()
r.Use(chimw.Logger)
r.Use(chimw.Recoverer)
r.Use(contentTypeJSON)

r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
handler.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
})

r.Route("/v1/machines", func(r chi.Router) {
r.Post("/", machineHandler.Create)
r.Get("/", machineHandler.List)
r.Get("/{id}", machineHandler.Get)
r.Put("/{id}", machineHandler.Update)
r.Delete("/{id}", machineHandler.Delete)
})

slog.Info("starting fakeprovider", "port", port)
if err := http.ListenAndServe(":"+port, r); err != nil {
slog.Error("server failed", "error", err)
os.Exit(1)
}
}

func contentTypeJSON(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
next.ServeHTTP(w, r)
})
}

func envOrDefault(key, fallback string) string {
if v := os.Getenv(key); v != "" {
return v
}
return fallback
}
