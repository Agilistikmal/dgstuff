package mail

type Mail struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       string
	Subject  string
	Body     string
	Template string
	Data     map[string]interface{}
}

type Mailer interface {
	Send() error
}

func NewMail(smtp bool, m *Mail) Mailer {
	if smtp {
		return NewSMTP(m)
	}
	return NewSendMail(m)
}
