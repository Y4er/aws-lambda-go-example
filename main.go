package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/issue9/watermark"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		Timeout: 5 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://y4er.com"+imgpath, nil)
	req.Header.Set("User-Agent", "netlify")
	req.Header.Set("Referer", "https://y4er.com"+imgpath)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		body = base64.StdEncoding.EncodeToString([]byte(err.Error()))
	} else {
		bs, _ := ioutil.ReadAll(resp.Body)
		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("%s", timestamp)
		file, _ := os.Create(filename)
		io.Copy(file, bytes.NewReader(bs))
		defer file.Close()
		w, err := watermark.New("watermark.png", 2, watermark.BottomRight)
		if err != nil {
			body = err.Error()
			fmt.Println(body)
		} else {
			err := w.MarkFile(filename)
			if err != nil {
				body = err.Error()
				fmt.Println(body)
				//body = base64.StdEncoding.EncodeToString([]byte(err.Error()))
			} else {
				content, _ := ioutil.ReadFile(filename)
				body = base64.StdEncoding.EncodeToString(content)
				fmt.Println(body)
			}
		}
	}

	dir, _ := os.Getwd()
	body = dir
	files, _ := ioutil.ReadDir("../")
	for _, f := range files {
		fmt.Println(f.Name())
		body += f.Name()
	}
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
