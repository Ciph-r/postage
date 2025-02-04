package health

type Config struct {
	Addr string `env:"ADDR" envDefault:":9000"`
}
