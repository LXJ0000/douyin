package util

import (
	"bytes"
	"douyin/config"
	"douyin/models"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
	"strings"
)

func NewFileName(userid int64) string {
	var count int64
	err := models.NewVideoDAO().QueryVideoCountByUserId(userid, &count)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d-%d", userid, count)
}

// GetSnapshot
// 传参 视频地址 封面保存地址 获取第几帧 eg static\cat.mp4 static\cat 5
// 返回 封面名 eg cat.png
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("解码缩略图失败：", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		log.Fatal("保存缩略图失败：", err)
		return "", err
	}

	fmt.Println("--snapshotPath--", snapshotPath)
	// --snapshotPath-- ./static/testImage

	names := strings.Split(snapshotPath, `\`)
	fmt.Println("----names----", names)
	// ----names---- [./static/testImage]
	// 这里把 snapshotPath 的 string 类型转换成 []string

	snapshotName = names[len(names)-1] + ".png"
	fmt.Println("----snapshotName----", snapshotName)
	// ----snapshotName---- ./static/testImage.png

	return snapshotName, nil
}

func GetFileUrl(fileName string) string {
	base := fmt.Sprintf("%s/static/%s", config.Info.URL, fileName)
	return base
}
