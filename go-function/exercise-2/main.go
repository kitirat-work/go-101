package main

import "os"

func main() {
	count, err := fileLen("./exercise-2/main.go")
	if err != nil {
		println("Error:", err.Error())
	} else {
		println("File length:", count)
	}

}

func fileLen(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	bytes := make([]byte, 1024)
	count, err := file.Read(bytes)
	return count, err
}
