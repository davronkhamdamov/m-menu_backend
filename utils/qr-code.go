package utils

import (
	"fmt"
	"log"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCodeForTable(tableID string) string {
	url := fmt.Sprint(tableID)
	filename := fmt.Sprintf("qrcodes/table_%s.png", tableID)
	err := qrcode.WriteFile(url, qrcode.Medium, 256, filename)
	if err != nil {
		log.Println("QR generation error:", err)
	}
	return filename
}
