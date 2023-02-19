package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
)

func ParseEmail(data []uint8) (email Email, result string) {
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
