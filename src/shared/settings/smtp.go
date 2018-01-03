package settings

type SMTPSettings struct {
	Host 		string	`msgpack:"host" json:"host"`
	Port 		int		`msgpack:"port" json:"port"`
	User 		string	`msgpack:"username" json:"username"`
	Password 	string	`msgpack:"password" json:"password"`
}
