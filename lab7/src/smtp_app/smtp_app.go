package main

import (
	"bytes"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"math/rand"
	"net/smtp"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Row struct {
	Name    string
	Email   string
	Message string
}

const INDEX_HTML = `
	<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
	<html>
		<head>
		</head>
		<body>
			<table cellpadding="0" cellspacing="0" border="0" width="100%"
				style="background: whitesmoke; min-width: 320px; font-size: 1px; line-height: normal;">
				<tr>
					<td align="center" valign="top">
						<table cellpadding="0" cellspacing="0" border="0" width="700"
							style="background: steelblue; color: whitesmoke; font-family: Arial, Helvetica, sans-serif;">
							<tr>
								<td align="center" valign="top">
									<span style="font-size: 20px; font-weight: bold;
										line-height: 40px; -webkit-text-size-adjust:none; display: block;">
										Здравствуйте, {{.Name}}!
									</span>
									<hr width="600" size="1" color="whitesmoke" noshade>
									<span style="font-size: 16px; font-style: italic;
										line-height: 40px; -webkit-text-size-adjust:none; display: block;">
										{{.Message}}
									</span>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
	</html>
	`

var (
	db_hostname = flag.String("db_host", "students.yss.su", "database host")
	db_name     = flag.String("db_name", "iu9networkslabs", "database name")
	db_username = flag.String("db_username", "iu9networkslabs", "database login")
	db_password = flag.String("db_password", "Je2dTYr6", "database password")

	hostname = flag.String("hostname", "smtp.gmail.com", "authentication hostname")
	hostport = flag.String("hostport", "587", "authentication host port")
	username = flag.String("username", "ilya.afanasyev.26@gmail.com", "authentication username")
	password = flag.String("password", "", "authentication password") // пароль указываю вручную

	indexHtml = template.Must(template.New("index").Parse(INDEX_HTML))
)

func main() {
	flag.Parse()
	auth := smtp.PlainAuth("", *username, *password, *hostname)

	log.Println("Opening database...")
	db, err := sql.Open("mysql", *db_username+":"+*db_password+"@tcp("+*db_hostname+")/"+*db_name)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	log.Println("Querying rows...")
	rows, err := db.Query("SELECT name, email, message FROM iu9afanasyev")
	if err != nil {
		log.Fatal("Failed to make a request: ", err)
	}
	defer rows.Close()

	rand.Seed(1234)
	for rows.Next() {
		var row Row
		err := rows.Scan(&row.Name, &row.Email, &row.Message)
		if err != nil {
			log.Fatal("Failed to scan row: ", err)
		}

		var buf bytes.Buffer
		err = indexHtml.Execute(&buf, row)
		if err != nil {
			log.Fatal("Failed to execute: ", err)
		}

		msg := "To: " + row.Email + "\r\n"
		msg += "Subject: " + "Афанасьев Илья, ИУ9-31Б" + "\r\n"
		msg += "Content-Type: text/html\r\n"
		msg += "\r\n"
		msg += buf.String() + "\r\n"

		dur := time.Duration(rand.Intn(60)) * time.Second
		log.Println("Sleeping for", dur, "...")
		time.Sleep(dur)

		err = smtp.SendMail(*hostname+":"+*hostport, auth, *username, []string{row.Email}, []byte(msg))
		if err != nil {
			log.Println(err)
		}
		log.Println("Successfully sent to " + row.Email)
	}
}
