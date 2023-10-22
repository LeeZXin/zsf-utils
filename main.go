package main

import "fmt"

func main() {
	c := make(chan string)
	close(c)
	for {
		select {
		case <-c:
			fmt.Println("xxxx")
			return
		}
	}
}
