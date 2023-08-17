# cronwrap 🕰️

Task scheduler written in go that wraps other programs.

---

## Installation

```bash
curl -sSL https://github.com/nicolasanjoran/cronwrap/releases/download/1.1.0/install.sh | sudo sh
```

## Features

- Schedule a program by wrapping the original command
- Send email on failure with logs
- Support for healthchecks.io URLs

## Usage

```bash
cronwrap "<cron expression>" <your command here>

# Example: echo "hello world every 5 secs"
cronwrap "0/5 * * * * *" echo "hello world"
```

## Environment variables

| Name             | Description                                                                                                                                                                                                    |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| JOB_NAME         | Give a job a name, will be used for reporting.                                                                                                                                                                 |
| SMTP_SERVER      | example: smtp.gmail.com                                                                                                                                                                                        |
| SMTP_PORT        | example: 587                                                                                                                                                                                                   |
| SMTP_USER        | your smtp username (for gmail: the email address)                                                                                                                                                              |
| SMTP_PASS        | your password                                                                                                                                                                                                  |
| EMAIL_FROM       | email address of the sender (the email address used above)                                                                                                                                                     |
| EMAIL_TO         | email address of the recipient                                                                                                                                                                                 |
| EMAIL_IF_SUCCESS | if "true": send a report by email on both success and failure, default: only on failure                                                                                                                        |
| HEALTHCHECK_URL  | URL of the healthcheck.io check (works with self-hosted instances). The program will automatically hit the /start endpoint before the task and will send the results after the process is completed or failed. |
