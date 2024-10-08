// Check every path level is exist or not
// If not exist, create it
// If exist, do nothing
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	path := "/a/b/c/d/e"
	createPath(path)

	path = "v/w/x/y/z"
	recursiveCreatePath(".", path)
}

func createPath(path string) {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			dir := path[:i]
			fmt.Println(dir)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.Mkdir(dir, 0755)
				fmt.Println("Create path: ", dir)
			}
		}
	}
}

func recursiveCreatePath(basePath, path string) {
	if len(path) == 0 {
		return
	}

	// base case
	if !strings.Contains(path, "/") {
		dir := basePath + "/" + path
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
			fmt.Println("Create path: ", dir)
		}
		return
	}

	// recursive case
	idx := strings.Index(path, "/")
	curDir := basePath + "/" + path[:idx]
	if _, err := os.Stat(curDir); os.IsNotExist(err) && len(curDir) > 0 {
		os.Mkdir(curDir, 0755)
		fmt.Println("Create path: ", curDir)
	}
	recursiveCreatePath(curDir, path[idx+1:])
}
