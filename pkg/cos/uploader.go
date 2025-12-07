package cos

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// UploadOptions 上传选项
type UploadOptions struct {
	AttachmentType string // voucher/invoice_file/red_invoice_file
	FileName       string // 原始文件名
	ContentType    string // MIME类型
}

// UploadResult 上传结果
type UploadResult struct {
	URL      string // 访问URL
	FilePath string // COS存储路径
	FileSize int64  // 文件大小
}

// Upload 上传文件到COS
func (c *COSClient) Upload(ctx context.Context, reader io.Reader, size int64, opts UploadOptions) (*UploadResult, error) {
	// 生成存储路径
	now := time.Now()
	uid := uuid.New().String()
	ext := filepath.Ext(opts.FileName)

	// 去除原文件名的扩展名
	nameWithoutExt := strings.TrimSuffix(opts.FileName, ext)
	filename := fmt.Sprintf("%s_%s%s", uid, nameWithoutExt, ext)

	cosPath := fmt.Sprintf("%s/%s/%d/%02d/%s",
		c.config.CompanyName,
		opts.AttachmentType,
		now.Year(),
		now.Month(),
		filename,
	)

	// 上传到COS
	_, err := c.client.Object.Put(ctx, cosPath, reader, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: opts.ContentType,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("上传失败: %w", err)
	}

	// 生成访问URL
	fileURL := fmt.Sprintf("%s/%s", strings.TrimRight(c.config.BaseURL, "/"), cosPath)

	return &UploadResult{
		URL:      fileURL,
		FilePath: cosPath,
		FileSize: size,
	}, nil
}

// Delete 删除COS文件
func (c *COSClient) Delete(ctx context.Context, filePath string) error {
	_, err := c.client.Object.Delete(ctx, filePath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// BatchDelete 批量删除文件
func (c *COSClient) BatchDelete(ctx context.Context, filePaths []string) error {
	if len(filePaths) == 0 {
		return nil
	}

	objects := make([]cos.Object, len(filePaths))
	for i, path := range filePaths {
		objects[i] = cos.Object{Key: path}
	}

	opt := &cos.ObjectDeleteMultiOptions{
		Objects: objects,
		Quiet:   true,
	}

	_, _, err := c.client.Object.DeleteMulti(ctx, opt)
	if err != nil {
		return fmt.Errorf("批量删除失败: %w", err)
	}

	return nil
}

// GetPresignedURL 生成预签名URL
func (c *COSClient) GetPresignedURL(ctx context.Context, filePath string, expire time.Duration) (string, error) {
	presignedURL, err := c.client.Object.GetPresignedURL(ctx, "GET", filePath, c.config.SecretID, c.config.SecretKey, expire, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}
	return presignedURL.String(), nil
}

// CopyObject 复制COS对象到新路径
func (c *COSClient) CopyObject(ctx context.Context, sourcePath, destPath string) error {
	// COS Copy API 的 source URL 格式: {bucket}.cos.{region}.myqcloud.com/{object-key}
	sourceURL := fmt.Sprintf("%s.cos.%s.myqcloud.com/%s",
		c.config.Bucket,
		c.config.Region,
		sourcePath,
	)

	_, _, err := c.client.Object.Copy(ctx, destPath, sourceURL, nil)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}
	return nil
}

// ListObjects 列出指定前缀的所有对象
func (c *COSClient) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	var allKeys []string
	marker := ""

	for {
		opt := &cos.BucketGetOptions{
			Prefix:  prefix,
			Marker:  marker,
			MaxKeys: 1000,
		}

		result, _, err := c.client.Bucket.Get(ctx, opt)
		if err != nil {
			return nil, fmt.Errorf("列出对象失败: %w", err)
		}

		for _, obj := range result.Contents {
			allKeys = append(allKeys, obj.Key)
		}

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}

	return allKeys, nil
}

// ObjectExists 检查对象是否存在
func (c *COSClient) ObjectExists(ctx context.Context, filePath string) (bool, error) {
	_, err := c.client.Object.Head(ctx, filePath, nil)
	if err != nil {
		if cos.IsNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
