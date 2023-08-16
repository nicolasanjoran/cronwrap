package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	smtpServer     = os.Getenv("SMTP_SERVER")
	smtpPort       = os.Getenv("SMTP_PORT")
	smtpUser       = os.Getenv("SMTP_USER")
	smtpPass       = os.Getenv("SMTP_PASS")
	emailFrom      = os.Getenv("EMAIL_FROM")
	emailTo        = os.Getenv("EMAIL_TO")
	emailIfSuccess = os.Getenv("EMAIL_IF_SUCCESS")
	jobName        = os.Getenv("JOB_NAME")
	healthcheckURL = os.Getenv("HEALTHCHECK_URL")

	taskRunning bool
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("You need to specify a cron schedule as the first argument followed by the command and its arguments.")
		return
	}

	cronSchedule := os.Args[1]
	cmdArgs := os.Args[2:]

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(cronSchedule, func() {
		if taskRunning {
			log.Println("Previous task is still running. Skipping the current schedule.")
			return
		}
		runTask(cmdArgs)
	})
	if err != nil {
		log.Fatalf("Failed to create cron job: %s", err)
		return
	}
	c.Start()
	select {} // Keep the program running
}

func runTask(cmdArgs []string) {
	taskRunning = true
	log.Printf("Running task: %s", strings.Join(cmdArgs, " "))
	defer func() { taskRunning = false }()

	if healthcheckURL != "" {
		_, err := http.Get(healthcheckURL + "/start")
		if err != nil {
			log.Printf("Failed to signal start to healthcheck.io: %s", err)
		}
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	var combinedOutput bytes.Buffer
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error obtaining stdout: %s", err.Error())
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error obtaining stderr: %s", err.Error())
	}

	go streamCopy(os.Stdout, stdoutPipe, &combinedOutput)
	go streamCopy(os.Stderr, stderrPipe, &combinedOutput)

	startTime := time.Now()
	err = cmd.Run()
	duration := time.Since(startTime)

	outputStr := fmt.Sprintf("%s\nTask Duration: %s", combinedOutput.String(), duration)

	if err != nil {
		log.Printf("Task failed: %s", err)
		outputStr = fmt.Sprintf("%s\nError: %s", outputStr, err)
		sendEmail(fmt.Sprintf("Task Failed: %s", jobName), outputStr)
		if healthcheckURL != "" {
			http.Post(healthcheckURL+"/fail", "text/plain", strings.NewReader(outputStr))
		}
	} else {
		log.Printf("Task succeeded in %s", duration)
		if emailIfSuccess == "true" {
			sendEmail(fmt.Sprintf("Task Succeeded: %s", jobName), outputStr)
		}
		if healthcheckURL != "" {
			http.Post(healthcheckURL, "text/plain", strings.NewReader(outputStr))
		}
	}
}

func streamCopy(dst io.Writer, src io.Reader, buf *bytes.Buffer) {
	buffer := make([]byte, 1024)
	for {
		n, err := src.Read(buffer)
		if n > 0 {
			dst.Write(buffer[:n])
			buf.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}
}

func sendEmail(subject, content string) {

	if smtpServer == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || emailFrom == "" || emailTo == "" {
		return
	}

	body := "Subject: " + subject + "\r\n\r\n" + content

	serverAddress := smtpServer + ":" + smtpPort

	// Setup authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)

	// Send the email. smtp.SendMail automatically starts a TLS session if the server supports it.
	err := smtp.SendMail(serverAddress, auth, emailFrom, []string{emailTo}, []byte(body))
	if err != nil {
		log.Println("Failed to send email:", err)
	}
}

