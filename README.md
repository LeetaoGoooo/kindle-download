# KindleDownload

> 下载亚马逊上你的 **电子书** 以及 **个人文档**

# 使用说明

## 配置说明

修改 `config.yaml` 文件中的对应参数

```yaml
common:
  workerNum: 5 # 最大协程数，协程数并不是越大越好
  fileDir: ebooks  # 下载文档保存的目录
  cookie: "" # 登录后的 cookie，需要手动登录，在浏览器将 cookies 复制到这里
  cn: true # 是否是国区
# 一下配置不要修改
cn:
  listUrl: https://www.amazon.cn/hz/mycd/myx#/home/content/booksAll
  ajaxUrl: https://www.amazon.cn/hz/mycd/ajax
com:
  listUrl: https://www.amazon.com/hz/mycd/myx#/home/content/booksAll
  ajaxUrl: https://www.amazon.com/hz/mycd/ajax
```

### cookie 的获取

登录亚马逊后，F12 打开浏览器控制台，然后找到任意请求，将下图的 cookie 对应的值，复制到 `config.yaml`

![](https://raw.githubusercontent.com/LeetaoGoooo/leetaogoooo.github.io/images/%E6%88%AA%E5%B1%8F2022-06-09%2021.16.06.png)

**注意**：cookie 里存在 “"” 的情况，需要在 “"” 前加 "\" 转义，如下:

```
sess-at-main="4EHeJ63CJDL9LtXbGwFXCGOKyI0sooxh381f3GjizmM="; 
转义后： 
sess-at-main=\"4EHeJ63CJDL9LtXbGwFXCGOKyI0sooxh381f3GjizmM=\"; 
```



## 使用

### 二进制文件

可以使用已经编译好的二进制文件

```
# mac 输入以下命令，回车运行
./kindle-download
```

### 直接运行

```golang
go run kindle-download.go
```

### 编译运行

```golang
go build kindle-download.go
```
编译成功后，当前目录会生成一个 `kindle-download` 的可执行程序


# 感谢

[Kindle_download_helper-Python](https://github.com/yihong0618/Kindle_download_helper)