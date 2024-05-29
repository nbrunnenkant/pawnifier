package simplelogin

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
)

const itemsPerPage = 20

type aliasResponse struct {
	Email string `json:"email"`
}

type aliasesResponse struct {
	Aliases []aliasResponse `json:"aliases"`
}

type StatsResponse struct {
	NumberAliases int `json:"nb_alias"`
}

type SimpleloginService struct {
}

func NewSimpleloginService() *SimpleloginService {
	return &SimpleloginService{}
}

func (sl *SimpleloginService) GetMails() []string {
	emails := make([]string, 0)
	client := http.Client{}

	for i := 0; i < calculateNumberOfCalls(); i++ {
		req, err := buildSimpleloginRequest("/v2/aliases?page_id=" + strconv.Itoa(i))

		if err != nil {
			fmt.Println(err)
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}

		bytes, err := io.ReadAll(resp.Body)
		aliasResp := &aliasesResponse{}

		err = json.Unmarshal(bytes, aliasResp)
		if err != nil {
			fmt.Println(err)
		}

		for _, alias := range aliasResp.Aliases {
			emails = append(emails, alias.Email)
		}
	}
	return emails
}

func calculateNumberOfCalls() int {
	req, err := buildSimpleloginRequest("stats")
	if err != nil {
		return 0
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0
	}

	stats := &StatsResponse{}
	err = json.Unmarshal(bytes, stats)
	if err != nil {
		return 0
	}

	numberOfCalls := math.Ceil(float64(stats.NumberAliases) / float64(itemsPerPage))
	return int(numberOfCalls)
}

func buildSimpleloginRequest(endpoint string) (*http.Request, error) {
	apiKey := os.Getenv("SIMPLELOGIN_API_KEY")
	req, err := http.NewRequest("GET", "http://app.simplelogin.io/api/"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authentication", apiKey)

	return req, nil
}
