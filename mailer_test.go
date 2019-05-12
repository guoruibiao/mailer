package mailer

import (
	"testing"
	"fmt"
)

func TestMailer_Auth(t *testing.T) {
	mailer := NewMailer("smtp.163.com", 25)
	err := mailer.Auth("username", "password")
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("authen success.")
	}
}
