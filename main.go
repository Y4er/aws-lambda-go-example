package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	parameters := request.PathParameters
	fmt.Println(len(parameters))
	for p := range parameters {
		fmt.Println(p)
	}
	path := request.Path
	imgpath := strings.ReplaceAll(path, "/.netlify/functions/test-lambda", "")
	body := fmt.Sprintf(
		"Path:%s\nID:%s\nIMG:%s",
		path,
		request.QueryStringParameters["id"],
		imgpath,
	)
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://y4er.com"+imgpath, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	req.Header.Set("Referer", "https://y4er.com"+imgpath)

	resp, err := client.Do(req)
	if err != nil {
		body = err.Error()
	} else {
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
