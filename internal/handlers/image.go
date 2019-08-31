package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	_ "image/gif" //
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"satellity/internal/configs"
	"satellity/internal/session"
)

func UploadImage(ctx context.Context, name, data string) (string, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(data)
	cfg, fmt, err := image.DecodeConfig(bytes.NewReader(imageBytes))
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	if cfg.Width < 600 || cfg.Height < 300 {
		return "", session.InvalidImageDataError(ctx)
	}

	fileName := name + "." + fmt
	file := filepath.Join(configs.AppConfig.System.Attachments.Path, fileName)
	err = os.MkdirAll(filepath.Dir(file), os.ModePerm)
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	defer f.Close()
	_, err = f.Write(imageBytes)
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	err = f.Sync()
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	return fileName, nil
}
