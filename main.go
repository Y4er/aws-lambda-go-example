package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
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
	imgpath:=path[index:]
	body := fmt.Sprintf(
		"Path:%s\nID:%s\nIMG:%s",
		path,
		request.QueryStringParameters["id"],
		imgpath,
	)
	resp, err := http.Get("https://y4er.com" + imgpath)
	if err != nil {
		body = err.Error()
	}else{
		bytes, _ := ioutil.ReadAll(resp.Body)
		body = string(bytes)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"Content-Type": "image/png"},
		MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
		Body:              body,
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
