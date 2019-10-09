package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync/atomic"
)

var (
	K = *flag.Int("k", 5, "Max degree of parallelism")
)

func main() {
	var input = make(chan string)
	go func() {
		var semaphore = make(chan bool, K)
		var total int32 = 0
		for {
			select {
			case item, ok := <-input:
				if ok {
					semaphore <- true
					go func(item string) {
						defer func() {
							<-semaphore
						}()
						work(item)
						atomic.AddInt32(&total, 1)
					}(item)
				} else {
					//nothing
				}
			default:
				if x := atomic.LoadInt32(&total); x != 0 && len(semaphore) == 0 {
					atomic.StoreInt32(&total, 0)
					log.Printf("Processed: %v urls.\n", x)
				}
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			input <- scanner.Text()
		}
	}
}

func work(item string) {
	var client http.Client
	resp, err := client.Get(item)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	re := regexp.MustCompile("\\b(Go)\\b")
	var cnt = re.FindAllIndex(bytes, len(bytes))
	log.Printf("%s: %v", item, len(cnt))
}
