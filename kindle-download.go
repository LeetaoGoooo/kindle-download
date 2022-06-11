package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"kindle-download/formatter"
	"kindle-download/tools"
	"os"
	"strings"
)

func main() {
	fmt.Println("***************************************************\n    欢迎使用 Kindle Downloader 工具    \n    项目开源在 github.com/leetaogoooo/kindle-download    \n    欢迎提交问题和建议    \n    别忘了 star 一下哦    \n***************************************************")

	if !checkCookieFile() {
		os.Exit(1)
	}

	cookie := ReadCookieFromFile()

	var workerNum int
	var fileDir string
	var cn bool

	flag.IntVar(&workerNum, "worker", 5, "最大协程数,默认为5")
	flag.StringVar(&fileDir, "dir", "ebooks", "文件保存目录,默认为当前目录下的ebooks")
	flag.BoolVar(&cn, "cn", true, "是否下载中国版的书籍,默认为 true")
	flag.Parse()

	var config formatter.Config = formatter.NewConfig(workerNum, fileDir, cookie, cn)

	kindle := tools.NewKindleClient(config)
	kindle.GetCsrfToken()

	// 获取设备列表
	devices := kindle.GetDevices()

	if devices.GetDevices.Count == 0 {
		fmt.Println("没有可用的设备")
		return
	}

	fsn := devices.GetDevices.Devices[0].DeviceSerialNumber
	deviceType := devices.GetDevices.Devices[0].DeviceType
	customerId := devices.GetDevices.Devices[0].CustomerId

	fileTypes := []string{"Ebook", "KindlePDoc"}

	ch := make(chan bool, config.Common.WorkerNum)

	for _, fileType := range fileTypes {
		startIndex := 0
		batchSize := 100
		totalContentCount := 0
		// 获取电子书列表
		resp := kindle.GetBookList(formatter.NewReqBookList(startIndex, batchSize, totalContentCount, fileType))
		Items := resp.OwnershipData.Items
		if resp.OwnershipData.HasMoreItems {
			startIndex += 18
			totalContentCount = resp.OwnershipData.NumberOfItems
			resp = kindle.GetBookList(formatter.NewReqBookList(startIndex, batchSize, totalContentCount, fileType))
			Items = append(Items, resp.OwnershipData.Items...)
		}

		// 获取电子书下载链接

		fTyppe := "EBOK"
		fTyppeMsg := "电子书"
		if fileType == "KindlePDoc" {
			fTyppe = "PDOC"
			fTyppeMsg = "个人文档"

		}

		fmt.Println("成功获取到【", fTyppeMsg, "】下载链接，总计", len(Items), "个")

		for _, item := range Items {

			url := fmt.Sprintf("https://cde-ta-g7g.amazon.com/FionaCDEServiceEngine/FSDownloadContent?type=%s&key=%s&fsn=%s&device_type=%s&customerId=%s", fTyppe, item.Asin, fsn, deviceType, customerId)

			if config.Common.CN {
				url += "&authPool=AmazonCN"
			}

			downloadFile := formatter.FileDownload{
				FileName: html.UnescapeString(item.Title),
				Url:      url,
			}
			ch <- true
			go func() {
				defer func() {
					<-ch
				}()

				kindle.DownloadFile(downloadFile)
			}()
		}
	}
}

// 检查 cookie 文件是否存在
// 并创建文件
func checkCookieFile() bool {
	_, err := os.Stat("cookie.txt")
	if os.IsNotExist(err) {
		fmt.Println("cookie.txt 文件不存在，将自动创建，将浏览器的 cookie 粘贴到 cookie.txt 文件中,重新运行文件即可")
		file, err := os.Create("cookie.txt")
		if err != nil {
			fmt.Println("创建 cookie.txt 文件失败,请手动创建 cookie.txt 文件,并将浏览器的 cookie 粘贴到 cookie.txt 文件中")
		}
		defer file.Close()
		return false
	}
	return true
}

// 从配置文件中获取 cookie
func ReadCookieFromFile() string {
	f, err := ioutil.ReadFile("cookie.txt")
	if err != nil {
		panic("读取 cookie 文件失败")
	}
	cookie := string(f)
	if len(strings.Trim(cookie, "")) == 0 {
		panic("cookie 为空, 请检查 cookie.txt 文件")
	}
	return cookie
}
