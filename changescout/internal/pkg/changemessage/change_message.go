package changemessage

type ChangeMessage struct {
	Title string
	Body  string
}

func NewChangeMessage(title, body string, vars map[string]string) ChangeMessage {
	return ChangeMessage{
		Title: title,
		Body:  body,
	}
}

func (c ChangeMessage) String() string {
	return c.Title + "\n\n" + c.Body
}
