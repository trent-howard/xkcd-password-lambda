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

func (app *App) Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	// parse query string and validate
	var length = 5
	lengthParam, ok := request.QueryStringParameters["length"]
	if ok {
		int, err := strconv.Atoi(lengthParam)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers: map[string]string{"Content-Type": "application/json"},
				Body: `{"error": "length param must be a number"}`,
			}, nil 
		}
		length = int
	}

	if length < 1 || length > 1000 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{"Content-Type": "application/json"},
			Body: `{"error": "length must be a number between 1 and 1000"}`,
		}, nil 
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

	responseJson, err := json.Marshal(responseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{"Content-Type": "application/json"},
			Body: `{"error": "internal server error"}`,
		}, nil 
	}
	
	response := events.APIGatewayProxyResponse{
		Body: string(responseJson),
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Methods": "GET",
			"Access-Control-Allow-Credentials": "true",
		},
	}
	return response, nil
}

func main() {
	id:= "someRandomString"
	app := newApp(id)

	lambda.Start(app.Handler)
}