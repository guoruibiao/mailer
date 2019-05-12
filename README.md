## mailer

yet, another mail helper in golang.

## dependency

> there are need no dependency expect the standard lib of golang.

## demonstration
```
import (
	mailer2 "github.com/guoruibiao/mailer"
	"fmt"
			"net/smtp"
		"strings"
	"io/ioutil"
	"encoding/base64"
	"net"
	"bytes"
)

func main() {
	mailer := mailer2.NewMailer("smtp.163.com", 25)
	err := mailer.Auth("mail-username", "authen-code")
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("authen result success ")
	}
	mailer.Send("from-mail-address", []string{"marksinoberg@gmail.com", "another-mail-address"}, []string{"cc-mail-address"}, "title of this mail", "subject of this mail", "content of this mail")
}

```

## todo
- [x] plain text mail
- [x] plain text file as the attachment
- [ ] attachment whose size > 50KB caused `421 Read data from client error`