package sockets

type Config struct {
	Addr string `env:"ADDR" envDefault:":8080"`
}
