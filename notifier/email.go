package notifier

type EmailNotifier struct {
	Host     string
	Port     int
	From     string
	Password string
}

func (notifier *EmailNotifier) send(message string, recipients []string) {

}
