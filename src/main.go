package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"time"
)

type User struct {
	Name string `json:"name,omitempty"`
}

type List struct {
	List []string `json:"list,omitempty"`
}

var (

	// Primary template for / serving
	homeTemplate = template.Must(template.New("index.html").Delims("[[", "]]").ParseFiles("templates/index.html"))

	// Global variable identifying the version of the app
	// Set to the Unix time at start-up
	VERSION int64

	// DEBUG flag
	DEBUG bool

	// PORT
	PORT string
)

// Main web page server (/ and all AngularJS routes)
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Trace("HomeHandler", r)
	PrintMemUsage()
	DebugInfo(r)

	// Execute and serve HTML template
	if err := homeTemplate.Execute(w, template.FuncMap{
		"Version": VERSION,
		"Debug": DEBUG,
		"User": User{
			Name: "Username",
		},
	}); err != nil {
		InternalServerError(w, "Error with homeTemplate: %v", err)
		return
	}
}


// API List handler (/api/list)
func APIListHandler(w http.ResponseWriter, r *http.Request) {
	Trace("APIListHandler", r)
	
	var list = List{
		List: []string{
			"AAA",
			"BBB",
		},
	}
	if err := WriteJSON(w, &list); err != nil {
		InternalServerError(w, "Error with APIListHandler: %v", err)
		return
	}

}

// Handler to debug an http.Request to the web user
func DumpHandler(w http.ResponseWriter, r *http.Request) {
	Trace("DumpHandler", r)
	PrintMemUsage()
	DebugInfo(r)

	fmt.Fprintf(w, "URL:%v \n", r.URL)
	fmt.Fprintf(w, "Method:%v \n", r.Method)
	fmt.Fprintf(w, "Proto:%v \n", r.Proto)
	fmt.Fprintf(w, "Header:%v \n", r.Header)
	fmt.Fprintf(w, "ContentLength:%v \n", r.ContentLength)
	fmt.Fprintf(w, "Host:%v \n", r.Host)
	fmt.Fprintf(w, "Referer:%v \n", r.Referer())
	fmt.Fprintf(w, "Form:%v \n", r.Form)
	fmt.Fprintf(w, "PostForm:%v \n", r.PostForm)
	fmt.Fprintf(w, "MultipartForm:%v \n", r.MultipartForm)
	fmt.Fprintf(w, "RemoteAddr:%v \n", r.RemoteAddr)
	fmt.Fprintf(w, "RequestURI:%v \n", r.RequestURI)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header %v = %v \n", k, v)
	}

	for _, v := range r.Cookies() {
		fmt.Fprintf(w, "Cookie %v = %v \n", v.Name, v.Value)
	}

	for _, v := range os.Environ() {
		fmt.Fprintf(w, "Env  %v \n", v)
	}
	request, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Fprintf(w, "Error while dumping request: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Request: %v\n", string(request))

}

// Ping Handler for /ping
func PingHandler(w http.ResponseWriter, r *http.Request) {
	Trace("PingHandler", r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "OK")
}

// Memory Handler for /memory
func MemoryHandler(w http.ResponseWriter, r *http.Request) {
	Trace("MemoryHandler", r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"alloc":%v,"totalalloc":%v,"sys":%v}`,
		m.Alloc, m.TotalAlloc, m.Sys)
}

// Worker Task - every 1 minute
func WorkerTask() {
	Debug(">>> WorkerTask")
	PrintMemUsage()
}

// Start cron worker (worker.go)
func StartWorker() {
	Debug(">>> StartWorker")
	gocron.Every(1).Minute().DoSafely(WorkerTask)
	<-gocron.Start()
}

// Main function to register handlers
func main() {

	// Define debug level across the application
	DEBUG = true
	if os.Getenv("DEBUG") != "" {
		DEBUG = os.Getenv("DEBUG") == "1"
	}

	// Define PORT
	PORT = "80"
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	// Set VERSION to current unix time
	VERSION = time.Now().Unix()

	// Start worker cron
	go StartWorker()

	// Define HTTP Router
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/dump", DumpHandler)
	r.HandleFunc("/memory", MemoryHandler)
	r.HandleFunc("/api/list", APIListHandler)	
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/").HandlerFunc(HomeHandler) // Catch-all
	http.Handle("/", r)

	// Start application - port 80 within the Docker image,
	Info("Starting up on %v", PORT)
	defer Info("Stopping application")
	PrintMemUsage()
	log.Fatal(http.ListenAndServe(":"+PORT, nil))

}
