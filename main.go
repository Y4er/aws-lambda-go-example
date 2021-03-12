package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	parameters := request.PathParameters
	for p := range parameters {
		fmt.Println(p)
	}
	body := fmt.Sprintf(
		"Path:%s\nID:%s",
		request.Path,
		request.QueryStringParameters["id"],
	)
	return events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"Content-Type": "text/plain"},
		MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
		Body:              body,
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
