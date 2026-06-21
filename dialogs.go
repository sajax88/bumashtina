package main

import (
	"context"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func ShowInfoDialog(c context.Context, title string, message string) {
	_, err := runtime.MessageDialog(c, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func ShowWarningDialog(c context.Context, title string, message string) {
	_, err := runtime.MessageDialog(c, runtime.MessageDialogOptions{
		Type:    runtime.WarningDialog,
		Title:   title,
		Message: message,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func ShowErrorDialog(c context.Context, title string, message string) {
	_, err := runtime.MessageDialog(c, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   title,
		Message: message,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func ShowQuestionDialog(c context.Context, title string, message string, defaultButton string) string {
	if defaultButton == "" {
		defaultButton = "No"
	}
	answer, err := runtime.MessageDialog(c, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         title,
		Message:       message,
		DefaultButton: defaultButton,
	})

	if err != nil {
		log.Fatal(err)
	}

	return answer
}
