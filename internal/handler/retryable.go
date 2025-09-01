package handler

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func NewRetryableClient() *retryablehttp.Client {
	client := retryablehttp.NewClient()
	
	client.RetryMax = 3
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 5 * time.Second
	return  client
}