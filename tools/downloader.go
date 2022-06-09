package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"kindle-download/formatter"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/cheggaaa/pb"
)

type KindleClient struct {
	Client    http.Client
	CsrfToken string
	Config    formatter.Config
}

func NewKindleClient(config formatter.Config) *KindleClient {

	_, err := os.Stat(config.Common.FileDir)
	if os.IsNotExist(err) {
		os.MkdirAll(config.Common.FileDir, 0755)
	}

	return &KindleClient{
		Client: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Config: config,
	}
}

func (kindleClient *KindleClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.Add("Cookie", kindleClient.Config.Common.Cookie)
	return req, nil
}

func (kindleClient *KindleClient) GetCsrfToken() {
	url := "https://www.amazon.cn/hz/mycd/myx#/home/content/booksAll/dateDsc/"
	if !kindleClient.Config.Common.CN {
		url = kindleClient.Config.COM.ListUrl
	}

	req, err := kindleClient.NewRequest("GET", url, nil)

	if err != nil {
		panic(fmt.Errorf("请求异常: %v", err))
	}
	resp, err := kindleClient.Client.Do(req)

	if err != nil {
		panic(fmt.Errorf("请求异常: %v", err))
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(fmt.Errorf("读取响应异常: %v", err))
	}
	r, _ := regexp.Compile("var csrfToken = \"(.+)\"")

	matches := r.FindStringSubmatch(string(respBody))

	if matches == nil {
		panic("获取csrfToken异常,请检查 cookie 是否正确")
	}

	kindleClient.CsrfToken = matches[1]
}

func (kindleClient *KindleClient) GetDevices() (result *formatter.DevicesResp) {
	body := url.Values{}
	body.Set("data", string("{\"param\":{\"GetDevices\":{}}}"))
	body.Set("csrfToken", kindleClient.CsrfToken)

	url := "https://www.amazon.cn/hz/mycd/ajax"

	if !kindleClient.Config.Common.CN {
		url = kindleClient.Config.COM.ListUrl
	}

	req, err := kindleClient.NewRequest("POST", url, strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		fmt.Println("构建请求异常", err)
		return
	}
	resp, err := kindleClient.Client.Do(req)

	if err != nil {
		fmt.Println("获取设备链接异常", err)
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应异常", err)
		return
	}
	json.Unmarshal(respBody, &result)
	return
}

// GetBookList 获取电子书列表
func (kindleClient *KindleClient) GetBookList(params formatter.ReqBookList) (result *formatter.RespBookList) {
	body := url.Values{}
	data, err := json.Marshal(params.Data)

	if err != nil {
		fmt.Println("序列化请求参数异常", err)
		return
	}

	body.Set("data", string(data))
	body.Set("csrfToken", kindleClient.CsrfToken)

	url := "https://www.amazon.cn/hz/mycd/ajax"

	if !kindleClient.Config.Common.CN {
		url = kindleClient.Config.COM.ListUrl
	}

	req, err := kindleClient.NewRequest("POST", url, strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		fmt.Println("构建请求异常", err)
		return
	}
	resp, err := kindleClient.Client.Do(req)

	if err != nil {
		fmt.Println("请求文档下载链接异常", err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应异常", err)
		return
	}
	json.Unmarshal(respBody, &result)
	return
}

//	GetDownloadLink 获取文档下载链接
func (kindleClient *KindleClient) GetDownloadLink(params formatter.DownLoadBookReq) (result *formatter.DownloadViaUSB) {
	body := url.Values{}
	data, err := json.Marshal(params.Data)

	if err != nil {
		fmt.Println("序列化请求参数异常", err)
		return
	}

	body.Set("data", string(data))
	body.Set("csrfToken", kindleClient.CsrfToken)

	url := "https://www.amazon.cn/hz/mycd/ajax"

	if !kindleClient.Config.Common.CN {
		url = kindleClient.Config.COM.ListUrl
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		fmt.Println("构建请求异常", err)
		return
	}

	resp, err := kindleClient.Client.Do(req)

	if err != nil {
		fmt.Println("请求文档下载链接异常", err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应异常", err)
		return
	}
	json.Unmarshal(respBody, &result)
	return
}

func (kindleClient *KindleClient) DownloadFile(fileDownload formatter.FileDownload) {

	req, err := kindleClient.NewRequest("GET", fileDownload.Url, nil)
	if err != nil {
		fmt.Println("构建请求异常", err)
		return
	}

	resp, err := kindleClient.Client.Do(req)

	if err != nil {
		fmt.Println("请求文档下载链接异常", err)
		return
	}

	// 文件后缀从 header 中获取 content-disposition
	if resp.StatusCode == 302 {
		fileDownload = formatter.FileDownload{
			Url:      resp.Header.Get("Location"),
			FileName: fileDownload.FileName,
		}
		kindleClient.DownloadFile(fileDownload)
	}
	contentDisposition := resp.Header.Get("Content-Disposition")
	r, _ := regexp.Compile(`filename\*=UTF-8''(.+)`)
	matches := r.FindStringSubmatch(contentDisposition)
	if len(matches) == 0 {
		return
	}

	fileName, err := url.QueryUnescape(matches[1])

	if err != nil {
		return
	}

	fmt.Printf("开始下载 【%s】\n", fileDownload.FileName)

	fileSize := int(resp.ContentLength)

	file, err := os.Create(fmt.Sprintf("%s/%s", kindleClient.Config.Common.FileDir, fileName))
	if err != nil {
		fmt.Println("创建文件异常", err)
		return
	}
	defer file.Close()
	bar := pb.New(fileSize).SetUnits(pb.U_BYTES)
	bar.Start()
	rd := bar.NewProxyReader(resp.Body)
	io.Copy(file, rd)

}
