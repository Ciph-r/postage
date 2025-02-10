package forwards

type Config struct {
	Addr string `env:"ADDR" envDefault:":9090"`
}
