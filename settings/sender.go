package settings

type SenderSettings struct {
	Enable				bool		`msgpack:"enable" json:"enable"`
	Email 				string		`msgpack:"email" json:"email" validate:"required,email"`
	ReportRecipients	[]string	`msgpack:"orderRecipients" json:"orderRecipients" validate:"required,gt=0"`
}
