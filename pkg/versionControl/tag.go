package versionControl

import (
	"time"

	"semtag/pkg/output"
)

type TagObj struct {
	Name    string
	Message string
}

func (t *TagObj) GenerateMessage() {
	var msg string
	msg += DescribeLong()
	msg += "-" + time.Now().UTC().Format("20060102150405")
	t.Message = msg
}

func (t *TagObj) Create() {
	if t.Message == "" {
		t.GenerateMessage()
	}
	output.Debug("create tag:", t)
	Tag(t.Name, t.Message)
}

func (t *TagObj) Push() {
	Push(t.Name)
}
