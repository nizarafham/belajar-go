package lib

import (
	"log"
	"os"
	"github.com/xendit/xendit-go/v3"
)

var XenditClient *xendit.APIClient

func InitXendit() {
	apiKey := os.Getenv("XENDIT_API_KEY")
	if apiKey == "" {
		log.Println("PERINGATAN: XENDIT_API_KEY tidak diatur. Pembayaran tidak akan berfungsi.")
		return
	}
	XenditClient = xendit.NewClient(apiKey)
}