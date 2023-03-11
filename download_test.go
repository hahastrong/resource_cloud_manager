package rcm

import (
	"fmt"
	"testing"
)

func TestFileDownload(t *testing.T) {
	d := SingleFileResource{
		Url: "https://hahastrong.com/weekvideo/audio/202303/Analysis-20230309.mp3",
		Path: "wv202303/",
		FileName: "Analysis-20230309.mp3",
	}

	d.Download()
}

func TestM3u8Download(t *testing.T) {

	d, err := ParseUrl("https://hahastrong.com/weekvideo/wv20230309/LT.m3u8")
	if err != nil {
		return
	}

	d.Download()

	DeleteExpiredResource()
}

func TestDownloadDir(t *testing.T) {
	err := DownloadDirFiles("http://127.0.0.1:8080/weekvideo/wv20220815/")
	if err != nil {
		fmt.Println(err)
	}
}

