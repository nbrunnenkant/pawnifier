package hibp

import (
	"net/http"
	"os"
	"time"
)

type HIBPService struct {
	queue    chan *http.Request
	Response chan bool
}

func NewHIBPService() *HIBPService {
	service := &HIBPService{queue: make(chan *http.Request), Response: make(chan bool)}
	go service.processQueue()

	return service
}

func (hibp *HIBPService) processQueue() {
	for {
		select {
		case activeMail := <-hibp.queue:
			hibp.Response <- checkMail(activeMail.PathValue("mail"))
			time.Sleep(time.Second * 6)
		}
	}
}

func (hibp *HIBPService) AddMail(req *http.Request) {
	hibp.queue <- req
}

func checkMail(email string) bool {
	client := &http.Client{}
	req, err := buildHIBPRequest(email)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	return resp.StatusCode == 404
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
