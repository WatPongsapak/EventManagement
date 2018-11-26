package handler

import (
	"strconv"
	"database/sql"
	"event-management/activity"
	"event-management/pinactivity"
	"event-management/admin"
	"html/template"
	"net/http"
	"fmt"
	"strings"
	"time"
	
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/securecookie"
)

type Handler struct {
	actManage *activity.Manager
	pinManage *pinactivity.Manager
	adminManage *admin.Manager
}

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

func (h *Handler)indexPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	session, _ := store.Get(r, "cookie-name")
	act, err := h.actManage.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	param := struct {
		Posts []activity.Activity
		Admin bool
	}{
		Posts: act,
		Admin: false,
	}
	if username, ok := session.Values["username"].(string); ok && username != "" {
		param.Admin = true
	}
	
	t,err := template.ParseFiles("html/index.html","html/template.html")
	err = t.ExecuteTemplate(w,"template", param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler)addActivityPage(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler)pinActivityPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t,err := template.ParseFiles("html/formpin.html","html/template.html")
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])
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



func (h *Handler)editActivityPage(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler)addActivity(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler)editActivity(w http.ResponseWriter, r *http.Request) {
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


func (h *Handler)delActivity(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler)pinActivity(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler)pinresultActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])
	t,err := template.ParseFiles("html/joinlist.html","html/template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p,err := h.pinManage.All(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		P []pinactivity.Pinactivity
	}{
		A : *a,
		Date: a.StartDate.Format("01/02/2006")+" - "+a.EndDate.Format("01/02/2006"),
		Time: a.StartTime.Format("03:04 pm")+" - "+a.EndTime.Format("03:04 pm"),
		ID:id,
		P : p,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler)generateExcel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id,_ := strconv.Atoi(vars["id"])
	p,err := h.pinManage.All(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a,err := h.actManage.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    xlsx := excelize.NewFile()
	xlsx.NewSheet("Sheet1")
	xlsx.SetCellValue("Sheet1", "A1", "ชื่อกิจกรรม")
	xlsx.SetCellValue("Sheet1", "B1", a.Name)
	xlsx.SetCellValue("Sheet1", "A2", "รายระเอียด")
	xlsx.SetCellValue("Sheet1", "B2", a.Description)
	xlsx.SetCellValue("Sheet1", "A3", "ผู้บรรยาย")
	xlsx.SetCellValue("Sheet1", "B3", a.Speaker)
	xlsx.SetCellValue("Sheet1", "A4", "จำนวนสูงสุด")
	xlsx.SetCellValue("Sheet1", "B4", a.Maxjoin)
	xlsx.SetCellValue("Sheet1", "A5", "สถานที่")
	xlsx.SetCellValue("Sheet1", "B5", a.Location)
	xlsx.SetCellValue("Sheet1", "A6", "วัน")
	xlsx.SetCellValue("Sheet1", "B6", a.StartDate.Format("01/02/2006")+" - "+a.EndDate.Format("01/02/2006"))
	xlsx.SetCellValue("Sheet1", "A7", "เวลา")
	xlsx.SetCellValue("Sheet1", "B7", a.StartTime.Format("03:04 pm")+" - "+a.EndTime.Format("03:04 pm"))

	xlsx.SetCellValue("Sheet1", "A9", "รหัสพนักงาน")
	xlsx.SetCellValue("Sheet1", "B9", "ชื่อนามสกิล")
	xlsx.SetCellValue("Sheet1", "C9", "เบอร์โทรศัพท์")
	
	rowIndex := 10
	for _, v := range p {
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(rowIndex), v.EmployeeCode)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(rowIndex), v.Name)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(rowIndex), v.Phone)
		rowIndex++
	}
    xlsx.SetActiveSheet(1)
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", "attachment; filename="+a.Name+" รุ่น"+strconv.Itoa(a.Round)+".xlsx")
    w.Header().Set("Content-Transfer-Encoding", "binary")
    w.Header().Set("Expires", "0")
    xlsx.Write(w)
}

func (h *Handler)loginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t,err := template.ParseFiles("html/formlogin.html","html/template.html")
	err = t.ExecuteTemplate(w,"template", struct {}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler)login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	username :=  r.FormValue("username")
	password :=  r.FormValue("password")
	fmt.Println(username,password)
	admin, err := h.adminManage.Login(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}else{
	session, _ := store.Get(r, "session-name")
	session.Values["username"] = admin.Username
	session.Save(r, w)
	}
	
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
		adminManage: &admin.Manager{
			DB: db,
		},
	}
	r.HandleFunc("/add_activity", h.addActivityPage).Methods("GET")
	r.HandleFunc("/edit_activity/{id}", h.editActivityPage).Methods("GET")
	r.HandleFunc("/pin_activity/{id}", h.pinActivityPage).Methods("GET")
	r.HandleFunc("/pinresult/{id}", h.pinresultActivity).Methods("GET")
	r.HandleFunc("/export_excel/{id}", h.generateExcel).Methods("GET")
	r.HandleFunc("/loginpage", h.loginPage).Methods("GET")


	r.HandleFunc("/activity/del/{id}", h.delActivity).Methods("GET")
	r.HandleFunc("/activity/{id}", h.editActivity).Methods("POST")
	r.HandleFunc("/activity", h.addActivity).Methods("POST")
	r.HandleFunc("/pinact/{id}", h.pinActivity).Methods("POST")
	r.HandleFunc("/login", h.login).Methods("POST")
	
	
	r.HandleFunc("/", h.indexPage).Methods("GET")

	return http.ListenAndServe(addr, r)

}
