package eqemuconfig

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	World     World     `xml:"world,omitempty"`
	Database  Database  `xml:"database,omitempty"`
	QuestsDir string    `xml:"directories>quests,omitempty"`
	Discord   Discord   `xml:"discord,omitempty"`
	Twitter   []Twitter `xml:"twitter,omitempty"`
	Github    Github    `xml:"github,omitempty"`
}

type World struct {
	Shortname string `xml:"shortname"`
	Longname  string `xml:"longname"`
	Tcp       Tcp    `xml:"tcp,omitempty"`
}

type Tcp struct {
	Ip     string `xml:"ip,attr"`
	Port   string `xml:"port,attr"`
	Telnet string `xml:"telnet,attr"`
}

type Database struct {
	Host     string `xml:"host"`
	Port     string `xml:"port"`
	Username string `xml:"username"`
	Password string `xml:"password"`
	Db       string `xml:"db"`
}

type Discord struct {
	Username       string    `xml:"username,omitempty"`
	Password       string    `xml:"password,omitempty"`
	ServerID       string    `xml:"serverid,omitempty"`
	ChannelID      string    `xml:"channelid,omitempty"`
	RefreshRate    int64     `xml:"refreshrate,omitempty"`
	Channels       []Channel `xml:"channel"`
	TelnetUsername string    `xml:"username,omitempty"`
	TelnetPassword string    `xml:"password,omitempty"`
}

type Channel struct {
	ChannelID   string `xml:"channelid,attr"`
	ChannelName string `xml:"channelname,attr"`
}

type Twitter struct {
	Owner             string `xml:"owner,attr,omitempty"`
	ConsumerKey       string `xml:"consumerkey,attr,omitempty"`
	ConsumerSecret    string `xml:"consumersecret,attr,omitempty"`
	AccessToken       string `xml:"accesstoken,attr,omitempty"`
	AccessTokenSecret string `xml:"accesstokensecret,attr,omitempty"`
}

type Github struct {
	PersonalAccessToken string `xml:"personalaccesstoken,attr,omitempty"`
	RepoUser            string `xml:"repouser,attr,omitempty"`
	RepoName            string `xml:"reponame,attr,omitempty"`
	RefreshRate         int64  `xml:"refreshrate,attr,omitempty"`
	IssueLabel          string `xml:"issuelabel,attr,omitempty"`
	ItemLabel           string `xml:"itemlabel,attr,omitempty"`
	NPCLabel            string `xml:"npclabel,attr,omitempty"`
	CharacterLabel      string `xml:"characterlabel,attr,omitempty"`
}

type Graph struct {
	TablePrefix string `xml:"tableprefix,attr,omitempty"`
}

var config *Config

func GetConfig() (respConfig *Config, err error) {
	if config != nil {
		respConfig = config
		return
	}

	f, err := os.Open("eqemu_config.xml")
	if err != nil {
		err = fmt.Errorf("Error opening config: %s", err.Error())
		return
	}
	config = &Config{}
	dec := xml.NewDecoder(f)
	err = dec.Decode(config)
	if err != nil {
		if !strings.Contains(err.Error(), "EOF") {
			err = fmt.Errorf("Error decoding config: %s", err.Error())
			return
		}

		//This may be a ?> issue, let's fix it.
		bConfig, rErr := ioutil.ReadFile("eqemu_config.xml")
		if rErr != nil {
			err = fmt.Errorf("Error reading config: %s", rErr.Error())
			return
		}
		strConfig := strings.Replace(string(bConfig), "<?xml version=\"1.0\">", "<?xml version=\"1.0\"?>", 1)
		err = xml.Unmarshal([]byte(strConfig), config)
		if err != nil {
			err = fmt.Errorf("Failed to unmarshal config: %s", err.Error())
			return
		}
	}
	err = f.Close()
	if err != nil {
		err = fmt.Errorf("Failed to close config: %s", err.Error())
		return
	}
	if config.QuestsDir == "" {
		config.QuestsDir = "quests"
	}
	respConfig = config
	return
}
