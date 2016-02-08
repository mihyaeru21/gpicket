package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mihyaeru21/gpicket/model"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Tokens []string `yaml:"tokens,flow"`
}

func CmdLog(c *cli.Context) {
	config := parseConfig()

	if len(config.Tokens) <= 0 {
		fmt.Println("There're no tokens in `~/.gpicket.yaml`.")
		os.Exit(1)
	}

	messages := make(chan model.Message, 100)
	for i := 0; i < len(config.Tokens); i++ {
		token := config.Tokens[i]
		go model.NewSlack(token).Start(messages)
	}

	for {
		log(<-messages)
	}
}

func parseConfig() Config {
	configPath := path.Join(os.Getenv("HOME"), ".gpicket.yaml")
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("Configuration file, `~/.gpicket.yaml` is required.")
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Configuration file has invalid format.")
		os.Exit(1)
	}

	return config
}

func log(message model.Message) {
	fmt.Printf("[%s][#%s][%s]%s\n", message.Team, message.Channel, message.UserID, message.Text)
}
