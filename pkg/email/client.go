package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

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

// SendTaskCreatedEmail 发送任务创建通知邮件
func (c *Client) SendTaskCreatedEmail(to, taskID string) error {
	// 检查邮件配置是否为空
	if c.host == "" || c.from == "" {
		fmt.Println("邮件配置为空，跳过发送邮件")
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.from, c.fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "下载任务已创建")

	// HTML 邮件正文
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>您的下载任务已提交</h2>
			<p>任务 ID: <strong>%s</strong></p>
			<p>您的任务已成功提交并进入队列，系统正在处理中。</p>
			<p>处理完成后，我们会再次发送邮件通知您。</p>
			<p>您可以使用任务 ID 在系统中查询任务状态。</p>
			<br>
			<p>感谢使用 Pathogen Query System!</p>
		</body>
		</html>
	`, taskID)

	m.SetBody("text/html", body)

	// 创建 SMTP 拨号器
	d := gomail.NewDialer(c.host, c.port, c.username, c.password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	fmt.Printf("✅ 任务创建通知邮件已发送至: %s\n", to)
	return nil
}

// SendTaskProcessingEmail 发送任务处理中通知邮件
func (c *Client) SendTaskProcessingEmail(to, taskID string) error {
	// 检查邮件配置是否为空
	if c.host == "" || c.from == "" {
		fmt.Println("邮件配置为空，跳过发送邮件")
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.from, c.fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "下载任务开始处理")

	// HTML 邮件正文
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>您的下载任务开始处理</h2>
			<p>任务 ID: <strong>%s</strong></p>
			<p>系统已开始处理您的下载任务，正在提取序列数据。</p>
			<p>处理完成后，我们会再次发送邮件通知您。</p>
			<p>您可以使用任务 ID 在系统中查询任务进度。</p>
			<br>
			<p>感谢使用 Pathogen Query System!</p>
		</body>
		</html>
	`, taskID)

	m.SetBody("text/html", body)

	// 创建 SMTP 拨号器
	d := gomail.NewDialer(c.host, c.port, c.username, c.password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	fmt.Printf("✅ 任务处理中通知邮件已发送至: %s\n", to)
	return nil
}

// SendTaskCompletedEmail 发送任务完成通知邮件
func (c *Client) SendTaskCompletedEmail(to, taskID, downloadURL string) error {
	// 检查邮件配置是否为空
	if c.host == "" || c.from == "" {
		fmt.Println("邮件配置为空，跳过发送邮件")
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.from, c.fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "下载任务完成通知")

	// HTML 邮件正文
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>您的下载任务已完成</h2>
			<p>任务 ID: <strong>%s</strong></p>
			<p>您可以通过以下方式获取文件:</p>
			<ul>
				<li>复制下载链接 %s/api/v1/tasks/%s/download 以下载文件</li>
				<li>查看任务详情: <a href="%s/tasks?query_type=task_id&task_ids=%s">任务详情</a></li>
			</ul>
			<p>您也可以复制任务 ID 在系统中查询下载。</p>
			<br>
			<p>感谢使用 Pathogen Query System!</p>
		</body>
		</html>
	`, taskID, downloadURL, taskID, downloadURL, taskID)

	m.SetBody("text/html", body)

	// 创建 SMTP 拨号器
	d := gomail.NewDialer(c.host, c.port, c.username, c.password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	fmt.Printf("✅ 邮件已发送至: %s\n", to)
	return nil
}

// SendTaskFailedEmail 发送任务失败通知邮件
func (c *Client) SendTaskFailedEmail(to, taskID, errorMsg, baseURL string) error {
	// 检查邮件配置是否为空
	if c.host == "" || c.from == "" {
		fmt.Println("邮件配置为空，跳过发送邮件")
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.from, c.fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "下载任务失败通知")

	// HTML 邮件正文
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>您的下载任务执行失败</h2>
			<p>任务 ID: <strong>%s</strong></p>
			<p>失败原因: <span style="color: red;">%s</span></p>
			<p>您可以 <a href="%s/tasks?query_type=task_id&task_ids=%s">查看任务详情</a> 或尝试重新提交任务。</p>
			<br>
			<p>Pathogen Query System</p>
		</body>
		</html>
	`, taskID, errorMsg, baseURL, taskID)

	m.SetBody("text/html", body)

	// 创建 SMTP 拨号器
	d := gomail.NewDialer(c.host, c.port, c.username, c.password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	fmt.Printf("✅ 邮件已发送至: %s\n", to)
	return nil
}
