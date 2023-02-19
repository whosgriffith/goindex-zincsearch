package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"indexing-zincsearch/helpers"
	"io/fs"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpu.prof", "", "write cpu profile to `file`")

func main() {
	// Profiling
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatal("Error closing file: ", err)
			}
		}(f) // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	const path string = "directory"

	var emailsSlice []helpers.Email

	successCount := 0
	partialSuccessCount := 0
	failCount := 0

	fileSystem := os.DirFS(path)
	err := fs.WalkDir(fileSystem, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if !d.IsDir() && d.Name() != ".DS_Store" {
			fullFilePath := fmt.Sprintf("%s/%s", path, filePath)
			data, _ := os.ReadFile(fullFilePath)
			email, result := helpers.ParseEmail(data)

			switch result {
			case "success":
				emailsSlice = append(emailsSlice, email)
				successCount++
			case "partialSuccess":
				emailsSlice = append(emailsSlice, email)
				partialSuccessCount++
			case "fail":
				failCount++
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Error walking trough directory: ", err)
	}

	fmt.Printf("%d total files:\n", successCount+partialSuccessCount+failCount)
	fmt.Printf("- %d successfull\n", successCount)
	fmt.Printf("- %d partially successfull\n", partialSuccessCount)
	fmt.Printf("- %d failed\n", failCount)

	file, err := json.MarshalIndent(emailsSlice, "", " ")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		_ = os.WriteFile("emails.json", file, 0644)
	}
}
