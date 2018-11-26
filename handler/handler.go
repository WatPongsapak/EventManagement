package handler

import (
	"strconv"
	"database/sql"
	"event-management/activity"
	"event-management/pinactivity"
	"html/template"
	"net/http"
	"fmt"
	"strings"
	"time"
	
	"github.com/gorilla/mux"
)

type Handler struct {
	actManage *activity.Manager
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
	err = t.ExecuteTemplate(w,"template", struct {
		Mode string
		A activity.Activity
		Date string
		Time string
	}{
		Mode : "add",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler)pinActivityPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t,err := template.ParseFiles("html/formpin.html","html/template.html")
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])
	fmt.Println(id)
	a,err := h.actManage.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w,"template", struct {
		A activity.Activity
		Date string
		Time string
		ID int
	}{
		A : *a,
		Date: a.StartDate.Format("01/02/2006")+" - "+a.EndDate.Format("01/02/2006"),
		Time: a.StartTime.Format("03:04 pm")+" - "+a.EndTime.Format("03:04 pm"),
		ID:id,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}



func (h *Handler)editActivityPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])
	fmt.Println(id)
	a,err := h.actManage.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(a)
	t,err := template.ParseFiles("html/formact.html","html/template.html")
	err = t.ExecuteTemplate(w,"template", struct {
		Mode string
		A activity.Activity
		Date string
		Time string
	}{
		Mode : "edit",
		A : *a,
		Date: a.StartDate.Format("01/02/2006")+" - "+a.EndDate.Format("01/02/2006"),
		Time: a.StartTime.Format("03:04 pm")+" - "+a.EndTime.Format("03:04 pm"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler)addActivityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.actManage.Insert(h.getForm(&activity.Activity{},w,r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler)getForm(a *activity.Activity, w http.ResponseWriter, r *http.Request) *activity.Activity{
	var name,location,description,speaker string
	var max int64
	datearr :=  strings.Split(r.FormValue("daterange"), " - ")
	timearr :=  strings.Split(r.FormValue("timerange"), " - ")
	starttime, _ := time.Parse("03:04 pm", timearr[0])
	endtime, _ := time.Parse("03:04 pm", timearr[1])
	startdate, _ := time.Parse("01/02/2006", datearr[0])
	enddate, _ := time.Parse("01/02/2006", datearr[1])
	name =  r.FormValue("name")
	description =  r.FormValue("description")
	speaker =  r.FormValue("speaker")
	max, _ = strconv.ParseInt(r.FormValue("max")[0:], 10, 64);
	location =  r.FormValue("location")
	a.Name = name
	a.Description = description
	a.Speaker = speaker
	a.Maxjoin = int(max)
	a.Location = location
	a.StartDate = startdate
	a.EndDate = enddate
	a.StartTime = starttime
	a.EndTime =  endtime
	return a
}

func (h *Handler)editActivityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])
	a,err := h.actManage.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.actManage.Update(h.getForm(a,w,r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	
}


func (h *Handler)delActivityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])
	a,err := h.actManage.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.actManage.Delete(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	
}

func (h *Handler)pinActivityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])
	name :=  r.FormValue("name")
	employeeid :=  r.FormValue("employeeid")
	phone :=  r.FormValue("phone")

	h.pinManage.Insert(&pinactivity.Pinactivity{
		ActivitiesID:id,
		EmployeeCode:employeeid,
		Name        :name,
		Phone       :phone,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func StartServer(addr string, db *sql.DB) error {

	r := mux.NewRouter()
	h := &Handler{
		actManage: &activity.Manager{
			DB: db,
		},
		pinManage: &pinactivity.Manager{
			DB: db,
		},
	}
	r.HandleFunc("/add_activity", h.addActivityPageHandler).Methods("GET")
	r.HandleFunc("/edit_activity/{id}", h.editActivityPageHandler).Methods("GET")
	r.HandleFunc("/pin_activity/{id}", h.pinActivityPageHandler).Methods("GET")

	r.HandleFunc("/activity/del/{id}", h.delActivityHandler).Methods("GET")
	r.HandleFunc("/activity/{id}", h.editActivityHandler).Methods("POST")
	r.HandleFunc("/activity", h.addActivityHandler).Methods("POST")
	r.HandleFunc("/pinact/{id}", h.pinActivityHandler).Methods("POST")
	
	r.HandleFunc("/", h.indexPageHandler).Methods("GET")

	return http.ListenAndServe(addr, r)

}
