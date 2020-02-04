package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 4 {
		fmt.Println("Input not valid")
		fmt.Println("Usage: ")
		fmt.Println(os.Args[0], " <subscription> <region> <file> <sample to run>")
		return
	}
	subscription := args[0]
	region := args[1]
	file := args[2]
	sample := args[3]

	switch sample {
	case "1":
		listenOnce(subscription, region, file)
	case "2":
		kws(subscription, region, file)
	}
}