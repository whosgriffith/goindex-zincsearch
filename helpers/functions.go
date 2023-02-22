package helpers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
)

func ProccessEmailFile(data []uint8) (email Email, result string) {
	msg, err := mail.ReadMessage(bytes.NewBuffer(data))
	if err != nil {
		return Email{}, "fail"
	}
	header := msg.Header
	result = "success"

	address, fromErr := header.AddressList("From")
	if fromErr != nil {
		address = nil
		result = "partialSuccess"
	}

	addresses, toErr := header.AddressList("To")
	if toErr != nil {
		addresses = nil
		result = "partialSuccess"
	}

	dateTime, dateErr := header.Date()
	if dateErr != nil {
		fmt.Println(dateErr)
		return Email{}, "fail"
	}

	body, bodyErr := io.ReadAll(msg.Body)
	if bodyErr != nil {
		body = nil
		result = "partialSuccess"
	}

	email = Email{
		MessageId: header.Get("Message-Id"),
		Subject:   header.Get("Subject"),
		Date:      dateTime,
		From:      address,
		To:        addresses,
		Content:   string(body),
	}

	return email, result
}

func ZincIngest(data []byte) {
	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulkv2", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "admin")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
