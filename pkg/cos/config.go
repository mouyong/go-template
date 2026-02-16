package cos

import (
	"fmt"
	"strings"

	"app/internal/models"

	"gorm.io/gorm"
)

// Config COS配置
type Config struct {
	Enable      bool   // 是否启用COS
	SecretID    string // 腾讯云SecretID
	SecretKey   string // 腾讯云SecretKey
	Bucket      string // 存储桶名称
	Region      string // 地域(如ap-guangzhou)
	CompanyName string // 公司标识
	BaseURL     string // 访问域名(可选CDN域名)
}

// LoadConfigFromDB 从数据库加载COS配置（旧系统 - system_settings）
func LoadConfigFromDB(db *gorm.DB, tenantID uint) (*Config, error) {
	// 查询COS配置
	var settings []models.SystemSetting
	err := db.Where("tenant_id = ? AND category = ?", tenantID, models.SettingCategoryCOS).
		Find(&settings).Error
	if err != nil {
		return nil, fmt.Errorf("查询COS配置失败: %w", err)
	}

	// 查询公司名称
	var tenant models.Tenant
	err = db.First(&tenant, tenantID).Error
	if err != nil {
		return nil, fmt.Errorf("查询租户信息失败: %w", err)
	}

	cfg := &Config{
		CompanyName: tenant.CompanyName,
	}

	// 解析配置
	for _, setting := range settings {
		switch setting.Key {
		case models.SettingKeyCOSEnable:
			cfg.Enable = setting.Value == "true"
		case models.SettingKeyCOSSecretID:
			cfg.SecretID = setting.Value
		case models.SettingKeyCOSSecretKey:
			cfg.SecretKey = setting.Value
		case models.SettingKeyCOSBucket:
			cfg.Bucket = setting.Value
		case models.SettingKeyCOSRegion:
			cfg.Region = setting.Value
		case models.SettingKeyCOSBaseURL:
			cfg.BaseURL = setting.Value
		}
	}

	// 验证必填配置
	if cfg.Enable {
		if cfg.SecretID == "" || cfg.SecretKey == "" || cfg.Bucket == "" || cfg.Region == "" {
			return nil, fmt.Errorf("COS配置不完整")
		}
	}

	return cfg, nil
}

// LoadConfigByScene 根据场景从数据库加载COS配置（新系统 - cos_configs）
func LoadConfigByScene(db *gorm.DB, scene string) (*Config, error) {
	var cosConfig models.CosConfig
	err := db.Where("scene = ?", scene).First(&cosConfig).Error
	if err != nil {
		return nil, fmt.Errorf("查询场景 %s 的 COS 配置失败: %w", scene, err)
	}

	cfg := &Config{
		Enable:    true, // cos_configs 表中的配置默认启用
		SecretID:  cosConfig.SecretID,
		SecretKey: cosConfig.SecretKey,
		Bucket:    cosConfig.Bucket,
		Region:    cosConfig.Region,
		BaseURL:   cosConfig.Domain,
	}

	// 验证必填配置
	if cfg.SecretID == "" || cfg.SecretKey == "" || cfg.Bucket == "" || cfg.Region == "" {
		return nil, fmt.Errorf("场景 %s 的 COS 配置不完整", scene)
	}

	return cfg, nil
}

// GetURLFromPath 根据文件路径构造完整的访问 URL
// scene: 场景名称（如 "template", "attachment"）
// path: 文件相对路径（如 "template/2026-02/xxx.xlsx"）
func GetURLFromPath(db *gorm.DB, scene string, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	// 如果已经是完整 URL，直接返回
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}

	// 尝试加载场景配置
	config, err := LoadConfigByScene(db, scene)
	if err != nil {
		// 如果场景配置不存在，返回相对路径
		return path, fmt.Errorf("无法加载场景 %s 的配置: %w", scene, err)
	}

	// 构造完整 URL
	if config.BaseURL != "" {
		// 使用自定义域名（CDN）
		return fmt.Sprintf("%s/%s", strings.TrimRight(config.BaseURL, "/"), path), nil
	}

	// 使用默认的 COS 域名
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s",
		config.Bucket,
		config.Region,
		path,
	), nil
}
