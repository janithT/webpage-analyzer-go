package channels

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/janithT/webpage-analyzer/config"
)

// WorkFunc is the callback for each executed URL
type WorkFunc func(url string, status int, latency int64)

type UrlWorker interface {
	Create() UrlWorker
	Build(url string, wg *sync.WaitGroup, fn WorkFunc) UrlWorker
	PushChannel()
}

type urlWorker struct {
	url string
	wg  *sync.WaitGroup
	fn  WorkFunc
}

var urlWorkChannel chan urlWorker
var ChanSvrStart = make(chan bool, 1)
var atomicVar int32

// NewUrlWorker returns a new UrlWorker instance
func NewUrlWorker() UrlWorker {
	return &urlWorker{}
}

func (u *urlWorker) Create() UrlWorker {
	return &urlWorker{}
}

func (u *urlWorker) Build(url string, wg *sync.WaitGroup, fn WorkFunc) UrlWorker {
	u.url = url
	u.wg = wg
	u.fn = fn
	return u
}

func (u *urlWorker) PushChannel() {
	urlWorkChannel <- *u
}

// InitializetPageUrlWorkerThreadPool starts 10 workers
func InitializetPageUrlWorkerThreadPool(threadCount int) {
	urlWorkChannel = make(chan urlWorker)

	for i := 1; i <= threadCount; i++ {
		go executeUrl(urlWorkChannel, i, threadCount)
	}
}

// executeUrl runs worker loop
func executeUrl(channel chan urlWorker, id int, maxId int) {
	log.Printf("Thread started: %d", id)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	atomic.AddInt32(&atomicVar, 1)
	if int(atomicVar) == maxId {
		ChanSvrStart <- true
	}

SignalBreakLabel:
	for {
		select {
		case call := <-channel:
			execute(call.url, call.wg, call.fn)
		case <-signals:
			break SignalBreakLabel
		}
	}
}

// execute performs the HTTP GET
func execute(url string, wg *sync.WaitGroup, fn WorkFunc) {
	startTime := time.Now()
	defer wg.Done()

	client := &http.Client{
		Timeout: config.GetAppConfig().GetLinkTimeout(),
	}

	res, err := client.Get(url)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fn(url, http.StatusGatewayTimeout, time.Since(startTime).Milliseconds())
		} else {
			fn(url, http.StatusBadRequest, time.Since(startTime).Milliseconds())
		}
		log.Println("Error in getting response", url, err)
		return
	}
	defer res.Body.Close()
	fn(url, res.StatusCode, time.Since(startTime).Milliseconds())
}
