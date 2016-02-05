package mail

import (
	"strconv"

	"./../pluginbase"
	"github.com/robertkrimen/otto"
	"gopkg.in/gomail.v1"
)

type SMTPSettings struct {
	Username   string
	Password   string
	Servername string
	Port       int
}

var config SMTPSettings

func InitPlugin() *pluginbase.Plugin {

	p1 := pluginbase.Plugin{
		Name: "mail",
		Init: func(vm *otto.Otto) {
			vm.Set("loadMailSettings", func(c otto.FunctionCall) otto.Value {
				loadSettings(vm)
				return otto.TrueValue()
			})
			vm.Set("send", func(c otto.FunctionCall) otto.Value {
				loadSettings(vm)
				recv, _ := c.Argument(0).ToString()
				subject, _ := c.Argument(1).ToString()
				msg, _ := c.Argument(2).ToString()

				err := sendmail(recv, subject, msg, "")

				return pluginbase.ToResult(vm, true, err)
			})
		},
	}

	return &p1
}

func loadSettings(vm *otto.Otto) {
	obj, _ := vm.Get("settings")
	plugins, _ := obj.Object().Get("Plugins")

	mailSettings, _ := plugins.Object().Get("mail")
	mailObj := mailSettings.Object()

	unameO, _ := mailObj.Get("username")
	pwO, _ := mailObj.Get("password")
	snameO, _ := mailObj.Get("servername")
	portO, _ := mailObj.Get("port")

	uname, _ := unameO.ToString()
	pw, _ := pwO.ToString()
	sname, _ := snameO.ToString()
	portStr, _ := portO.ToString()

	config = SMTPSettings{
		Username:   uname,
		Password:   pw,
		Servername: sname,
	}

	port, _ := strconv.Atoi(portStr)
	config.Port = port
}

func sendmail(email string, subject string, messageString string, txtAttachment string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Username)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", messageString)
	if len(txtAttachment) > 0 {
		f := gomail.CreateFile("attached.txt", []byte(txtAttachment))
		msg.Attach(f)
	}

	mailer := gomail.NewMailer(config.Servername, config.Username, config.Password, config.Port)
	err := mailer.Send(msg)
	return err
}