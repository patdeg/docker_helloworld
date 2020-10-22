package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strconv"
)

func S2I(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func Debug(format string, a ...interface{}) {
	if !DEBUG {
		return
	}
	fmt.Printf(format+"\n", a...)
}

func Trace(function string, r *http.Request) {
	Info(">>> %v: %v%v (%v)", function, r.Host, r.RequestURI, r.RemoteAddr)
}

func Info(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func DebugOS() {
	Debug("Environment variables:")
	for _, e := range os.Environ() {
		Debug("%v", e)
	}
	Debug("Process id: %v", os.Getpid())
	Debug("Parent Process id: %v", os.Getppid())
	if host, err := os.Hostname(); err == nil {
		Debug("Hostname: %v", host)
	}
}

func Error(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", a...)
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	Debug("Alloc = %v MiB \t TotalAlloc = %v MiB \t Sys = %v MiB", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys))
}

func DebugInfo(r *http.Request) {
	Debug("URL:%v ", r.URL)
	Debug("Method:%v ", r.Method)
	Debug("Proto:%v ", r.Proto)
	Debug("Header:%v ", r.Header)
	Debug("ContentLength:%v ", r.ContentLength)
	Debug("Host:%v ", r.Host)
	Debug("Referer:%v ", r.Referer())
	Debug("Form:%v ", r.Form)
	Debug("PostForm:%v ", r.PostForm)
	Debug("MultipartForm:%v ", r.MultipartForm)
	Debug("RemoteAddr:%v ", r.RemoteAddr)
	Debug("RequestURI:%v ", r.RequestURI)
	for k, v := range r.Header {
		Debug("Header %v = %v ", k, v)
	}

	for _, v := range r.Cookies() {
		Debug("Cookie %v = %v", v.Name, v.Value)
	}

	request, err := httputil.DumpRequest(r, true)
	if err != nil {
		Debug("Error while dumping request: %v", err)
		return
	}
	Debug("Request: %v", string(request))
}

func DebugRequest(r *http.Request) {
	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		Debug("Error while dumping request: %v", err)
		return
	}
	Debug("Request: %q", dump)
}

func DebugResponse(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		Debug("Error while dumping response: %v", err)
		return
	}
	Debug("Response: %q", dump)
}

func GetBody(r *http.Request) []byte {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		Debug("Error while dumping request: %v", err)
		return []byte{}
	}
	return buffer.Bytes()
}

func GetBodyResponse(r *http.Response) []byte {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		Debug("Error while reading body: %v", err)
		return []byte{}
	}
	return buffer.Bytes()
}

func WriteJSON(w http.ResponseWriter, d interface{}) error {
	jsonData, err := json.Marshal(d)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", jsonData)
	return nil
}

func ReadJSON(b []byte, d interface{}) error {
	return json.Unmarshal(b, d)
}

func WriteXML(w http.ResponseWriter, d interface{}) error {
	xmlData, err := xml.Marshal(d)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, "%s", xmlData)
	return nil
}

func ReadXML(b []byte, d interface{}) error {
	return xml.Unmarshal(b, d)
}

func UnmarshalRequest(r *http.Request, value interface{}) error {

	body := GetBody(r)
	Debug("Response: %s", body)

	err := json.Unmarshal(body, value)
	if err != nil {
		return err
	}

	return nil
}

func UnmarshalResponse(r *http.Response, value interface{}) error {

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	Debug("Response: %s", body)

	err = json.Unmarshal(body, value)
	if err != nil {
		return err
	}

	return nil
}

func InternalServerError(w http.ResponseWriter, format string, a ...interface{}) {
	errorMessage := fmt.Sprintf(format, a...)
	Error(errorMessage)
	http.Error(w, errorMessage, http.StatusInternalServerError)
}

func BadRequestError(w http.ResponseWriter, format string, a ...interface{}) {
	errorMessage := fmt.Sprintf(format, a...)
	Error(errorMessage)
	http.Error(w, errorMessage, http.StatusBadRequest)
}

func UnauthorizedError(w http.ResponseWriter, format string, a ...interface{}) {
	errorMessage := fmt.Sprintf(format, a...)
	Error(errorMessage)
	http.Error(w, errorMessage, http.StatusUnauthorized)
}
