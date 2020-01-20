package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type counters struct {
	sync.Mutex
	view  int
	click int
}

var (
	c = counters{}

	content = []string{"sports", "entertainment", "business", "education"}

	c_arr [4]counters
	reqIP []string

	logFile        = "log.txt"
	requestCounter = 0
	secondsCounter = 0
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf(">>>>>>>>> Welcome >>>>>>>>>\n")
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var randContent = rand.Intn(len(content))
	data := content[randContent]

	c_arr[randContent].Lock()
	c_arr[randContent].view++
	c_arr[randContent].Unlock()

	// fmt.Printf("\n>>>>>>>>> VIEWED >>>>>>>>>\n")

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(data, randContent)
	}

	var outputLog string
	var t = time.Now().Format("2006-01-02 15:04")

	outputLog = data + " " + t + " {view: " + strconv.Itoa(c_arr[randContent].view) + " clicks: " + strconv.Itoa(c_arr[randContent].click) + "}\n"
	fmt.Print(outputLog)
	fmt.Fprint(w, outputLog)
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string, randContent int) error {
	c_arr[randContent].Lock()
	c_arr[randContent].click++
	c_arr[randContent].Unlock()
	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	// capturedIP := getIP(r)
	// print(capturedIP, "\n")

	requestCounter++

	if !isAllowed() {
		fmt.Printf("TOO MANY REQUESTS\n")
		w.WriteHeader(429)
		return
	}

	fmt.Printf("Request Accepted\n")

	fmt.Fprint(w, "Statistics Collected up to now \n")
	fmt.Fprint(w, "------------------------------ \n \n")

	for i := 0; i < 4; i++ {

		var outputLog string
		var t = time.Now().Format("2006-01-02 15:04")
		outputLog = content[i] + " " + t + " {view: " + strconv.Itoa(c_arr[i].view) + " clicks: " + strconv.Itoa(c_arr[i].click) + "}\n"

		fmt.Fprint(w, outputLog)
	}
}

// func requestAdder(ip string) {

// }

func isAllowed() bool {
	if requestCounter > 5 {
		return false
	} else {
		return true
	}
}

func uploadCounters() error {
	return nil
}

// func getIP(r *http.Request) string {
// 	forwarded := r.Header.Get("X-FORWARDED-FOR")
// 	print(forwarded)
// 	if forwarded != "" {
// 		return forwarded
// 	}
// 	return r.RemoteAddr
// }

func writeFile(input string) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	datawriter.WriteString(input)

	datawriter.Flush()
	file.Close()
}

func tracker() {
	for {
		if secondsCounter == 10 {
			secondsCounter = 0
			requestCounter = 0
			fmt.Print("----------RESET---------------- SECONDS ACTIVE: ", secondsCounter, " | requestCounter: ", requestCounter, "\n")
			time.Sleep(1000 * time.Millisecond)
		} else {
			secondsCounter++
			fmt.Print("------------------------------- SECONDS ACTIVE: ", secondsCounter, " | requestCounter: ", requestCounter, "\n")
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func logger() {
	for {
		fmt.Print("LOGGING\n")

		os.Remove(logFile) // resetting the log file to take new entries
		writeFile("LOGS FOR CLICKS AND VIEWS...\n")

		for i := 0; i < 4; i++ {

			var outputLog string
			var t = time.Now().Format("2006-01-02 15:04")
			outputLog = content[i] + " " + t + " {view: " + strconv.Itoa(c_arr[i].view) + " clicks: " + strconv.Itoa(c_arr[i].click) + "}\n"

			writeFile(outputLog)
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	fmt.Printf("Outputting in terminal as a guide\n\n")

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)

	go logger()
	go tracker()

	log.Fatal(http.ListenAndServe(":8070", nil))
}
