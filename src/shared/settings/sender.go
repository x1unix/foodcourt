package settings

type SenderSettings struct {
	Email 				string		`msgpack:"email" json:"email"`
	ReportRecipients	[]string	`msgpack:"orderRecipients" json:"orderRecipients"`
}
