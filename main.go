package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ci", LighthouseCi).Methods("POST")
	router.HandleFunc("/chrome", HeadlessChrome).Methods("POST")

	log.Printf("Listening on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type Lighthouse struct {
	Output string `json:"output"`
	Url    string `json:"url"`
}

func LighthouseCi(w http.ResponseWriter, r *http.Request) {
	var chromeFlags = "--chrome-flags=--headless --no-sandbox"
	if len(os.Getenv("HTTPPROXY")) > 0 {
		chromeFlags = "--chrome-flags=\"--headless --no-sandbox --proxy-server=" + os.Getenv("HTTPPROXY") + "\""
	}

	var lighthouse Lighthouse
	json.NewDecoder(r.Body).Decode(&lighthouse)

	if !(len(lighthouse.Output) > 0) {
		lighthouse.Output = "json"
	}

	args := []string{chromeFlags, "--output=" + lighthouse.Output, "--emulated-form-factor=mobile", "--port=0", lighthouse.Url}
	cmd := exec.Command("lighthouse", args...)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	fmt.Println(args)
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, string(cmdOutput.Bytes()))
	}
}

func HeadlessChrome(w http.ResponseWriter, r *http.Request) {}
