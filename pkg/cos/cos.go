package cos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// COSClient COS客户端封装
type COSClient struct {
	client *cos.Client
	config *Config
}

var instance *COSClient

// InitCOSClient 初始化COS客户端(单例)
func InitCOSClient(cfg *Config) error {
	if !cfg.Enable {
		return fmt.Errorf("COS未启用")
	}

	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region))
	if err != nil {
		return fmt.Errorf("解析COS URL失败: %w", err)
	}

	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	instance = &COSClient{
		client: client,
		config: cfg,
	}

	return nil
}

// GetInstance 获取COS客户端实例
func GetInstance() (*COSClient, error) {
	if instance == nil {
		return nil, fmt.Errorf("COS客户端未初始化")
	}
	return instance, nil
}

// IsEnabled 检查COS是否已启用
func IsEnabled() bool {
	return instance != nil && instance.config.Enable
}

// HealthCheck 健康检查
func (c *COSClient) HealthCheck(ctx context.Context) error {
	_, _, err := c.client.Bucket.Get(ctx, &cos.BucketGetOptions{
		MaxKeys: 1,
	})
	return err
}
