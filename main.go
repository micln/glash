package main

import (
	"fmt"
	"net/http"
	"os"
	. "pgo"

	"flag"
	"io"
	"io/ioutil"
	"os/signal"
)

var defaultQueue *Queue
var defaultAddr = `0.0.0.0:10087`

func init() {
	defaultQueue = &Queue{
		LimitQueueCapacity: 10000,
		LimitConsumer:      50,
	}

	go defaultQueue.Listen()
}

func webHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get(`url`)
	fileName := r.URL.Query().Get(`filename`)

	job := defaultQueue.Push(url, fileName)

	w.Write([]byte(JsonEncode(job)))
}

func main() {

	startServer := flag.Bool(`d`, false, "start http server.\n\tlisten on http://"+defaultAddr+`?url=xx&filename=xx`)
	isJson := flag.Bool(`json`, false, `parse the file as json`)
	isTab := flag.Bool(`t`, false, `parse the file like 'url(\t)filename'`)

	filename := flag.String(`f`, ``, `filename. the file will be download to "Downloads/`+Date(`Ymd`)+`/filename"`)

	tool := flag.String(`tool`, ``, `choose the download tools. support [curl, aria2]`)

	flag.Parse()

	defaultQueue.Crawler = *tool

	switch {

	case *isTab:
		DownloadFromTabFile(*filename)

	case *isJson:
		DownloadFromJsonFile(*filename)

	case *startServer:
		http.HandleFunc(`/`, webHandler)
		http.ListenAndServe(defaultAddr, nil)

	default:
		flag.Usage()
		os.Exit(0)
	}
}

func DownloadFromJsonFile(filename string) {
	content, err := ioutil.ReadFile(filename)
	Err.Fatal(err)

	files := JsonDecode(content).([]interface{})

	for idx, val := range files {
		fmt.Println(idx)

		file := val.(map[string]interface{})

		defaultQueue.Push(file["url"].(string), file["filename"].(string))
	}

	c := make(chan bool)
	<-c
}

func DownloadFromTabFile(file string) {
	fi, err := os.OpenFile(file, os.O_RDONLY, 0775)
	Err.Fatal(err)

	var url string
	var filename string

	for {
		_, err := fmt.Fscanf(fi, "%s\t%s", &url, &filename)
		if err == io.EOF {
			break
		}

		defaultQueue.Push(url, filename)
	}

	c := make(chan bool)
	<-c
}

func ListenKill() {
	c := make(chan os.Signal)
	signal.Notify(c)

	s := <-c
	fmt.Println(s)
	fmt.Println(defaultQueue)

	signal.Stop(c)
}
