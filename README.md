# cronwrap 🕰️

Task scheduler written in go that wraps other programs.

---

## Installation

### Using the 1-liner install script
```bash
curl -sSL https://raw.githubusercontent.com/nicolasanjoran/cronwrap/main/install.sh | sudo sh
```

### Manual installation
```bash
git clone git@github.com:nicolasanjoran/cronwrap.git
cd cronwrap

# Optional: re-build the binary
# sh ./build.sh

# replace <binary> with your OS and architecture
sudo cp release/<binary> /usr/local/bin/cronwrap
sudo chmod +x /usr/local/bin/cronwrap
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
