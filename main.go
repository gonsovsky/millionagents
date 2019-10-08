package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)
var (
	K = *flag.Int("k", 5, "Max degree of parallelism")
)

func main(){
	var input = make(chan string)

	go func(){
		var semaphore = make(chan bool, K)
		for item := range input {
			semaphore <- true
			go func(item string) {
				defer func() { <-semaphore }()
				work(item)
			}(item)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input <- scanner.Text()
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
	re := regexp.MustCompile( "\\b(Go)\\b")
	var cnt = re.FindAllIndex(bytes, len(bytes))
	log.Printf("%s: %v", item, len(cnt))
}