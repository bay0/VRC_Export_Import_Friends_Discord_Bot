package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Init initialize logrus and the logs.log
func Init() {
	var filename = "logs.log"
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(prefixed.TextFormatter)
	// You can change the Timestamp format. But you have to use the same date and time.
	// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
	// ie "Mon Jan 2 15:04:05 MST 2006" is the reference time. You can't change it
	// Formatter.DisableColors = true
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	Formatter.ForceFormatting = true
	Formatter.ForceColors = true

	log.SetFormatter(Formatter)
	mw := io.MultiWriter(os.Stdout, f)
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		log.SetOutput(ansicolor.NewAnsiColorWriter(mw))
	}
}