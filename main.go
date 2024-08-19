package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"ascii-web/funcs"
)

var (
	tpl        *template.Template
	port       = "8080"
	text       []string
	banner     []string
	err        error
	ErrType    string
	outputEdit string
	inputs     = []string{"test", "Hello\r\n\r\nThere!", "\r\n\r\n\r\n\r\nI'm after 4 new lines", "I'm before 2 new lines\r\n\r\n", "ABCDEFghij123@", "123,;:!456", "Special Characters\r\n{:>|}{}", ">>NonPrintable:é<<"}
	values     struct {
		Btn        string
		Output     string
		BannerTest []string
	}
)

func main() {
	fmt.Println("http://localhost:" + port)
	http.HandleFunc("/", IndexHanlde)
	http.HandleFunc("/static", StaticHandle)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ascii-art", OutputHandle)
	http.HandleFunc("/download", DownloadHandle)
	http.ListenAndServe(":"+port, nil)
}

// a func for return the statut error code
func ErrorHandle(myStrStatus string) int {
	if myStrStatus == "internal" {
		return http.StatusInternalServerError
	} else if myStrStatus == "method" {
		return http.StatusMethodNotAllowed
	} else if myStrStatus == "request" {
		return http.StatusBadRequest
	} else if myStrStatus == "notfound" {
		return http.StatusNotFound
	}
	return 0
}

func ExecError(ErrType string, w http.ResponseWriter) {
	tpl, _ = template.ParseGlob("static/templates/error.html")
	w.WriteHeader(ErrorHandle(ErrType))
	tpl.ExecuteTemplate(w, "error.html", http.StatusText(ErrorHandle(ErrType)))
}

// a func return the file html who's not exist in the templates directory
func InternalErr() string {
	MySlice := []string{"error", "output", "index"}
	for _, v := range MySlice {
		tpl, err = template.ParseGlob("static/templates/" + v + ".html")
		if err != nil {
			return v
		}
	}
	return ""
}

func ExecutingTemplate(w http.ResponseWriter, tplt string, data any) error {
	tpl, _ = template.ParseGlob("static/templates/" + tplt)
	err := tpl.ExecuteTemplate(w, tplt, data)
	if err != nil {
		return err
	}
	return nil
}

func IndexHanlde(w http.ResponseWriter, r *http.Request) {
	// ---------------------check if html files is exist-------------------------
	if InternalErr() == "error" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if InternalErr() == "index" || InternalErr() == "output" {
		ExecError("internal", w)
		return
	}
	// --------------------------------------------------------------------------

	// ----------------------check the banners if exist----------------------------
	if !funcs.BannerExist() {
		ExecError("internal", w)
		return
	}

	// -----------------------------handling path---------------------------------
	if r.URL.Path != "/" {
		ExecError("notfound", w)
		return
	}

	// ------------------------check the method request-----------------------------
	if r.Method != http.MethodGet {
		ExecError("method", w)
		return
	}

	// ---------------if all good the index template execute---------------------
	ExecutingTemplate(w, "index.html", nil)
}

func OutputHandle(w http.ResponseWriter, r *http.Request) {
	// ---------------------check html files is exist----------------------------
	if InternalErr() == "error" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if InternalErr() == "index" || InternalErr() == "output" {
		ExecError("internal", w)
		return
	}
	// --------------------------------------------------------------------------

	// ----------------------check the banners if exist----------------------------
	if !funcs.BannerExist() {
		ExecError("internal", w)
		return
	}

	// ------------------------check the method request-----------------------------
	if r.Method != http.MethodPost {
		ExecError("method", w)
		return
	}

	data := values
	data.BannerTest = []string{}                   // for store the output of liste of tests
	data.Output = ""                               // for store the output of user input
	banner = append(banner, r.FormValue("banner")) // store the banner
	text = append(text, r.FormValue("textName"))   // store the input text
	data.Btn = r.FormValue("ButtonName")           // store which button the user click
	lastBanner := banner[len(banner)-1]
	lastText := text[len(text)-1]
	if len(lastText) > 1000 {
		ExecError("request", w)
		return
	}
	if data.Btn == "submit" {
		data.Output = funcs.Print(lastText, lastBanner)
		// ------------------check if the user enter the text and banner--------------------
		if len(lastText) != 0 && data.Output != "err" {
			ExecutingTemplate(w, "output.html", data)
		} else {
			ExecError("request", w)
		}

	} else if data.Btn == "test" {
		for _, value := range inputs {
			data.BannerTest = append(data.BannerTest, funcs.Print(value, lastBanner))
		}
		ExecutingTemplate(w, "output.html", data)
	} else {
		ExecError("request", w)
		return
	}
}

func DownloadHandle(w http.ResponseWriter, r *http.Request) {
	//------------------Handle case if the path is "localhost:port/download"------------------
	if r.Method != http.MethodPost {
		ExecError("method", w)
		return
	}

	outputEdit = ""
	for _, v := range r.FormValue("output") {
		if v != 13 {
			outputEdit += string(v)
		}
	}

	//------------------this lines below to range of the outputs to download the rigth one------------------
	for i := len(text) - 1; i >= 0; i-- {
		if funcs.Print(text[i], banner[i]) == outputEdit {
			file_name := "ascii-art.txt"
			w.Header().Set("Content-Disposition", "attachment; filename="+file_name)
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", strconv.Itoa(len(outputEdit)))
			w.Write([]byte(outputEdit))
			break
		} else {
			//------------------if range of all outputs and can't find the similar output in the <form> that mean an error of request happen------------------
			if i == 0 {
				ExecError("request", w)
			}
		}
	}
}

func StaticHandle(w http.ResponseWriter, r *http.Request) {
	ExecError("notfound", w)
}
