package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var reqnum int
var port string
var help = flag.Bool("help", false, "Show help")
var wg = &sync.WaitGroup{}

func init() {
	flag.IntVar(&reqnum, "reqnum", 1, "Specify a number of requests")
	flag.StringVar(&port, "port", "9090", "Specify port")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	rand.Seed(time.Now().Unix())
}

func main() {
	defer timeTrack(time.Now())

	for i := 0; i < reqnum; i++ {
		wg.Add(1)
		go sendRequestWithValue(rand.Intn(100), wg)

		duration := time.Duration(rand.Intn(20) + 10) // Define random interval from 10 to 30 ms (provides ~ 50 RPS)
		time.Sleep(time.Millisecond * duration)
	}

	wg.Wait()
}

func sendRequestWithValue(value int, group *sync.WaitGroup) {
	defer group.Done()

	jsonBody := []byte(fmt.Sprintf("{\"value\":%v}", value))
	bodyReader := bytes.NewReader(jsonBody)
	res, err := http.Post(fmt.Sprintf("http://127.0.0.1:%v", port), "application/json", bodyReader)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.StatusCode)
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("Execution took %s\n", elapsed)
}
