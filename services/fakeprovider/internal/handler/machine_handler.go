package handler

import (
"encoding/json"
"net/http"

"github.com/go-chi/chi/v5"
"github.com/showbiz-io/showbiz/services/fakeprovider/internal/service"
)

type MachineHandler struct {
svc *service.MachineService
}

func NewMachineHandler(svc *service.MachineService) *MachineHandler {
return &MachineHandler{svc: svc}
}

func (h *MachineHandler) Create(w http.ResponseWriter, r *http.Request) {
var input service.CreateMachineInput
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
return
}

if input.Name == "" || input.Namespace == "" || input.CPU <= 0 || input.MemoryMB <= 0 || input.Image == "" {
Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "name, namespace, cpu, memoryMB, and image are required")
return
}

m := h.svc.Create(input)
JSON(w, http.StatusCreated, m)
}

func (h *MachineHandler) List(w http.ResponseWriter, r *http.Request) {
machines := h.svc.List()
JSON(w, http.StatusOK, machines)
}

func (h *MachineHandler) Get(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")
m, err := h.svc.Get(id)
if err != nil {
Error(w, http.StatusNotFound, "NOT_FOUND", "Machine not found")
return
}
JSON(w, http.StatusOK, m)
}

func (h *MachineHandler) Update(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")
var input service.UpdateMachineInput
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
return
}

m, err := h.svc.Update(id, input)
if err != nil {
Error(w, http.StatusNotFound, "NOT_FOUND", "Machine not found")
return
}
JSON(w, http.StatusOK, m)
}

func (h *MachineHandler) Delete(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")
if err := h.svc.Delete(id); err != nil {
Error(w, http.StatusNotFound, "NOT_FOUND", "Machine not found")
return
}
w.WriteHeader(http.StatusNoContent)
}
