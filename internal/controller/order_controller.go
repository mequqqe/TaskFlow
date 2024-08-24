package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"taskflow/internal/delivery"
	"taskflow/internal/domain"
	"time"
)

type OrderController struct {
	delivery *delivery.OrderService
}

func NewOrderController(service *delivery.OrderService) *OrderController {
	return &OrderController{delivery: service}
}

func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) 

	var order domain.Order
	order.EntrepreneurName = r.FormValue("entrepreneur_name")
	order.Theme = r.FormValue("theme")
	order.Amount, _ = strconv.ParseFloat(r.FormValue("amount"), 64)
	order.Requirements = r.FormValue("requirements")
	deadline, _ := time.Parse("2006-01-02", r.FormValue("deadline"))
	order.Deadline = deadline
	order.Status = r.FormValue("status")


	file, fileHeader, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		order.FileData = fileBytes
		order.FileName = fileHeader.Filename
		order.FileType = fileHeader.Header.Get("Content-Type")
	}

	if err := c.delivery.CreateOrder(&order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
