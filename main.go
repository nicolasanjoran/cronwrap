package main

import (
	"bytes"
	"io"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"strings"

	"github.com/robfig/cron/v3"
)

var (
	smtpServer = os.Getenv("SMTP_SERVER")
	smtpPort   = os.Getenv("SMTP_PORT")
	smtpUser   = os.Getenv("SMTP_USER")
	smtpPass   = os.Getenv("SMTP_PASS")
	fromEmail  = os.Getenv("FROM_EMAIL")
	toEmail    = os.Getenv("TO_EMAIL")
)

type CmdError struct {
	err string
	log string
}

func (e *CmdError) Error() string {
	return e.err + ": " + e.log
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("You need to provide the cron schedule and the command to run.")
	}

	cronSpec := os.Args[1]

	c := cron.New(cron.WithSeconds())
	c.AddFunc(cronSpec, func() {
		err := runTask(os.Args[2:])
		if err != nil {
			sendEmailWithLog(err.Error())
		}
	})

	c.Start()
	select {}
}

func runTask(commandArgs []string) error {
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutDone := streamCopy(&stdoutBuf, stdoutPipe, os.Stdout)
	stderrDone := streamCopy(&stderrBuf, stderrPipe, os.Stderr)

	err = cmd.Start()
	if err != nil {
		return err
	}

	// Wait for the copy operations to complete
	<-stdoutDone
	<-stderrDone

	err = cmd.Wait()
	if err != nil {
		return &CmdError{err: err.Error(), log: "STDOUT:\n" + stdoutBuf.String() + "\n\nSTDERR:\n" + stderrBuf.String()}
	}
	return nil
}

func streamCopy(dst *bytes.Buffer, src io.Reader, additionalDst io.Writer) <-chan bool {
	done := make(chan bool)

	go func() {
		defer close(done)

		_, err := io.Copy(io.MultiWriter(dst, additionalDst), src)
		if err != nil {
			log.Println("Failed to copy stream:", err)
		}
	}()

	return done
}

func sendEmailWithLog(logContent string) {
	body := "Subject: Task Result\r\n\r\n" + logContent

	auth := smtp.PlainAuth("", smtpUser, smtpPass, strings.Split(smtpServer, ":")[0])
	err := smtp.SendMail(smtpServer, auth, fromEmail, []string{toEmail}, []byte(body))
	if err != nil {
		log.Println("Failed to send email:", err)
	}
}
