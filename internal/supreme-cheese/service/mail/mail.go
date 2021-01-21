package mail

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types/config"
	"github.com/pkg/errors"
	"strconv"
)

type Request struct {
	from  string
	to    []string
	body  string
	email *config.ConfigForSendEmail
}

func NewRequest(to []string, cnf *config.ConfigForSendEmail) *Request {
	return &Request{to: to, email: cnf}
}

//func (r *Request) parseTemplate(fileName string, data interface{}) error {
//	t, err := template.ParseFiles(fileName)
//	if err != nil {
//		return errors.Wrap(err, "err while ParseFiles")
//	}
//	buffer := new(bytes.Buffer)
//	dataStr := struct {
//		Data string
//	}{
//		Data: data.(string),
//	}
//
//	if err = t.Execute(buffer, dataStr); err != nil {
//		return errors.Wrap(err, "err while Execute")
//	}
//	r.body = buffer.String()
//	return nil
//}

func (r *Request) sendMail(template string, data string) error {

	//body := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\r\n%s\r\n", r.email.EmailLogin, r.to[0], template, data)
	//addr := fmt.Sprintf("%s:%s", r.email.EmailHost, r.email.EmailPort)
	//if err := smtp.SendMail(addr, smtp.PlainAuth("", r.email.EmailLogin, r.email.EmailPass, r.email.EmailHost), r.email.EmailLogin, r.to, []byte(body)); err != nil {
	//	return errors.Wrap(err, "err while SendMail")
	//}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "support@supreme-cheese.ru", "support@supreme-cheese.ru")
	m.SetAddressHeader("To", r.to[0], r.to[0])
	m.SetHeader("From", "support@supreme-cheese.ru")
	m.SetHeader("To", r.to[0])
	m.SetHeader("Subject", template)
	m.SetHeader("MIME-Version:", "1.0")
	m.SetHeader("Reply-To", r.to[0])
	//	m.SetHeader()
	m.SetBody("text/plain", data)
	port, err := strconv.Atoi(r.email.EmailPort)
	if err != nil {
		return errors.Wrap(err, "err while Atoi ")
	}
	d := gomail.NewDialer(r.email.EmailHost, port, r.email.EmailLogin, r.email.EmailPass)

	if err := d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "err while DialAndSend ")
	}

	return nil
}

func (r *Request) Send(template string, data string) error {
	//err := r.parseTemplate(templateName, data)
	//if err != nil {
	//	return errors.Wrap(err, "err while parseTemplate")
	//}
	if err := r.sendMail(template, data); err != nil {
		return errors.Wrap(err, "err while sendMail")
	}

	return nil
}

func (r *Request) SendCert(template string, userID int, pathForCert string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "support@supreme-cheese.ru", "support@supreme-cheese.ru")
	m.SetAddressHeader("To", r.to[0], r.to[0])
	m.SetHeader("From", "support@supreme-cheese.ru")
	m.SetHeader("To", r.to[0])
	m.SetHeader("Subject", template)
	m.SetBody("text/plain", "Добрый день!\n\nВаш сертификат во вложении!")
	m.Attach(fmt.Sprintf("%s/%d/Certificate.pdf", pathForCert, userID))
	port, err := strconv.Atoi(r.email.EmailPort)
	if err != nil {
		return err
	}
	d := gomail.NewDialer(r.email.EmailHost, port, r.email.EmailLogin, r.email.EmailPass)

	if err := d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "err with d.DialAndSend in SendCert")
	}

	return nil
}

func (r *Request) SendReklama(template string, data string) error {
	fmt.Printf("%s, send\n", r.to[0])
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "support@supreme-cheese.ru", "support@supreme-cheese.ru")
	m.SetAddressHeader("To", r.to[0], r.to[0])
	m.SetHeader("From", "support@supreme-cheese.ru")
	m.SetHeader("To", r.to[0])
	m.SetHeader("Subject", template)
	m.SetHeader("MIME-Version:", "1.0")
	m.SetHeader("Reply-To", "support@supreme-cheese.ru")
	m.SetHeader("List-Unsubscribe", fmt.Sprintf("<mailto: unsubscribe-supreme@yandex.ru>, <https://supreme-cheese.ru:8080/mail/unsubscribe/%v>", r.to[0])) //r.to[0])) //, <https://supreme-cheese.ru:8080/mail/unsubscribe/%s>
	m.SetHeader("List-Unsubscribe-Post", "List-Unsubscribe=One-Click")

	m.SetBody("text/html", data)
	port, err := strconv.Atoi(r.email.EmailPort)
	if err != nil {
		return errors.Wrap(err, "err while Atoi ")
	}
	d := gomail.NewDialer(r.email.EmailHost, port, r.email.EmailLogin, r.email.EmailPass)

	if err := d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "err while DialAndSend ")
	}

	return nil
}
