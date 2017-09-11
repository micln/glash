package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"time"

	"glash/crawler"
	"strings"

	. "github.com/micln/go-utils"
)

type Job struct {
	Url      string
	Folder   string
	FileName string
}

type EventFunc func(*Job)

type Queue struct {
	IsRunning bool

	LimitConsumer      int
	LimitQueueCapacity int

	jobs chan *Job

	eventBeforeDownload []EventFunc
	eventAfterDownload  []EventFunc

	Crawler string
}

func (q *Queue) Log(format string, args ...interface{}) {
	log.Printf("[Queue] "+format+"\n", args...)
}

func (this *Queue) Listen() {
	this.Log(`Starting...`)

	this.jobs = make(chan *Job, this.LimitQueueCapacity)
	this.Log(`LimitQueueCapacity: %d`, this.LimitQueueCapacity)

	wg := sync.WaitGroup{}
	wg.Add(this.LimitConsumer)

	for i := 0; i < this.LimitConsumer; i++ {
		go func(idx int) {
			wg.Done()

			var job *Job
			for {
				job = <-this.jobs
				var err error

				this.Log(`Consumer[%d]: %s`, idx, JsonEncode(job))

				if len(job.Folder) > 0 {
					_, err = os.Open(job.Folder)
					if err != nil {
						os.MkdirAll(job.Folder, 0755)
					}
				}

				//cmd := exec.Command(`axel`, job.Url, `-o`, path.Join(job.Folder, job.FileName))
				//crawler := crawler.Which(`aria2`)
				crawler := crawler.Which(this.Crawler)
				crawler.SetUrl(job.Url).
					SetPath(path.Join(job.Folder, job.FileName))
				cmd := exec.Command(crawler.Cmd(), crawler.Args()...)
				//log.Fatalln(cmd)
				output, err := cmd.Output()
				fmt.Printf(`%s`, output)
				if err != nil {
					log.Println(`Download Error:`, err)
				}

				for _, evt := range this.eventAfterDownload {
					this.FireEvent(evt, job)
				}
			}
		}(i)
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			this.Log(`Status : [%d/%d] Jobs`, len(this.jobs), cap(this.jobs))
		}
	}()

	wg.Wait()

	this.IsRunning = true

	this.Log(`Listening...`)
}

func (q *Queue) WaitingStarted() {
	for {
		if q.IsRunning {
			return
		}
		runtime.Gosched()
	}
}

func (q *Queue) PushJob(job *Job) {
	q.WaitingStarted()

	go func(job *Job) {
		job.Folder, job.FileName = parsePath(job.FileName)

		job.Folder = path.Join(`Downloads`, Date(`Ymd`), job.Folder)

		q.jobs <- job
	}(job)
}

func (q *Queue) Push(url string, filename string) *Job {
	job := &Job{
		Url:      url,
		FileName: filename,
	}

	defaultQueue.PushJob(job)

	return job
}

func (q *Queue) AfterDownload(f EventFunc) {
	q.eventAfterDownload = append(q.eventAfterDownload, f)
}

func (q *Queue) FireEvent(f EventFunc, job *Job) {
	f(job)
}

//  把一个路径拆分为 文件夹 和 文件名
func parsePath(filePath string) (folder string, filename string) {
	paths := strings.Split(strings.TrimRight(filePath, `/`), `/`)
	num := len(paths)
	for idx, seq := range paths {
		if idx+1 == num {
			filename = seq
		} else {
			folder = path.Join(folder, seq)
		}
	}
	return
}
