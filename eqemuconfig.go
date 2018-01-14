package eqemuconfig

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Config struct {
	World     World     `json:"world,omitempty" xml:"world,omitempty"`
	Database  Database  `json:"database,omitempty" xml:"database,omitempty"`
	QuestsDir string    `json:"directories>quests,omitempty" xml:"directories>quests,omitempty"`
	Discord   Discord   `json:"discord,omitempty" xml:"discord,omitempty"`
	Twitter   []Twitter `json:"twitter,omitempty" xml:"twitter,omitempty"`
	Github    Github    `json:"github,omitempty" xml:"github,omitempty"`
}

type World struct {
	Shortname string `json:"shortname" xml:"shortname"`
	Longname  string `json:"longname" xml:"longname"`
	Tcp       Tcp    `json:"tcp,omitempty" xml:"tcp,omitempty"`
	Telnet    Telnet `json:"telnet,omitempty" xml:"telnet,omitempty"`
}

//This was used by configs prior to 4/16 patch
type Tcp struct {
	Ip     string `json:"ip,attr" xml:"ip,attr"`
	Port   string `json:"port,attr" xml:"port,attr"`
	Telnet string `json:"telnet,attr" xml:"telnet,attr"`
}

//This is used by configs after 4/16 patch
type Telnet struct {
	Ip      string `json:"ip,attr" xml:"ip,attr"`
	Port    string `json:"port,attr" xml:"port,attr"`
	Enabled string `json:"enabled,attr" xml:"enabled,attr"`
}

type Database struct {
	Host     string `json:"host" xml:"host"`
	Port     string `json:"port" xml:"port"`
	Username string `json:"username" xml:"username"`
	Password string `json:"password" xml:"password"`
	Db       string `json:"db" xml:"db"`
}

type Discord struct {
	Username       string        `json:"username,omitempty" xml:"username,omitempty"`
	Password       string        `json:"password,omitempty" xml:"password,omitempty"`
	ServerID       string        `json:"serverid,omitempty" xml:"serverid,omitempty"`
	ChannelID      string        `json:"channelid,omitempty" xml:"channelid,omitempty"`
	RefreshRate    time.Duration `json:"refreshrate,omitempty" xml:"refreshrate,omitempty"`
	ItemUrl        string        `json:"itemurl,omitempty" xml:"itemurl,omitempty"`
	Channels       []Channel     `json:"channel" xml:"channel"`
	Admins         []Admin       `json:"admin" xml:"admin"`
	TelnetUsername string        `json:"telnetusername,omitempty" xml:"telnetusername,omitempty"`
	TelnetPassword string        `json:"telnetpassword,omitempty" xml:"telnetpassword,omitempty"`
}

type Channel struct {
	ChannelID   string `json:"channelid,attr" xml:"channelid,attr"`
	ChannelName string `json:"channelname,attr" xml:"channelname,attr"`
}

type Admin struct {
	Name string `json:"name,attr" xml:"name,attr"`
	Id   string `json:"id,attr" xml:"id,attr"`
}

type Twitter struct {
	Owner             string `json:"owner,attr,omitempty" xml:"owner,attr,omitempty"`
	ConsumerKey       string `json:"consumerkey,attr,omitempty" xml:"consumerkey,attr,omitempty"`
	ConsumerSecret    string `json:"consumersecret,attr,omitempty" xml:"consumersecret,attr,omitempty"`
	AccessToken       string `json:"accesstoken,attr,omitempty" xml:"accesstoken,attr,omitempty"`
	AccessTokenSecret string `json:"accesstokensecret,attr,omitempty" xml:"accesstokensecret,attr,omitempty"`
}

type Github struct {
	PersonalAccessToken string `json:"personalaccesstoken,attr,omitempty" xml:"personalaccesstoken,attr,omitempty"`
	RepoUser            string `json:"repouser,attr,omitempty" xml:"repouser,attr,omitempty"`
	RepoName            string `json:"reponame,attr,omitempty" xml:"reponame,attr,omitempty"`
	RefreshRate         int64  `json:"refreshrate,attr,omitempty" xml:"refreshrate,attr,omitempty"`
	IssueLabel          string `json:"issuelabel,attr,omitempty" xml:"issuelabel,attr,omitempty"`
	ItemLabel           string `json:"itemlabel,attr,omitempty" xml:"itemlabel,attr,omitempty"`
	NPCLabel            string `json:"npclabel,attr,omitempty" xml:"npclabel,attr,omitempty"`
	CharacterLabel      string `json:"characterlabel,attr,omitempty" xml:"characterlabel,attr,omitempty"`
}

type Graph struct {
	TablePrefix string `json:"tableprefix,attr,omitempty" xml:"tableprefix,attr,omitempty"`
}

var config *Config

func GetConfig() (respConfig *Config, err error) {
	if config != nil {
		respConfig = config
		return
	}

	respConfig, err = loadJson()
	if err != nil {
		respConfig, err = loadXML()
		if err != nil {
			return
		}
	}
	return
}

func loadJson() (respConfig *Config, err error) {
	f, err := os.Open("eqemu_config.json")
	if err != nil {
		//try to load via env variable
		if err.Error() == "no such file or directory" {
			err = fmt.Errorf("Error opening config: %s", err.Error())
			return
		}

		path := os.Getenv("EQEMU_CONFIG")
		if f, err = os.Open(path); err != nil {
			err = fmt.Errorf("Error opening config: %s", err.Error())
			return
		}
	}
	config = &Config{}
	dec := json.NewDecoder(f)
	err = dec.Decode(config)
	if err != nil {
		if !strings.Contains(err.Error(), "EOF") {
			err = fmt.Errorf("Error decoding config: %s", err.Error())
			return
		}

		//This may be a ?> issue, let's fix it.
		bConfig, rErr := ioutil.ReadFile("eqemu_config.json")
		if rErr != nil {
			err = fmt.Errorf("Error reading config: %s", rErr.Error())
			return
		}
		err = json.Unmarshal(bConfig, config)
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

func loadXML() (respConfig *Config, err error) {
	f, err := os.Open("eqemu_config.xml")
	if err != nil {
		//try to load via env variable
		if err.Error() == "no such file or directory" {
			err = fmt.Errorf("Error opening config: %s", err.Error())
			return
		}

		path := os.Getenv("EQEMU_CONFIG")
		if f, err = os.Open(path); err != nil {
			err = fmt.Errorf("Error opening config: %s", err.Error())
			return
		}
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
