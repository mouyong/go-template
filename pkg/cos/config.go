package cos

import (
	"fmt"

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

// LoadConfigFromDB 从数据库加载COS配置
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
