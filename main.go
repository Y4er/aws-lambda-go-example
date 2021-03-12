package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"strings"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	parameters := request.PathParameters
	fmt.Println(len(parameters))
	for p := range parameters {
		fmt.Println(p)
	}
	path := request.Path
	index := strings.Index(
		path,
		"/img/uploads",
	)
	body := fmt.Sprintf(
		"Path:%s\nID:%s\nIMG:%s",
		path,
		request.QueryStringParameters["id"],
		path[index:],
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
