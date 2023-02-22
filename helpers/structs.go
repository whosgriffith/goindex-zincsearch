package helpers

import (
	"net/mail"
	"time"
)

type Email struct {
	MessageId string
	Subject   string
	Date      time.Time
	From      []*mail.Address
	To        []*mail.Address
	Content   string
}

type ZincBodyBulkV2 struct {
	Index   string
	Records []Email
}
