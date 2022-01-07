package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// Waitgroup for syncronizing threading
var wg = sync.WaitGroup{}

// AMG Hello Go code to try different things out

func main() {

	args := os.Args
	var options bool

	if len(os.Args) >= 2 {

		fmt.Printf("This is AMG's Go Hello with command line options: %v\n", args)
		options = true

	} else {

		fmt.Printf("This is AMG's Go Hello without command line options: %v\n", args)
		options = false

	}

	if options {

		// Say hello to all of the names
		for i, v := range args {

			if i == 0 {

				fmt.Printf("Saying Hello to:\n")

			} else {

				fmt.Printf("Hello Argument %v %v\n", i, v)

				// Post a Goodbye message to the argument option
				wg.Add(1)
				go postGoodbyeMessage((v))
			}

		}
	}

	// Lets also try a simple HTTP server

	http.HandleFunc("/hello", helloHttpServer)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8091", nil)

	// Make sure to wait for all the threads to stop before actually exiting
	wg.Wait()
}

// Example of a function call with a delay

func postGoodbyeMessage(name string) {

	time.Sleep(60 * time.Second)

	var ticket = fmt.Sprintf("Goodbye %v", name)
	fmt.Println("*********************")
	fmt.Printf("Sending threaded message: %v \n", ticket)
	fmt.Println("*********************")

	wg.Done()
}

// Other interesting links
// https://gobyexample.com/command-line-flags
// https://gobyexample.com/command-line-subcommands

func helloHttpServer(w http.ResponseWriter, req *http.Request) {

	// Create a context to handle the request
	ctx := req.Context()
	fmt.Println("helloHttpServer: Hello handler started")
	defer fmt.Println("helloHttpServer handler ended")

	select {
	case <-time.After(10 * time.Second):
		fmt.Fprintf(w, "Hello\n")
	case <-ctx.Done():

		err := ctx.Err()
		fmt.Println("helloHttpServer: ", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}

}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
