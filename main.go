package main

import (
	"fmt"
)

var eula = ""

var ovfTemplate = ""

var numbs = 10

func main() {
	a := 0
	for a = 0; a <= numbs; {
		fmt.Println(a)
		a++
	}

	fmt.Printf(sha256("http://test.com \n"))
	fmt.Println(getVmdkFiles("in the list \n"))
	fmt.Println(streamOptimizeVmdkFiles("in the list \n"))
	fmt.Println(createOva("github.com \n", "gitlab.com \n", "data.db \n"))
	fmt.Println(createOvf("github.com \n", "gitlab.com \n", "data.db \n"))
}

func sha256(path string) string {
	return path
}

func createOva(ovaPath string, ovfPath string, ovaFiles string) string {
	return ovaPath + ovfPath + ovaFiles
}

func createOvf(path string, data string, ovfTemplate string) string {
	return path + data + ovfTemplate
}

func createOvaManifest(path string, infilePaths string) string {
	return path + infilePaths
}

func getVmdkFiles(inList string) string {
	return inList
}

func streamOptimizeVmdkFiles(inList string) string {
	return inList
}
