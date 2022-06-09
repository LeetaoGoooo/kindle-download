package formatter

// 配置文件
type Config struct {
	Common struct {
		WorkerNum int    `yaml:"workerNum"`
		FileDir   string `yaml:"fileDir"`
		Cookie    string `yaml:"cookie"`
		CN        bool   `yaml:"cn"`
	} `yaml:"common"`
	CN struct {
		ListUrl string `yaml:"listUrl"`
		AjaxUrl string `yaml:"ajaxUrl"`
	} `yaml:"cn"`
	COM struct {
		ListUrl string `yaml:"listUrl"`
		AjaxUrl string `yaml:"ajaxUrl"`
	} `yaml:"com"`
}

// 设备列表结果
type DevicesResp struct {
	GetDevices struct {
		Count   int `json:"count"`
		Devices []struct {
			DeviceSerialNumber string `json:"deviceSerialNumber"`
			DeviceType         string `json:"deviceType"`
			CustomerId         string `json:"customerId"`
		} `json:"devices"`
	}
}

/// 请求获取书籍列表的结构体
type ReqBookList struct {
	Data map[string]interface{} `json:"data"`
}

type OwnershipData struct {
	HasMoreItems  bool                `json:"hasMoreItems"`
	NumberOfItems int                 `json:"numberOfItems"`
	Success       bool                `json:"success"`
	Items         []OwnershipDataItem `json:"items"`
}

type OwnershipDataItem struct {
	Asin  string `json:"asin"`
	Title string `json:"title"`
}

/// 电子书书籍列表
type RespBookList struct {
	OwnershipData OwnershipData `json:"OwnershipData"`
}

/// 请求下载书籍的结构体
type DownLoadBookReq struct {
	Data     map[string]interface{}
	FileName string
}

type DownloadBookResp struct {
	DownloadViaUSB DownloadViaUSB `json:"DownloadViaUSB"`
}

type DownloadViaUSB struct {
	Success bool   `json:success`
	URL     string `json:URL`
}

type FileDownload struct {
	Url      string `json:"url"`
	FileName string `json:"fileName"`
}

func NewDownLoadBookReq(asin string, targetDevice string, fileName string) DownLoadBookReq {
	return DownLoadBookReq{
		Data: map[string]interface{}{
			"param": map[string]interface{}{
				"DownloadViaUSB": map[string]string{"contentName": asin, "encryptedDeviceAccountId": targetDevice, "originType": "Purchase"},
			},
		},
		FileName: fileName,
	}
}

func NewReqBookList(startIndex, batchSize, totalContentCount int, fileType string) ReqBookList {
	OwnershipData := map[string]interface{}{
		"sortOrder":   "DESCENDING",
		"sortIndex":   "DATE",
		"startIndex":  startIndex,
		"batchSize":   batchSize,
		"contentType": fileType,
		"itemStatus":  [...]string{"Active"},
	}

	if fileType == "Ebook" {
		OwnershipData["originType"] = [...]string{"Purchase"}
		OwnershipData["showSharedContent"] = true
	} else {
		OwnershipData["batchSize"] = 18
		OwnershipData["isExtendedMYK"] = false
	}

	return ReqBookList{
		Data: map[string]interface{}{
			"param": map[string]interface{}{
				"OwnershipData": OwnershipData,
			},
		},
	}
}
