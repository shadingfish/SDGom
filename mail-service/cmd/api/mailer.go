package main

// 整个邮件发送服务通过SMTP协议发送HTML和纯文本格式的邮件，并支持自定义模板和附件。

// The entire mail sending service sends HTML and plain text formatted emails via SMTP protocol, and supports custom templates and attachments.

// 该实现支持不同的加密方式（TLS、SSL、无加密），并能够将CSS样式转换为内联样式以确保兼容性。

// The implementation supports different encryption methods (TLS, SSL, no encryption), and can convert CSS styles to inline styles to ensure compatibility.

// 通过 getEncryption 函数，可以将字符串加密类型转换为具体的 mail.Encryption 类型, 以便在SMTP配置中使用.

// The getEncryption function can convert string encryption types to specific mail.Encryption types for use in SMTP configuration.

import (
	"bytes"         // 引入 bytes 包，用于操作字节缓冲区
	"html/template" // 引入 html/template 包，用于解析和渲染HTML模板
	"log"           // 引入 log 包，用于记录日志信息
	"time"          // 引入 time 包，用于时间操作

	"github.com/vanng822/go-premailer/premailer" // 引入 premailer 包，用于将CSS样式转换为内联样式
	mail "github.com/xhit/go-simple-mail/v2"     // 引入 go-simple-mail/v2 包，用于发送邮件
)

// 定义了一个 Mail 结构体，用于存储邮件相关的配置参数
// Defines a Mail struct to store configuration parameters related to email.
type Mail struct {
	Domain      string // 邮件域名，如 example.com
	Host        string // 邮件服务器的主机名或IP地址
	Port        int    // 邮件服务器的端口号
	Username    string // 邮件服务器登录的用户名
	Password    string // 邮件服务器登录的密码
	Encryption  string // 邮件加密类型，如TLS或SSL
	FromAddress string // 发送者的邮箱地址
	FromName    string // 发送者的名称
}

// 定义了一个 Message 结构体，用于封装邮件的具体内容
// Defines a Message struct to encapsulate the specific content of an email.
type Message struct {
	From        string   // 邮件发送者的地址（可选，如果为空则使用默认地址）
	FromName    string   // 邮件发送者的名称（可选，如果为空则使用默认名称）
	To          string   // 邮件接收者的地址
	Subject     string   // 邮件主题
	Attachments []string // 邮件附件列表
	// []string 是 Go 语言中的字符串切片（slice of strings），表示附件的文件路径列表。
	// [] 表示切片类型，string 表示切片中每个元素都是字符串类型。
	// 切片（slice）是一种动态数组类型，它的长度可以改变，并且提供了更为灵活的元素操作方式。

	Data any // 邮件正文的数据，可以是任意类型
	// any 是 Go 1.18 及以上版本中引入的新类型，表示任意类型的数据。它是 interface{} 的别名，可以用来存储任何类型的值。
	// any 类型的字段可以接收和存储任意数据类型的值，如字符串、整数、布尔值、结构体、切片、映射等。

	DataMap map[string]any // 数据映射，用于传递到模板中
	// map[string]any 是 Go 中的映射类型（map），表示键值对的集合，其中键是 string 类型，值是 any 类型。
	// map 是 Go 中的一种内置数据结构，类似于其他语言中的字典（dictionary）或哈希表（hash table）。
	// any 表示映射的值可以是任意类型。
}

// SendSMTPMessage 方法用于通过SMTP协议发送邮件
// The SendSMTPMessage method is used to send an email via SMTP protocol.
func (m *Mail) SendSMTPMessage(msg Message) error {
	// 如果消息中的发件人地址为空，则使用Mail结构体中配置的默认发件人地址
	// If the sender's address in the message is empty, use the default sender's address in the Mail struct.
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	// 如果消息中的发件人名称为空，则使用Mail结构体中配置的默认发件人名称
	// If the sender's name in the message is empty, use the default sender's name in the Mail struct.
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// 创建一个数据映射，将消息内容存入DataMap字段中
	// Create a data map and store the message content in the DataMap field.
	data := map[string]any{
		"message": msg.Data,
	}

	// 将数据映射赋值给消息的 DataMap 字段
	// Assign the data map to the message's DataMap field.
	msg.DataMap = data

	// 构建 HTML 格式的邮件内容
	// Build the HTML formatted email content.
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	// 构建纯文本格式的邮件内容
	// Build the plain text formatted email content.
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// 创建一个新的 SMTP 服务器客户端，并配置其属性
	// Create a new SMTP server client and configure its properties.
	server := mail.NewSMTPClient()
	server.Host = m.Host                              // 设置邮件服务器主机地址
	server.Port = m.Port                              // 设置邮件服务器端口号
	server.Username = m.Username                      // 设置邮件服务器登录用户名
	server.Password = m.Password                      // 设置邮件服务器登录密码
	server.Encryption = m.getEncryption(m.Encryption) // 设置加密方式
	server.KeepAlive = false                          // 不保持连接
	server.ConnectTimeout = 10 * time.Second          // 设置连接超时时间为10秒
	server.SendTimeout = 10 * time.Second             // 设置发送邮件超时时间为10秒

	// 尝试连接到SMTP服务器，并获取SMTP客户端
	// Attempt to connect to the SMTP server and get the SMTP client.
	smtpClient, err := server.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	// 创建一个新的邮件消息，并设置相关的邮件属性
	// Create a new email message and set related email properties.
	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	// 设置邮件的主体内容，包括纯文本格式和HTML格式
	// Set the body content of the email, including plain text and HTML formats.
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	// 如果消息中包含附件，则将附件添加到邮件中
	// If the message contains attachments, add them to the email.
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// 发送邮件
	// Send the email.
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	// 返回nil表示邮件发送成功
	// Return nil to indicate that the email was sent successfully.
	return nil
}

// buildHTMLMessage 构建 HTML 格式的邮件内容
// buildHTMLMessage builds the HTML formatted email content.
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	// 定义要渲染的HTML模板路径
	// Define the path to the HTML template to render.
	templateToRender := "./templates/mail.html.gohtml"

	// 解析指定的HTML模板文件
	// Parse the specified HTML template file.
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// 使用字节缓冲区存储渲染后的HTML内容
	// Use a byte buffer to store the rendered HTML content.
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// 获取渲染后的HTML字符串，并将CSS样式转换为内联样式
	// Get the rendered HTML string and convert the CSS styles to inline styles.
	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	// 返回最终的HTML格式邮件内容
	// Return the final HTML formatted email content.
	return formattedMessage, nil
}

// buildPlainTextMessage 构建纯文本格式的邮件内容
// buildPlainTextMessage builds the plain text formatted email content.
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	// 定义要渲染的纯文本模板路径
	// Define the path to the plain text template to render.
	templateToRender := "./templates/mail.plain.gohtml"

	// 解析指定的纯文本模板文件
	// Parse the specified plain text template file.
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// 使用字节缓冲区存储渲染后的纯文本内容
	// Use a byte buffer to store the rendered plain text content.
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// 获取渲染后的纯文本字符串
	// Get the rendered plain text string.
	plainMessage := tpl.String()

	// 返回最终的纯文本格式邮件内容
	// Return the final plain text formatted email content.
	return plainMessage, nil
}

// inlineCSS 将HTML内容中的CSS样式转换为内联样式
// inlineCSS converts CSS styles in the HTML content to inline styles.
func (m *Mail) inlineCSS(s string) (string, error) {
	// 配置 premailer 的选项
	// Configure the options for premailer.
	options := premailer.Options{
		RemoveClasses:     false, // 是否移除HTML中的类选择器
		CssToAttributes:   false, // 是否将CSS属性转换为HTML属性
		KeepBangImportant: true,  // 是否保留 !important 声明
	}

	// 使用传入的HTML内容创建一个新的 Premailer 实例
	// Create a new Premailer instance using the provided HTML content.
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	// 将HTML内容转换为内联样式
	// Transform the HTML content to inline styles.
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	// 返回转换后的HTML内容
	// Return the transformed HTML content.
	return html, nil
}

// getEncryption 根据传入的字符串返回对应的加密方式
// getEncryption returns the corresponding encryption type based on the input string.
func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		// 如果传入的字符串是 "tls"，则返回 STARTTLS 加密类型
		// If the input string is "tls", return the STARTTLS encryption type.
		return mail.EncryptionSTARTTLS
	case "ssl":
		// 如果传入的字符串是 "ssl"，则返回 SSL/TLS 加密类型
		// If the input string is "ssl", return the SSL/TLS encryption type.
		return mail.EncryptionSSLTLS
	case "none", "":
		// 如果传入的字符串是 "none" 或为空，则不使用任何加密
		// If the input string is "none" or empty, use no encryption.
		return mail.EncryptionNone
	default:
		// 默认返回 STARTTLS 加密类型
		// By default, return the STARTTLS encryption type.
		return mail.EncryptionSTARTTLS
	}
}
