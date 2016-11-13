package notifier

type PushbulletNotifier struct {
	ApiKey string
}

func (notifier *PushbulletNotifier) send(message string, recipients []string) {

}
