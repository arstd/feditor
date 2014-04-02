package main


import (
	"net/http"
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"io"
	"sort"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/nav/sub", navSubHandler)
	http.HandleFunc("/main/view", mainViewHandler)
	http.HandleFunc("/main/save", mainSaveHandler)
	http.ListenAndServe(":3080", nil)
}

type fileInfo struct {
	Name string
	IsDir bool
	Path string
}

type fileInfos []fileInfo
func (s fileInfos) Len() int {
    return len(s)
}
func (s fileInfos) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s fileInfos) Less(i, j int) bool {
    return s[i].Name < s[j].Name
}


func navSubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	r.ParseForm()
	// dir := "/opt/gowork/src/github.com/gavinsh/gate"
	dirname := r.Form["dir"][0]

	dir, _ := os.Open(dirname)
	defer dir.Close()

	fs, _:= dir.Readdir(0)
	fis := []fileInfo{}
	for _, fi := range fs {
		filename := fi.Name()
		if fi.IsDir() {
			filename += "/"
		}
//		fmt.Println(filename)
		fis = append(fis, fileInfo{
			Name: filename,
			IsDir: fi.IsDir(),
			Path: dirname + filename,
		})
	}
	sort.Sort(fileInfos(fis))
	subs, _ := json.Marshal(fis)

	fmt.Fprint(w, string(subs))
}

func mainViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	r.ParseForm()
	// dir := "/opt/gowork/src/github.com/gavinsh/gate"
	filename := r.Form["file"][0]

	content, _ := ioutil.ReadFile(filename)

	fmt.Fprint(w, string(content))
}

func mainSaveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	r.ParseForm()
	// dir := "/opt/gowork/src/github.com/gavinsh/gate"
	filename := r.Form["file"][0]

	defer r.Body.Close()
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_SYNC, 0644)
	if err != nil {
		// panic?
	}
	defer file.Close()
	written, _ := io.Copy(file, r.Body) // err

	fmt.Fprint(w, "Write ", written, " bytes")
}
