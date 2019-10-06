package clouds

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif" //
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"satellity/internal/configs"
	"satellity/internal/session"
)

// UploadImage upload image to storage
func UploadImage(ctx context.Context, name, data string) (string, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(data)
	h := fmt.Sprintf("%x", md5.Sum(imageBytes))
	cfg, ext, err := image.DecodeConfig(bytes.NewReader(imageBytes))
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	if cfg.Width < 256 || cfg.Height < 256 {
		return "", session.InvalidImageDataError(ctx)
	}

	fileName := name + "." + ext
	file := filepath.Join(configs.AppConfig.System.Attachments.Path, fileName)
	err = os.MkdirAll(filepath.Dir(file), os.ModePerm)
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	defer f.Close()
	_, err = f.WriteAt(imageBytes, 0)
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	err = f.Sync()
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	return fmt.Sprintf("%s/attachments%s?v=%s", configs.AppConfig.HTTP.Host, fileName, h), nil
}
