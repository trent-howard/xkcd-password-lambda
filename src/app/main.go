package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	w "aws-lambda-go/wordlist"
)

type App struct {
	id string
}

func newApp(id string) *App {
	return &App{
		id: id,
	}
}

type ErrorBody struct {
	ErrorMsg string `json:"error,omitempty"`
}

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Methods":     "GET",
		"Access-Control-Allow-Credentials": "true",
	}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

func (app *App) Handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// parse query string and validate
	var length = 5
	lengthParam, ok := request.QueryStringParameters["length"]
	if ok {
		int, err := strconv.Atoi(lengthParam)
		if err != nil {
			return apiResponse(http.StatusBadRequest,
				ErrorBody{"length param must be a number"},
			)
		}
		length = int
	}

	if length < 1 || length > 1000 {
		return apiResponse(http.StatusBadRequest, ErrorBody{"length must be a number between 1 and 1000"})
	}

	// generate unique random indexes we'll use to lookup words from the word list
	// less secure than if we allow words more than once, but I prefer the aesthetic 🤣
	uniqueRandomInts := make(map[int]bool)
	for len(uniqueRandomInts) < length {
		i := rand.Int() % len(w.WordList)
		if !uniqueRandomInts[i] {
			uniqueRandomInts[i] = true
		}
	}

	// build the password and send it!
	pieces := []string{}
	for idx := range uniqueRandomInts {
		pieces = append(pieces, w.WordList[idx])
	}
	result := strings.Join(pieces, "-")

	responseBody := map[string]string{
		"password": result,
	}

	return apiResponse(http.StatusOK, responseBody)
}

func main() {
	app := newApp("xkcd-password-gen")

	lambda.Start(app.Handler)
}
