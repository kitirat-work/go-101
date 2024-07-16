package exportlargeexcel

import (
	"net/http"
	"time"
)

type IHandler interface {
	Export(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	service IService
}

// Export implements IHandler.
func (h *Handler) Export(w http.ResponseWriter, r *http.Request) {
	// success
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	date := time.Now().Format("02-01-2006")
	w.Header().Set("Content-Disposition", "attachment; filename=report-database-"+date+".xlsx")
	err := h.service.Export(r.Context(), w)
	if err != nil {
		// internal server error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func NewHandler(iService IService) IHandler {
	return &Handler{
		service: iService,
	}
}
