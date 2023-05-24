package main

import (
	"fmt"
	"os"
)

func WriteToFile(text string) {
  file, err := os.OpenFile("example.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  if err != nil {
    fmt.Println("Error opening file:", err)
    return
  }
	
	defer file.Close()

  data := []byte(text + "\n")
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File written successfully.")
}

func main() {
  WriteToFile("this is a test!!!")
  WriteToFile("2")
  WriteToFile("3")
}

