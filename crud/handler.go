package crud

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	service1 *service
}

func NewHandler(service *service) *handler {
	return &handler{service1: service}
}
func (h *handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.select1(w, r)

	}
}

func (h *handler) select1(w http.ResponseWriter, r *http.Request) {
	products, err := h.service1.select1()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *handler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.selectAllProducts(w, r)
	default:

	}
}

func (h *handler) selectAllProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	products, err := h.service1.selectAllProducts(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.selectAll(w, r)
	case http.MethodPost:
		h.insert(w, r)
	}
}

// https://m-nhasibuan5181-xe4oymdo.leapcell.dev/api/categories
// https://mmmmmm.zeabur.app/api/categories
func (h *handler) selectAll(w http.ResponseWriter, r *http.Request) {
	c, err := h.service1.selectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// [
//     {
//         "id": 1,
//         "nama": "Food",
//         "description": "makanan"
//     },
//     {
//         "id": 4,
//         "nama": "Drinks",
//         "description": "minuman"
//     }
// ]

// https://m-nhasibuan5181-xe4oymdo.leapcell.dev/api/categories
func (h *handler) insert(w http.ResponseWriter, r *http.Request) {
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service1.insert(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *handler) HandleByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.selectByID(w, r)
	case http.MethodPut:
		h.update(w, r)
	case http.MethodDelete:
		h.delete(w, r)
	}
}

// https://m-nhasibuan5181-xe4oymdo.leapcell.dev/api/categories/1
func (h *handler) selectByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.service1.selectByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// https://m-nhasibuan5181-xe4oymdo.leapcell.dev/api/categories/1
func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var category Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category.ID = id
	err = h.service1.update(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// https://m-nhasibuan5181-xe4oymdo.leapcell.dev/api/categories/1
func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service1.delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "delete successful",
	})
}

func (h *handler) InsertShoppingCart(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var co CheckoutItem1
		err := json.NewDecoder(r.Body).Decode(&co)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.service1.insertShoppingCart(&co)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(co)

	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}

func (h *handler) Checkout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req CheckoutRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		transaction, err := h.service1.checkout(req.Items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transaction)
	default:
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	}
}

func (h *handler) HandleTodayReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	values := r.URL.Query()
	startStr := values.Get("start_date")
	endStr := values.Get("end_date")
	var (
		startDate time.Time
		endDate   time.Time
		err       error
	)
	if startStr == "" || endStr == "" {
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate
		// http.Error(w, "", http.StatusBadRequest)
		// return
	} else {

		startDate, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		endDate, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	report, err := h.service1.GetTodayReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(report)
}
