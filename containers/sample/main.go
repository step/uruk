package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		if err != nil {
			return err
		}
		*files = append(*files, path)
		return nil
	}
}

func main() {
	var files []string
	sourcePath := "/source"
	fmt.Println("source path is", sourcePath)
	err := filepath.Walk(sourcePath, visit(&files))
	if err != nil {
		fmt.Println(err)
	}

	result := make(map[string][]string)
	result["files"] = files
	bytes, _ := json.Marshal(result)

	resultPath := "/results"
	resultsFilename := filepath.Join(resultPath, "result.json")
	os.MkdirAll("/results", 0777)
	resultsFile, _ := os.OpenFile(resultsFilename, os.O_CREATE|os.O_WRONLY, 0777)
	defer resultsFile.Close()
	fmt.Fprintf(resultsFile, "%s", string(bytes))
}
