package rcm

import (
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type Downloader interface {
	Download() error
}

type M3u8Downloader struct {
	Url      string `json:"url"`
	BaseUrl  string
	Path     string
	FileName string
}


func ParseUrl(rawURL string) (d Downloader, err error) {
	// 解析 URL
	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// 判断 URL 格式是否正确
	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return nil, errors.New("Invalid URL: " + rawURL)
	}

	idx := len(rawURL) - 1
	for rawURL[idx] != '/' {
		idx--
	}

	secondIdx := idx-1
	for rawURL[secondIdx] != '/' {
		secondIdx--
	}

	filename := rawURL[idx+1:]

	if strings.Contains(rawURL, ".m3u8") {
		d = M3u8Downloader{
			Url: rawURL,
			FileName: filename,
			BaseUrl: rawURL[:idx+1],
			Path: rawURL[secondIdx+1:idx+1],
		}

		return d, nil
	}

	d = SingleFileResource{
		Url: rawURL,
		FileName: filename,
		Path: rawURL[secondIdx+1:idx+1],
	}
	return d, nil
}

func (d M3u8Downloader) Download() error {

	if err := CreateDir(d.Path); err != nil {
		return err
	}

	res, err := http.Get(d.Url)
	if err != nil {
		return err
	}

	// 解析 m3u8 文件
	playlist, listType, err := m3u8.DecodeFrom(res.Body, true)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}


	if listType != m3u8.MEDIA {
		return errors.New("mismatch type")
	}

	mediaList := playlist.(*m3u8.MediaPlaylist)

	// 遍历 m3u8 中的分段文件列表
	for _, segment := range mediaList.Segments {
		if segment == nil {
			break
		}
		// 获取分段文件 URL
		segmentUrl := d.BaseUrl + segment.URI

		// 发送 HTTP GET 请求，获取分段文件内容
		resp, err := http.Get(segmentUrl)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// 读取分段文件内容到内存中
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// 生成分段文件的保存路径
		filename := fmt.Sprintf("%s%03d.ts", d.Path, segment.SeqId)

		// 保存分段文件到本地文件系统中
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Downloaded", filename)
	}

	return nil
}

type SingleFileResource struct {
	Url      string `json:"url"`
	Path     string
	FileName string
}

func (d SingleFileResource) Download() error {

	if err := CreateDir(d.Path); err != nil {
		return err
	}


	res, err := http.Get(d.Url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return err
	}

	err = os.WriteFile(d.Path+d.FileName, body, 0644)

	fmt.Println("download finished!")

	return err
}


func CreateDir(path string) error {
	// 判断文件夹是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

const TimeLayout = "20060102"

func DeleteExpiredResource() {
	date := time.Now().AddDate(0,0,-15).Format(TimeLayout)
	dir := fmt.Sprintf("wv%s", date)

	// 删除文件夹
	_ = os.RemoveAll(dir)
}

func ParseDir(path string) ([]string, error) {
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return nil, err
	}

	r, _ := regexp.Compile("href=\"(.*?)\"")
	files := r.FindAllStringSubmatch(string(body),-1)

	var fileList []string
	for _, file := range files {
		if len(file) > 1 {
			fileList = append(fileList, file[1])
		}
	}

	return fileList, nil
}

func DownloadDirFiles(path string) (err error) {
	if path[len(path)-1] != '/' {
		path += "/"
	}

	files, err := ParseDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file[len(file)-1] == '/' {
			// actually doesn't support the branch
			err = DownloadDirFiles(path + file)
			if err != nil {
				return err
			}
			continue
		}

		// download single file
		idx := len(path) - 2
		for path[idx] != '/' {
			idx--
		}

		d := SingleFileResource{
			Url: path + file,
			FileName: file,
			Path: path[idx+1:],
		}

		err = d.Download()
		if err != nil {
			return err
		}
	}

	return
}


