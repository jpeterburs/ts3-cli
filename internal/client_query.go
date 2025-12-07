package client_query

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
)

type ClientQuery struct {
	conn net.Conn
}

func Dial() (*ClientQuery, error) {
	conn, err := net.Dial("tcp", viper.GetString("host"))
	if err != nil {
		return nil, err
	}

	return &ClientQuery{conn: conn}, nil
}

func (c *ClientQuery) Authenticate() error {
	if !viper.IsSet("apikey") && viper.GetString("apikey") == "" {
		return fmt.Errorf("apikey is not set or empty")
	}

	fmt.Fprintf(c.conn, "auth apikey=%v\n", viper.GetString("apikey"))

	return nil
}

func (c *ClientQuery) Do(command string) {
	fmt.Fprintln(c.conn, command)
}

func (c *ClientQuery) Quit() {
	c.conn.Close()
}
