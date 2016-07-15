package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/cagnosolutions/repono"
	"github.com/cagnosolutions/web"
)

var mux *web.Mux
var tmpl *web.TmplCache
var db *repono.Client

func init() {
	web.STATIC_PATH = "assets"
	web.Funcs["prettyDate"] = prettyDate
	tmpl = web.NewTmplCache()
	mux = web.NewMux()
	db = repono.Embed(repono.EmbedAndServeDB(":7777"))
	db.AddStore("news")
	db.AddStore("message")
	db.AddStore("employee")

}

func main() {
	mux.AddRoutes(index, login, loginpage, send, perform, article, articleView, uploadServer, whatWeDo, terms, privacy, faq, logout)
	mux.AddSecureRoutes(ADMIN, inbox, news, newsView, addNews, addEmployee, calendar, msgView, delNews, delMsg, delEmployee, settings, employee, employeeView, question)
	fmt.Println("------------------------------------------REMEMBER TO REGISTER ALL NEW ROUTES")
	log.Fatal(http.ListenAndServe(":8888", mux))
}

// -----------------USER PAGES-------------------

var index = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	var employee []Employee
	db.GetAll("employee", &employee)
	var news NewsSorted
	db.GetAll("news", &news)
	sort.Stable(news)
	tmpl.Render(w, r, "index.tmpl", web.Model{
		"allNews":     news,
		"allEmployee": employee,
		"footerNews":  getFooterNews(),
	})
}}

var perform = web.Route{"GET", "/perform", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "perform.tmpl", web.Model{
		"footerNews": getFooterNews(),
	})
}}

var whatWeDo = web.Route{"GET", "/whatwedo", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "whatwedo.tmpl", web.Model{
		"footerNews": getFooterNews(),
	})
}}

var article = web.Route{"GET", "/article", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "article.tmpl", web.Model{
		"footerNews": getFooterNews(),
	})
}}

var loginpage = web.Route{"GET", "/login", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "login.tmpl", web.Model{
		"footerNews": getFooterNews(),
	})
}}

var terms = web.Route{"GET", "/terms", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "terms.tmpl", web.Model{
		"footerNews": getFooterNews(),
	})
}}

var privacy = web.Route{"GET", "/privacy", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "privacy.tmpl", web.Model{
		"footerNews": getFooterNews(),
	})
}}

var faq = web.Route{"GET", "/faq", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "faq.tmpl", web.Model{
		"footerNews": getFooterNews(),
	})
}}

// -----------------END USER PAGES-------------------

// -----------------USER CONTROLLERS-------------------

var login = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("username")
	password := r.FormValue("password")
	if email == "admin" && password == "admin" {
		web.Login(w, r, "ADMIN")
		web.SetSuccessRedirect(w, r, "/inbox", "Welcome Admin!")
		return
	}
	web.SetErrorRedirect(w, r, "/login", "Incorrect Email or Password")
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	web.SetSuccessRedirect(w, r, "/", "Logged Out")
	return
}}

var send = web.Route{"POST", "/send", func(w http.ResponseWriter, r *http.Request) {

	id := strconv.Itoa(int(time.Now().UnixNano()))
	question, _ := strconv.ParseBool(r.FormValue("question"))
	m := Message{
		Id:       id,
		FullName: r.FormValue("fullname"),
		Email:    r.FormValue("email"),
		Phone:    r.FormValue("phone"),
		Subject:  r.FormValue("subject"),
		Message:  r.FormValue("message"),
		Question: question,
	}
	db.Add("message", id, m)
	web.SetSuccessRedirect(w, r, r.FormValue("redirect"), "Message Sent")
}}

// -----------------END USER CONTROLLERS-------------------

// -----------------ADMIN PAGES-------------------

var inbox = web.Route{"GET", "/inbox", func(w http.ResponseWriter, r *http.Request) {
	var allMsgs []Message
	db.GetAll("message", &allMsgs)

	var filteredMsgs []Message

	for _, m := range allMsgs {
		if !m.Question {
			filteredMsgs = append(filteredMsgs, m)
		}
	}

	tmpl.Render(w, r, "inbox.tmpl", web.Model{
		"allMsgs": filteredMsgs,
		"page":    "inbox",
	})
}}

var question = web.Route{"GET", "/inbox/question", func(w http.ResponseWriter, r *http.Request) {
	var allMsgs []Message
	db.GetAll("message", &allMsgs)

	var filteredMsgs []Message

	for _, m := range allMsgs {
		if m.Question {
			filteredMsgs = append(filteredMsgs, m)
		}
	}

	tmpl.Render(w, r, "question.tmpl", web.Model{
		"allMsgs": filteredMsgs,
		"page":    "inbox",
	})
}}

var msgView = web.Route{"GET", "/message/:id", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(":id")
	var msg Message
	db.Get("message", id, &msg)
	tmpl.Render(w, r, "msg-view.tmpl", web.Model{
		"msg":  msg,
		"page": "inbox",
	})
}}

var news = web.Route{"GET", "/news", func(w http.ResponseWriter, r *http.Request) {
	var allNews []News
	db.GetAll("news", &allNews)
	tmpl.Render(w, r, "news.tmpl", web.Model{
		"allNews": allNews,
		"page":    "news",
	})
}}

var newsView = web.Route{"GET", "/news/:id", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(":id")
	var news News
	if id != "add" {
		db.Get("news", id, &news)
	}
	tmpl.Render(w, r, "news-view.tmpl", web.Model{
		"news": news,
		"page": "news",
	})
}}

var articleView = web.Route{"GET", "/article/:id", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(":id")
	var news News
	var allNews NewsSorted
	db.GetAll("news", &allNews)
	db.Get("news", id, &news)
	sort.Stable(allNews)
	tmpl.Render(w, r, "article.tmpl", web.Model{
		"news":       news,
		"allNews":    allNews,
		"footerNews": getFooterNews(),
	})
}}

var employee = web.Route{"GET", "/employee", func(w http.ResponseWriter, r *http.Request) {
	var allEmployees []Employee
	db.GetAll("employee", &allEmployees)
	tmpl.Render(w, r, "employee.tmpl", web.Model{
		"allEmployees": allEmployees,
		"page":         "employee",
	})
}}

var employeeView = web.Route{"GET", "/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(":id")
	var employee Employee
	if id != "add" {
		db.Get("employee", id, &employee)
	}
	tmpl.Render(w, r, "employee-view.tmpl", web.Model{
		"employee": employee,
		"page":     "employee",
	})
}}

var calendar = web.Route{"GET", "/calendar", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "calendar.tmpl", web.Model{
		"page": "calendar",
	})
}}

var settings = web.Route{"GET", "/settings", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "settings.tmpl", web.Model{
		"page": "settings",
	})
}}

// -----------------END ADMIN PAGES-------------------

// -----------------ADMIN CONTROLLERS-------------------

var addNews = web.Route{"POST", "/news", func(w http.ResponseWriter, r *http.Request) {
	date := time.Now().UnixNano()
	id := strconv.Itoa(int(date))
	var news News
	if r.FormValue("id") != "" {
		id = r.FormValue("id")
		db.Get("news", id, &news)
	}
	path := "upload/news/"
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("main.go >> addNews >> os.MkdirAll() >> %v\n", err)
		fmt.Fprintf(w, "error saving")
		return
	}

	news.Id = id
	news.Title = r.FormValue("title")
	news.FirstName = r.FormValue("firstname")
	news.LastName = r.FormValue("lastname")
	news.Description = r.FormValue("description")
	news.Body = r.FormValue("body")
	news.Date = date

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err == nil {
		if len(handler.Header["Content-Type"]) < 1 {
			log.Printf("main.go >> addNews >> %v\n", err)
			fmt.Fprintf(w, "error saving")
			return
		}
		defer file.Close()
		if handler.Header["Content-Type"][0] != "image/png" && handler.Header["Content-Type"][0] != "image/jpeg" {
			log.Printf("main.go >> addNews >> did not recieve image\n")
			fmt.Fprintf(w, "error saving")
			return
		}
		f, err := os.OpenFile(path+id+".png", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("main.go >> addNews >>  os.OpenFile() >> %v\n", err)
			fmt.Fprintf(w, "error saving")
			return
		}
		defer f.Close()
		io.Copy(f, file)
		news.Image = true
	}

	db.Set("news", id, news)
	web.SetSuccessRedirect(w, r, "/news", "Article Added")
}}

var delNews = web.Route{"POST", "/delnews", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var news News
	db.Get("news", id, &news)
	if news.Image {
		path := "upload/news/"
		if err := os.Remove(path + id + ".png"); err != nil {
			log.Printf("main.go ->> delNews ->> os.Remove() ->> %v\n", err)
			web.SetErrorRedirect(w, r, "/news/"+id, "Error Deleting Article")
			return
		}
	}

	db.Del("news", id)
	web.SetSuccessRedirect(w, r, "/news", "Article Deleted")
	return
}}

var addEmployee = web.Route{"POST", "/employee", func(w http.ResponseWriter, r *http.Request) {
	date := time.Now().UnixNano()
	id := strconv.Itoa(int(date))
	var employee Employee
	if r.FormValue("id") != "" {
		id = r.FormValue("id")
		db.Get("employee", id, &employee)
	}
	path := "upload/employee/"
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("main.go >> addEmployee >> os.MkdirAll() >> %v\n", err)
		fmt.Fprintf(w, "error saving")
		return
	}

	employee.Id = id
	employee.FirstName = r.FormValue("firstname")
	employee.LastName = r.FormValue("lastname")
	employee.Title = r.FormValue("title")
	employee.Description = r.FormValue("description")
	employee.Location = r.FormValue("location")

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err == nil {
		if len(handler.Header["Content-Type"]) < 1 {
			log.Printf("main.go >> addEmployee >> %v\n", err)
			fmt.Fprintf(w, "error saving")
			return
		}
		defer file.Close()
		if handler.Header["Content-Type"][0] != "image/png" && handler.Header["Content-Type"][0] != "image/jpeg" {
			log.Printf("main.go >> addEmployee >> did not recieve image\n")
			fmt.Fprintf(w, "error saving")
			return
		}
		f, err := os.OpenFile(path+id+".png", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("main.go >> addEmployee >>  os.OpenFile() >> %v\n", err)
			fmt.Fprintf(w, "error saving")
			return
		}
		defer f.Close()
		io.Copy(f, file)
		employee.Image = true
	}

	db.Set("employee", id, employee)
	web.SetSuccessRedirect(w, r, "/employee", "Employee Added")
}}

var delEmployee = web.Route{"POST", "/delemployee", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var employee Employee
	db.Get("employee", id, &employee)
	if employee.Image {
		path := "upload/employee/"
		if err := os.Remove(path + id + ".png"); err != nil {
			log.Printf("main.go ->> delEmployee ->> os.Remove() ->> %v\n", err)
			web.SetErrorRedirect(w, r, "/employee/"+id, "Error Deleting Employee")
			return
		}
	}

	db.Del("employee", id)
	web.SetSuccessRedirect(w, r, "/employee", "Employee Deleted")
	return
}}

var delMsg = web.Route{"POST", "/delmsg", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	db.Del("message", id)
	web.SetSuccessRedirect(w, r, "/inbox", "Message Deleted")
	return
}}

var uploadServer = web.Route{"GET", "/upload/:folder/:file", func(w http.ResponseWriter, r *http.Request) {
	dir := "upload/" + r.FormValue(":folder")
	server := http.StripPrefix("/"+dir, http.FileServer(http.Dir(dir)))
	server.ServeHTTP(w, r)
}}

func prettyDate(s string) string {
	ts, err := strconv.Atoi(s)
	if err != nil {
		return ""
	}
	t := time.Unix(0, int64(ts))
	return t.Format("January 2, 2006")
}

func getFooterNews() NewsSorted {
	var news NewsSorted
	db.GetAll("news", &news)
	sort.Stable(news)
	if len(news) > 3 {
		news = news[:3]
	}
	return news
}

// -----------------END ADMIN CONTROLLERS-------------------t
