package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var queues map[string][]string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("port not set")
		return
	}
	port := os.Args[1]
	queues = make(map[string][]string)
	http.HandleFunc("/", handle)
	http.ListenAndServe(":"+port, nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	queue := r.URL.Path
	if r.Method == http.MethodPut {
		msg := r.FormValue("v")
		if msg == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		putMsg(queue, msg)
	}
	if r.Method == http.MethodGet {
		timeout, _ := strconv.ParseInt(r.FormValue("timeout"), 0, 64)
		msg, err := getMsg(queue)
		if err != nil {
			if timeout > 0 {
				for range time.Tick(time.Second) {
					if msg, err = getMsg(queue); err == nil {
						w.Write([]byte(msg))
						return
					}
					if timeout--; timeout == 0 {
						break
					}
				}
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(msg))
	}
}

func getMsg(queue string) (string, error) {
	if len(queues[queue]) > 0 {
		msg := queues[queue][0]
		queues[queue] = queues[queue][1:]
		return msg, nil
	}
	return "", errors.New("empty")
}

func putMsg(queue, msg string) {
	queues[queue] = append(queues[queue], msg)
}
