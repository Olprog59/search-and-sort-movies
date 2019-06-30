package myapp

import (
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

//func SendMail(subject, body string) {
//	from := "kameleon836@gmail.com"
//	pass := "jdsqinkdjuvaifoa"
//	to := GetEnv("email")
//	msg := "From: " + from + "\n" +
//		"To: " + to + "\n" +
//		"Bcc: " + from + "\n" +
//		"Subject: " + subject + "\n\n" +
//		body
//	err := smtp.SendMail("smtp.gmail.com:587",
//		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
//		from, []string{to}, []byte(msg))
//	if err != nil {
//		log.Printf("smtp error: %s", err)
//		return
//	}
//	log.Print("sent, visit http://foobarbazz.mailinator.com")
//}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// return "LOGIN", []byte{}, nil
	return "LOGIN", []byte(a.username), nil
}
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}

//func EnvoiDeMail(subject, fileName string) {
//	user := "admin@olprog.fr"
//	host := "127.0.0.1:1025"
//	to := GetEnv("mail")
//	body := `
//        <html>
//        <body>
//        <h3>
//        ` + fileName + `
//        </h3>
//        </body>
//        </html>
//        `
//	log.Println("send email")
//	err := SendMailLocal(host, user, subject, body, []string{to})
//	if err != nil {
//		log.Println("Send mail error!")
//		log.Println(err)
//	} else {
//		log.Println("Send mail success!")
//	}
//}

//func SendMailLocal(addr, from, subject, body string, to []string) error {
//	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")
//
//	c, err := smtp.Dial(addr)
//	if err != nil {
//		return err
//	}
//	defer c.Close()
//	if err = c.Mail(r.Replace(from)); err != nil {
//		return err
//	}
//	for i := range to {
//		to[i] = r.Replace(to[i])
//		if err = c.Rcpt(to[i]); err != nil {
//			return err
//		}
//	}
//
//	w, err := c.Data()
//	if err != nil {
//		return err
//	}
//
//	msg := "To: " + strings.Join(to, ",") + "\r\n" +
//		"From: " + from + "\r\n" +
//		"Subject: " + subject + "\r\n" +
//		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
//		"Content-Transfer-Encoding: base64\r\n" +
//		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
//
//	_, err = w.Write([]byte(msg))
//	if err != nil {
//		return err
//	}
//	err = w.Close()
//	if err != nil {
//		return err
//	}
//	return c.Quit()
//}
