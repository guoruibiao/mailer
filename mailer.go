package mailer

import (
	"strconv"
	"net"
		"github.com/pkg/errors"
	"net/smtp"
	"encoding/base64"
	"bytes"
	"fmt"
	"strings"
	"io/ioutil"
	)

// reference from http://www.361way.com/golang-email-attachment/5856.html

// ServerConfig: config your smtp mail server
type ServerConfig struct {
	host string
	port int
}


// Mailer play as the master of email dispatcher
type Mailer struct{
	serverConfig ServerConfig
	auth smtp.Auth
}


// NewMailer return the instance by configing your SMTP server
func NewMailer(host string, port int) *Mailer {
	return &Mailer{
		ServerConfig{
			host:host,
			port:port,
		},
		nil,
	}
}

// Auth just only send the credentials if the connection is using TLS
// or is connected to localhost. Otherwise authentication will fail with an
// error, without sending the credentials.
// for more details which defined in RFC 4616.
func (mailer *Mailer)Auth(username, authencode string) error {
	address := mailer.serverConfig.host + ":" + strconv.Itoa(mailer.serverConfig.port)
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	if _, err = strconv.Atoi(port); err != nil {
		return errors.New("mial server port is not a valid number, error:" + err.Error())
	}
	auth := smtp.PlainAuth("", username, authencode, host)
	mailer.auth = auth
	return nil
}



// Send send your plain text email or simple attachment contains plain text
// todo handle the attachment whose's size > 50KB, which will cause `421 Read data from client error`
func (mailer *Mailer)Send(from string, to, cc []string, title, subject, body string) (bool, error) {
	// 设置邮件header啥的
	message := NewMailMessage("MY_BOUNDARY_DELIMITER")
	message.addHeader(from, to, cc, title, subject)
	// 添加邮件正文
	message.addContent(body)
	//message.addAttach("/Users/biao/go/src/github.com/guoruibiao/commands/commands.go", "text/plain")
	//message.addAttach("/Users/biao/Desktop/d6154ea0a1cfdb98e9bafc03bd31ffcb.png", "image/png")
	//message.addAttach("/Users/biao/Desktop/example.tar.gz.png", "application/octet-stream")

	serverName := mailer.serverConfig.host + ":" + strconv.Itoa(mailer.serverConfig.port)
	err := smtp.SendMail(serverName, mailer.auth, from, to, message.outlet())
	if err != nil {
		fmt.Println(err)
	}
	return true, nil
}

// MIMEMessage MIME，Multipurpose Internet Mail Extensions
// which extends the standard internet protocol, so we can free to send multi data
type MIMEMessage struct {
	boundary string
	container *bytes.Buffer
}


func NewMailMessage(boundary string) *MIMEMessage{
	return &MIMEMessage{
		boundary:boundary,
		container: bytes.NewBuffer(nil),
	}
}

func (mime *MIMEMessage)addHeader(from string, to []string, cc []string, title, subject string) {
	tolist := strings.Join(to, ",")
	cclist := strings.Join(cc, ",")
	header := fmt.Sprintf("From: %s<%s>\r\nTo: %s\r\nCC: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\n", from, from, tolist, cclist, subject)
	mime.container.WriteString(header)
	mime.container.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", mime.boundary))
	mime.container.WriteString(fmt.Sprintf("Content-Description: %s\r\n", title))
}

// addContent append the content to your email
func (mime *MIMEMessage)addContent(content string) {
	mime.container.WriteString(fmt.Sprintf("--%s\r\n", mime.boundary))
	mime.container.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	mime.container.WriteString(content)
}

// addAttach TODO: fix `421 Read data from client error` which caused by whose size > 50KB
func (mime *MIMEMessage)addAttach(filepath string, mimetype string) error{
	pathParams := strings.Split(filepath, "/")
	if len(pathParams) <= 0 {
		return errors.New("No such an attach file in path: " + filepath)
	}
	filename := pathParams[len(pathParams)-1]
	// 添加 附件栏内容相关元数据
	mime.container.WriteString(fmt.Sprintf("\r\n--%s\r\n", mime.boundary))
	//mime.container.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimetype))
	mime.container.WriteString("image/png\r\n")
	mime.container.WriteString(fmt.Sprintf("Content-Description: %s\r\n", filename))
	mime.container.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filename))
	// 读取文件并进行编码处理
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	if mimetype != "text/plain" {
		mime.container.WriteString(fmt.Sprintf("Content-Transfer-Encoding: base64\r\n"))
		edata := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(edata, data)
		// 切一下，复用下变量
		data = edata
	}
	mime.container.Write(data)
	fmt.Println(string(data))
	// 防止添加多个附件时缺少额外的分隔
	mime.container.WriteString("\r\n")
	return nil
}

// outlet take all sections of this mail into bytes for sending.
func (mime *MIMEMessage) outlet() []byte {
	if mime.container != nil && mime.container.Len()>0 {
		// 邮件正式结束
		mime.container.WriteString(fmt.Sprintf("\r\n--%s--\r\n\r\n", mime.boundary))
	}else {
		mime.container.WriteString(fmt.Sprintf("\r\n--%s--\r\n\r\n", mime.boundary))
	}
	return mime.container.Bytes()
}