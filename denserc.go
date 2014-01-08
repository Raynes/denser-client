// A client for the denser server. Basically just a daemon that pings the
// server every 5 minutes with the current IP address.
package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"net/http"
	"os/user"
	"path/filepath"
	"time"
)

type Config struct {
	Endpoint string
}

func ConfigPath() string {
	user, _ := user.Current()
	return filepath.Join(user.HomeDir, ".denser")
}

func ReadConfig() (config Config) {
	data, _ := ioutil.ReadFile(ConfigPath())
	toml.Decode(string(data), &config)
	return
}

func IpAddress() string {
	resp, _ := http.Get("http://icanhazip.com")
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func SetIpAddress(config Config) {
	ip := IpAddress()
	fmt.Println("Setting IP address to", ip)
	endpoint := fmt.Sprintf("http://%v/%v", config.Endpoint, ip)
	req, _ := http.NewRequest("PUT", endpoint, nil)

	http.DefaultClient.Do(req)
}

func main() {
	config := ReadConfig()
	for {
		SetIpAddress(config)
		time.Sleep(time.Duration(5) * time.Minute)
	}
}
