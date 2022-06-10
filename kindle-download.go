package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"kindle-download/formatter"
	"kindle-download/tools"

	"gopkg.in/yaml.v2"
)

func main() {
	fmt.Println("***************************************************\n    欢迎使用 Kindle Downloader 工具    \n    项目开源在 github.com/leetaogoooo/kindle-download    \n    欢迎提交问题和建议    \n    别忘了 star 一下哦    \n***************************************************")

	var config formatter.Config
	File, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("读取配置文件失败 #%v", err)
		return
	}
	err = yaml.Unmarshal(File, &config)
	if err != nil {
		fmt.Printf("解析失败: %v", err)
		return
	}

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
