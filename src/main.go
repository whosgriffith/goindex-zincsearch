package main

import (
	"encoding/json"
	"fmt"
	"indexing-zincsearch/helpers"
	"io/fs"
	"log"
	"os"
)

func main() {
	//const path string = "/Users/risker/Documents/enron_mail_20110402/maildir"
	const path string = "/Users/risker/Documents/enron_mail_full/maildir"
	//const path string = "/Users/risker/Documents/enron_mail_full/maildir/baughman-d/calendar/19."

	var emailsSlice []helpers.Email

	successCount := 0
	partialSuccessCount := 0
	failCount := 0

	fileSystem := os.DirFS(path)
	fs.WalkDir(fileSystem, ".", func(filePath string, d fs.DirEntry, err error) error {
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
				fmt.Println(filePath)
				failCount++
			}
		}
		return nil
	})

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
