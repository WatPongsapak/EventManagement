package handler

import (
	"strconv"
	"database/sql"
	"event-management/activity"
	"event-management/date"
	"event-management/pinactivity"
	"html/template"
	"net/http"
	"fmt"
	
	"github.com/gorilla/mux"
)

type Handler struct {
	actManage *activity.Manager
	dateManage *date.Manager
	pinManage *pinactivity.Manager
}

func (h *Handler)indexPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	act, err := h.actManage.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t,err := template.ParseFiles("html/index.html","html/template.html")
	err = t.ExecuteTemplate(w,"template", struct {
		Posts []activity.Activity
	}{
		Posts: act,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler)addActivityPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t,err := template.ParseFiles("html/formact.html","html/template.html")
	err = t.ExecuteTemplate(w,"template", struct {}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler)addActivityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var name ,location ,description,speaker string
	var max int
	var datetime []string
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	for key, values := range r.PostForm {
		switch key {
		case "name":
			name = fmt.Sprintf("%s",values)
		case "location":
			location = fmt.Sprintf("%s",values)
		case "description":
			description = fmt.Sprintf("%s",values)
		case "speaker":
			speaker = fmt.Sprintf("%s",values)
		case "max":
			max,_ = strconv.Atoi(fmt.Sprintf("%s",values))
		case "datetimes[]":
			datetime = append(datetime,fmt.Sprintf("%s",values))
		}
	}
	h.actManage.Insert(&activity.Activity{
		Name        : name,
		Location    : location,
		Speaker     : speaker,
		Description : description,
		Maxjoin     : max,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}




func StartServer(addr string, db *sql.DB) error {

	r := mux.NewRouter()
	h := &Handler{
		actManage: &activity.Manager{
			DB: db,
		},
		dateManage: &date.Manager{
			DB: db,
		},
		pinManage: &pinactivity.Manager{
			DB: db,
		},
	}
	r.HandleFunc("/add_activity", h.addActivityPageHandler).Methods("GET")
	r.HandleFunc("/activity", h.addActivityHandler).Methods("POST")
	r.HandleFunc("/", h.indexPageHandler).Methods("GET")

	return http.ListenAndServe(addr, r)

}
