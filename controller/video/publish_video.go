package video

import (
	"douyin/config"
	"douyin/models"
	"douyin/service/video"
	"douyin/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

// key value 判断是否是视频/图片
var (
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
	//pictureIndexMap = map[string]struct{}{
	//	".jpg": {},
	//	".bmp": {},
	//	".png": {},
	//	".svg": {},
	//}
)

func PublishVideoHandler(c *gin.Context) {
	//准备参数
	rawId, _ := c.Get("user_id")
	//判断是否是int64
	userId, ok := rawId.(int64)
	if !ok {
		c.JSON(http.StatusOK, models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  "解析UserId出错",
		})
		return
	}

	//form-data里拿出数据来
	title := c.PostForm("title")
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	files := form.File["data"]
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if _, ok := videoIndexMap[ext]; !ok {
			c.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "不支持该视频格式",
			})
			continue
		}

		//  生成唯一的文件名用于保存
		fileName := util.NewFileName(userId)
		fullName := fileName + ext

		//	写入static
		savePath := filepath.Join(config.Info.StaticSourcePath, fullName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			continue
		}
		//  获取视频封面并写入static
		snapshotPath := filepath.Join(config.Info.StaticSourcePath, fileName)
		coverName, err := util.GetSnapshot(savePath, snapshotPath, 5)
		if err != nil {
			c.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			continue
		}
		//  数据持久化
		if err := video.PostVideo(userId, fullName, coverName, title); err != nil {
			c.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			continue
		}
		c.JSON(http.StatusOK, models.CommonResponse{
			StatusCode: 0,
			StatusMsg:  file.Filename + "上传成功",
		})
	}
}
