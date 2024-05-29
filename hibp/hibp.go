package hibp

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func CheckMails(emails ...string) {
	client := &http.Client{}
	for _, mail := range emails {
		req, err := buildHIBPRequest(mail)
		if err != nil {
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		fmt.Println(resp.Status)
		time.Sleep(6 * time.Second)
	}
}

func buildHIBPRequest(mail string) (*http.Request, error) {
	apiKey := os.Getenv("HIBP_API_KEY")
	req, err := http.NewRequest("GET", "https://haveibeenpwned.com/api/v3/breachedaccount/"+mail, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("hibp-api-key", apiKey)

	return req, nil
}
