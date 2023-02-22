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

var profile = flag.String("profile", "", "write cpu profile to `file`")
var directory = flag.String("directory", "", "Path to mails folders.")

func main() {

	flag.Parse()
	if *profile != "" {
		f, err := os.Create(*profile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatal("Error closing file: ", err)
			}
		}(f)
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	if *directory == "" {
		fmt.Println("Error: Need -directory flag with path to start indexation")
		os.Exit(2)
	}

	var successCount, partialSuccessCount, failCount int
	var emailsSlice []helpers.Email

	fmt.Println("Processing files from directory...")
	fileSystem := os.DirFS(*directory)
	err := fs.WalkDir(fileSystem, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if !d.IsDir() && d.Name() != ".DS_Store" {
			fullFilePath := fmt.Sprintf("%s/%s", *directory, filePath)
			data, _ := os.ReadFile(fullFilePath)

			email, result := helpers.ProccessEmailFile(data)
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

	fmt.Println("Encode JSON...")
	zincBodyBulkV2 := helpers.ZincBodyBulkV2{Index: "emails", Records: emailsSlice}
	data, err := json.Marshal(zincBodyBulkV2)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Println("Indexing data...")
	helpers.ZincIngest(data)

	fmt.Println("Summary:")
	fmt.Printf("%d total files\n", successCount+partialSuccessCount+failCount)
	fmt.Printf("- %d successfull\n", successCount)
	fmt.Printf("- %d partially successfull\n", partialSuccessCount)
	fmt.Printf("- %d failed\n", failCount)
}
