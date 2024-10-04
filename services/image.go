package services

import (
	"fmt"
	"github.com/562589540/Go-Midjourney/gclient"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadImage 下载图片并保存到本地
func DownloadImage(url, filepath string) error {

	client := gclient.GetGclient().GetHTTPClient()

	// 发送HTTP请求获取图片
	resp, err := client.Get(url)

	if err != nil {
		return fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: received status code %d", resp.StatusCode)
	}

	// 创建本地文件
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// 将图片数据写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	fmt.Println("Image successfully downloaded to", filepath)
	return nil
}

// SplitImage 将四宫格图片切割成四张单独的图片
func SplitImage(inputFile string, outputDir string) error {
	// 打开四宫格图片
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// 获取图片的宽度和高度
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// 计算每张子图片的宽度和高度
	subImgWidth := imgWidth / 2
	subImgHeight := imgHeight / 2

	// 定义四张图片的位置（左上、右上、左下、右下）
	positions := []image.Rectangle{
		image.Rect(0, 0, subImgWidth, subImgHeight),                // 左上
		image.Rect(subImgWidth, 0, imgWidth, subImgHeight),         // 右上
		image.Rect(0, subImgHeight, subImgWidth, imgHeight),        // 左下
		image.Rect(subImgWidth, subImgHeight, imgWidth, imgHeight), // 右下
	}

	// 遍历四个位置并将每张图片保存
	for i, pos := range positions {
		// 裁剪出子图
		subImg := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(pos)

		// 创建输出文件
		outputFile := filepath.Join(outputDir, fmt.Sprintf("subimage_%d.jpg", i+1))
		out, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer out.Close()

		// 保存子图片为JPEG格式
		err = jpeg.Encode(out, subImg, nil)
		if err != nil {
			return fmt.Errorf("failed to encode subimage: %v", err)
		}
	}

	fmt.Println("Image successfully split into four parts")
	return nil
}
