package reconciliation
import (
 "encoding/json"
 "net/http"
 pkghttp "apexpay/internal/platform/http"
 mw "apexpay/internal/platform/middleware"
 "github.com/go-chi/chi/v5"
)
type Handler struct{ svc *Service }
func NewHandler(s *Service)*Handler{return &Handler{svc:s}}
func(h *Handler) Routes(r chi.Router){r.Get("/payment-reconciliation",h.List);r.Post("/payment-reconciliation/{merchantID}/{key}/decision",h.Decide)}
func(h *Handler) List(w http.ResponseWriter,r *http.Request){cases,err:=h.svc.ListOpen(r.Context(),50);if err!=nil{pkghttp.WriteError(w,r,err);return};pkghttp.WriteJSON(w,r,200,cases)}
func(h *Handler) Decide(w http.ResponseWriter,r *http.Request){var req struct{Decision string `json:"decision"`; Note string `json:"note"`};if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{pkghttp.WriteErrorWithBody(w,r,400,"validation_error","invalid json");return}; reviewer:=mw.APIKeyID(r.Context());if err:=h.svc.Decide(r.Context(),chi.URLParam(r,"merchantID"),chi.URLParam(r,"key"),req.Decision,reviewer,req.Note);err!=nil{pkghttp.WriteErrorWithBody(w,r,409,"reconciliation_conflict",err.Error());return};pkghttp.WriteJSON(w,r,200,map[string]string{"status":req.Decision})}
