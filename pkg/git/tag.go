package git

import (
	"time"
)

type TagObj struct {
	Name    string
	Message string
}

func (t *TagObj) SetMessage() {
	var msg string
	msg += DescribeLong()
	msg += "-" + time.Now().UTC().Format("20060102150405")
	t.Message = msg
}

func (t *TagObj) Push() {
	if t.Message == "" {
		t.SetMessage()
	}
	Tag(t.Name, t.Message)
	Push(t.Name)
}
