package client_query

import (
	"bufio"
	"fmt"
	"net"
	"strings"

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

	// read message for each connection, but don't process any further
	// this allows reading the response messages more easily
	reader := bufio.NewReader(conn)
	for range make([]int, 4) {
		_, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
	}

	return &ClientQuery{conn: conn}, nil
}

func (c *ClientQuery) Authenticate() error {
	if !viper.IsSet("apikey") || viper.GetString("apikey") == "" {
		return fmt.Errorf("apikey is not set or empty")
	}

	if err := c.Do(fmt.Sprintf("auth apikey=%v\n", viper.GetString("apikey"))); err != nil {
		return err
	}

	return nil
}

func (c *ClientQuery) Do(command string) error {
	fmt.Fprintln(c.conn, command)
	status, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return err
	}
	rsp := parseResponse(status)
	if rsp != "ok" {
		return fmt.Errorf("%s => %s", command, rsp)
	}

	return nil
}

func (c *ClientQuery) Quit() {
	c.conn.Close()
}

func parseResponse(response string) string {
	msg := strings.TrimSpace(strings.ReplaceAll(strings.Join(strings.Split(response, "msg=")[1:], " "), "\\s", " "))

	return msg
}
