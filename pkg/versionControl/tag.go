package versionControl

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
)

type Tag struct {
	Name    string
	Message string
}

// SetMessage generates a message for the tag
func (t *Tag) SetMessage() error {
	var msg string
	describe, err := DescribeLong()
	if err != nil {
		return err
	}
	msg += describe
	msg += "-" + time.Now().UTC().Format("20060102150405")
	t.Message = msg
	output.Logger().WithFields(logrus.Fields{
		"tagName":    t.Name,
		"tagMessage": t.Message,
	}).Debug("generated a message for the git tag")
	return nil
}

// Create the tag
func (t *Tag) Create() error {
	if isTagged := IsAlreadyTagged(t.Name); isTagged != false {
		return fmt.Errorf("the current commit has already been tagged with tag %q", t.Name)
	}

	if t.Message == "" {
		if err := t.SetMessage(); err != nil {
			return err
		}
	}

	if err := TagCommit(t.Name, t.Message); err != nil {
		return err
	}

	output.Logger().WithFields(logrus.Fields{
		"tag": fmt.Sprintf("%#v", t),
	}).Debug("git tag created")
	return nil
}

// Push the tag
func (t *Tag) Push() error {
	if err := Push(t.Name); err != nil {
		return err
	}

	output.Logger().WithFields(logrus.Fields{
		"tag": fmt.Sprintf("%#v", t),
	}).Info("git tag pushed to remote")
	return nil
}
