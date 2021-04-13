package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:9090/chunks"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer response.Body.Close()
	for i := 0; i < 4; i++ {
		p := make([]byte, 512)
		response.Body.Read(p)
		fmt.Println(string(p))
		time.Sleep(time.Second)
	}
}
