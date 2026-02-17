package web

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"crypto/subtle"

	"github.com/berkaycubuk/subtrack/internal/utils"
)

type pageData struct {
	Title string
	Error string
	Data  any
}

func (s *Server) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	// If already logged in, redirect to dashboard
	if cookie, err := r.Cookie("session"); err == nil && s.sessions.validateSession(cookie.Value) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl := parseTemplate("login.html")
	tmpl.Execute(w, pageData{Title: "Login"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(s.username)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(s.password)) == 1

	if !usernameMatch || !passwordMatch {
		tmpl := parseTemplate("login.html")
		tmpl.Execute(w, pageData{Title: "Login", Error: "Invalid username or password"})
		return
	}

	token, err := s.sessions.createSession()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(24 * time.Hour / time.Second),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		s.sessions.destroySession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	subs, err := s.subSvc.ListSubscriptions()
	if err != nil {
		log.Printf("Error listing subscriptions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl := parseTemplate("dashboard.html")
	tmpl.Execute(w, pageData{Title: "Dashboard", Data: subs})
}

func (s *Server) handleAddForm(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplate("form.html")
	tmpl.Execute(w, pageData{Title: "Add Subscription"})
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	price := r.FormValue("price")
	currency := r.FormValue("currency")
	cycle := r.FormValue("cycle")
	paymentDate := r.FormValue("payment_date")

	if err := s.subSvc.AddSubscription(name, price, currency, cycle, paymentDate); err != nil {
		tmpl := parseTemplate("form.html")
		tmpl.Execute(w, pageData{Title: "Add Subscription", Error: err.Error(), Data: map[string]string{
			"Name": name, "Price": price, "Currency": currency, "Cycle": cycle, "PaymentDate": paymentDate,
		}})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	sub, err := s.subSvc.GetSubscription(uint(id))
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	tmpl := parseTemplate("form.html")
	tmpl.Execute(w, pageData{Title: "Edit Subscription", Data: map[string]string{
		"ID":          strconv.FormatUint(uint64(sub.ID), 10),
		"Name":        sub.Name,
		"Price":       strconv.FormatFloat(sub.Price, 'f', 2, 64),
		"Currency":    sub.Currency,
		"Cycle":       sub.Cycle,
		"PaymentDate": utils.FormatDate(sub.PaymentDate),
	}})
}

func (s *Server) handleEdit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	price := r.FormValue("price")
	currency := r.FormValue("currency")
	cycle := r.FormValue("cycle")
	paymentDate := r.FormValue("payment_date")

	if err := s.subSvc.UpdateSubscription(uint(id), name, price, currency, cycle, paymentDate); err != nil {
		tmpl := parseTemplate("form.html")
		tmpl.Execute(w, pageData{Title: "Edit Subscription", Error: err.Error(), Data: map[string]string{
			"ID": strconv.FormatUint(id, 10), "Name": name, "Price": price, "Currency": currency, "Cycle": cycle, "PaymentDate": paymentDate,
		}})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleDeleteForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	sub, err := s.subSvc.GetSubscription(uint(id))
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	tmpl := parseTemplate("delete.html")
	tmpl.Execute(w, pageData{Title: "Delete Subscription", Data: sub})
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := s.subSvc.DeleteSubscription(uint(id)); err != nil {
		http.Error(w, "Failed to delete subscription", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
