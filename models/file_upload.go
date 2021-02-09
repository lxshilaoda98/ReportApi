package models

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	helper "github.com/n1n1n1_owner/ReportApi/models/Helper"
)

type ResultClass struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}


func (r *ResultClass) Up(file multipart.File, header *multipart.FileHeader) {
	config:= helper.GetIVRConfig()
	var FsDir = config.GetString("IVRConfig.dir")
	filename := header.Filename
	out, err := os.Create(FsDir+"/sounds/en/us/callie/IVRWav/" + filename)
	if err == nil {
		defer out.Close()
		_, err = io.Copy(out, file)
		if err == nil {
			res := map[string]interface{}{
				"filePath": "./static/res/uploadFile/ivrWav/" + filename,
				"fileName": filename,
			}
			r.Code = 0
			r.Data = res
			r.Msg = "上传成功"
		} else {
			r.Code = -3
			r.Msg = "复制文件出错"
		}
	} else {
		fmt.Println(err)
		r.Code = -2
		r.Msg = "创建文件出错"
	}
}
