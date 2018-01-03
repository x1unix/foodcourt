package settings

type SMTPSettings struct {
	Host 		string	`msgpack:"host" json:"host" validate:"required,hostname|ip_addr"`
	Port 		int		`msgpack:"port" json:"port" validate:"required,max=65535"`
	User 		string	`msgpack:"username" json:"username" validate:"required,gt=1"`
	Password 	string	`msgpack:"password" json:"password"`
}
