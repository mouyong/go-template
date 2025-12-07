package email

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	Enable   bool   `json:"enable"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	FromName string `json:"from_name"`
	SSL      bool   `json:"ssl"`
}

// Client 邮件客户端
type Client struct {
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

// NewClient 创建邮件客户端
func NewClient(host string, port int, username, password, from, fromName string) *Client {
	return &Client{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		fromName: fromName,
	}
}

// NewClientFromConfig 从配置对象创建邮件客户端
func NewClientFromConfig(config *EmailConfig) *Client {
	if config == nil {
		return nil
	}

	// 转换端口号
	port := 465 // 默认 SSL 端口
	if config.Port != "" {
		if p, err := strconv.Atoi(config.Port); err == nil {
			port = p
		}
	}

	return &Client{
		host:     config.Host,
		port:     port,
		username: config.Username,
		password: config.Password,
		from:     config.From,
		fromName: config.FromName,
	}
}

// SendHTML 发送 HTML 格式邮件
func (c *Client) SendHTML(to, subject, htmlBody string) error {
	// 检查邮件配置是否为空
	if c == nil || c.host == "" || c.from == "" {
		fmt.Println("邮件配置为空，跳过发送邮件")
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.from, c.fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// 创建 SMTP 拨号器
	d := gomail.NewDialer(c.host, c.port, c.username, c.password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	fmt.Printf("✅ 邮件已发送至: %s (主题: %s)\n", to, subject)
	return nil
}

// downloadFile 下载远程文件到临时目录
func downloadFile(url string) (string, error) {
	// 创建临时文件
	tmpDir := os.TempDir()
	fileName := filepath.Base(url)
	// 如果 URL 中包含查询参数,移除它们
	if idx := strings.Index(fileName, "?"); idx != -1 {
		fileName = fileName[:idx]
	}
	tmpFile := filepath.Join(tmpDir, fileName)

	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载文件失败: HTTP %d", resp.StatusCode)
	}

	// 创建本地文件
	out, err := os.Create(tmpFile)
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	return tmpFile, nil
}

// SendHTMLWithAttachments 发送带附件的 HTML 邮件
func (c *Client) SendHTMLWithAttachments(to, subject, htmlBody string, attachments []string) error {
	// 检查邮件配置是否为空
	if c == nil || c.host == "" || c.from == "" {
		fmt.Println("邮件配置为空，跳过发送邮件")
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.from, c.fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// 下载并添加附件
	var tmpFiles []string
	for _, filePath := range attachments {
		if filePath == "" {
			continue
		}

		// 判断是本地文件还是远程 URL
		var localPath string
		if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
			// 远程文件,需要先下载
			tmpPath, err := downloadFile(filePath)
			if err != nil {
				fmt.Printf("⚠️  下载附件失败: %s, 错误: %v\n", filePath, err)
				continue
			}
			localPath = tmpPath
			tmpFiles = append(tmpFiles, tmpPath)
		} else {
			// 本地文件
			localPath = filePath
		}

		m.Attach(localPath)
	}

	// 发送邮件后清理临时文件
	defer func() {
		for _, tmpFile := range tmpFiles {
			os.Remove(tmpFile)
		}
	}()

	// 创建 SMTP 拨号器
	d := gomail.NewDialer(c.host, c.port, c.username, c.password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	fmt.Printf("✅ 邮件已发送至: %s (主题: %s, 附件数: %d)\n", to, subject, len(attachments))
	return nil
}
