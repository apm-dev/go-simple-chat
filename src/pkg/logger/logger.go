package logger

import (
	"log"
)

func Error(err error) {
	print("Error", err)
}

func Debug(err error) {
	print("Debug", err)
}

func Info(err error) {
	print("Info", err)
}

func print(lvl string, err error) {
	log.Printf("\n%s:\n%+v\n", lvl, err)
}
