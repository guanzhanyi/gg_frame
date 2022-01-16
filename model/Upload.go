package model

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
	"goblog/settings"
	"goblog/utils/errmsg"
	"mime/multipart"
)

var (
	secretKey string
	accessKey string
	bucket    string
	imgUrl    string
)

func InitOSS(cfg *settings.OSSConfig) {
	secretKey = cfg.SecretKey
	accessKey = cfg.AccessKey
	bucket = cfg.Bucket
	imgUrl = cfg.ServerAddr
}

func UpLoadFile(file multipart.File, fileSize int64) (string, int) {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	zap.L().Info(accessKey)
	zap.L().Info(secretKey)
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	zap.L().Info(upToken)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		zap.L().Error("upload failed", zap.Error(err))
		return "", errmsg.ERROR
	}
	url := imgUrl + ret.Key
	return url, errmsg.SUCCSE

}
