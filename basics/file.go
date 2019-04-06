package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	fin, err := os.Open("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fin.Close()

	fout, err := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fout.Close()

	reader := bufio.NewReader(fin)
	writer := bufio.NewWriter(fout)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		writer.WriteString(line)
	}

	writer.Flush()
}
