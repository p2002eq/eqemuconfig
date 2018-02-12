package eqemuconfig

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type JsonConfig struct {
	Config *Config `json:"server,omitempty"`
}

type Config struct {
	World     World     `json:"world,omitempty" xml:"world,omitempty"`
	Database  Database  `json:"database,omitempty" xml:"database,omitempty"`
	QuestsDir string    `json:"directories>quests,omitempty" xml:"directories>quests,omitempty"`
	Discord   Discord   `json:"discord,omitempty" xml:"discord,omitempty"`
	Twitter   []Twitter `json:"twitter,omitempty" xml:"twitter,omitempty"`
	Github    Github    `json:"github,omitempty" xml:"github,omitempty"`
	NATS      NATS      `json:"nats,omitempty" xml:"nats,omitempty"`
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

type NATS struct {
	Host string `json:"host" xml:"host"`
	Port string `json:"port" xml:"port"`
}

type Discord struct {
	Username          string `json:"username,omitempty" xml:"username,omitempty"`
	Password          string `json:"password,omitempty" xml:"password,omitempty"`
	ClientID          string `json:"clientid,omitempty" xml:"clientid,omitempty"`
	ServerID          string `json:"serverid,omitempty" xml:"serverid,omitempty"`
	ChannelID         string `json:"channelid,omitempty" xml:"channelid,omitempty"`
	CommandChannelID  string `json:"commandchannelid,omitempty" xml:"commandchannelid,omitempty"`
	RefreshRate       time.Duration
	RefreshRateString string    `json:"refreshrate,omitempty" xml:"refreshrate,omitempty"`
	ItemUrl           string    `json:"itemurl,omitempty" xml:"itemurl,omitempty"`
	Channels          []Channel `json:"channel" xml:"channel"`
	Admins            []Admin   `json:"admin" xml:"admin"`
	TelnetUsername    string    `json:"telnetusername,omitempty" xml:"telnetusername,omitempty"`
	TelnetPassword    string    `json:"telnetpassword,omitempty" xml:"telnetpassword,omitempty"`
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

var configInstance *Config

func GetConfig() (config *Config, err error) {
	if configInstance != nil {
		config = configInstance
		return
	}

	config, err = loadJson()
	if err != nil {
		lastErr := err
		config, err = loadXML()
		if err != nil {
			err = errors.Wrapf(err, "failed to load json: %s,  xml", lastErr.Error())
			return
		}
	}
	configInstance = config
	return
}

func loadJson() (config *Config, err error) {
	f, err := os.Open("eqemu_config.json")
	if err != nil {
		//how about via env variable?
		path := os.Getenv("EQEMU_CONFIG")
		if f, err = os.Open(path); err != nil {
			err = errors.Wrap(err, "failed to open config")
			return
		}
		err = nil
	}
	jsonConfig := &JsonConfig{}
	dec := json.NewDecoder(f)
	err = dec.Decode(jsonConfig)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshal config")
		return
	}
	config = jsonConfig.Config
	err = f.Close()
	if err != nil {
		err = errors.Wrap(err, "failed to close config")
		return
	}
	if config.QuestsDir == "" {
		config.QuestsDir = "quests"
	}

	//first see if it's an integer
	refreshRate, err := strconv.ParseInt(config.Discord.RefreshRateString, 64, 10)
	if err != nil { //not an integer
		//is it a proper duration?
		config.Discord.RefreshRate, err = time.ParseDuration(config.Discord.RefreshRateString)
		//not a proper duration
		if err != nil {
			//use default 5 seconds
			config.Discord.RefreshRate = 5 * time.Second
			err = nil
			return
		}
		//use parsed time
		return
	}
	config.Discord.RefreshRate, _ = time.ParseDuration(fmt.Sprintf("%ds", refreshRate))
	return
}

func loadXML() (config *Config, err error) {
	f, err := os.Open("eqemu_config.xml")
	if err != nil {
		path := os.Getenv("EQEMU_CONFIG")
		if f, err = os.Open(path); err != nil {
			err = errors.Wrap(err, "failed to open config")
			return
		}
	}
	config = &Config{}
	dec := xml.NewDecoder(f)
	err = dec.Decode(config)
	if err != nil {
		if !strings.Contains(err.Error(), "EOF") {
			err = errors.Wrap(err, "failed to decode config")
			return
		}

		//This may be a ?> issue, let's fix it.
		bConfig, rErr := ioutil.ReadFile("eqemu_config.xml")
		if rErr != nil {
			err = errors.Wrap(err, "failed to read config for repair")
			return
		}
		strConfig := strings.Replace(string(bConfig), "<?xml version=\"1.0\">", "<?xml version=\"1.0\"?>", 1)
		err = xml.Unmarshal([]byte(strConfig), config)
		if err != nil {
			err = errors.Wrap(err, "failed to unmarshal config after repair")
			return
		}
	}
	err = f.Close()
	if err != nil {
		err = errors.Wrap(err, "failed to close config")
		return
	}
	if config.QuestsDir == "" {
		config.QuestsDir = "quests"
	}
	//first see if it's an integer
	refreshRate, err := strconv.ParseInt(config.Discord.RefreshRateString, 64, 10)
	if err != nil { //not an integer
		//is it a proper duration?
		config.Discord.RefreshRate, err = time.ParseDuration(config.Discord.RefreshRateString)
		//not a proper duration
		if err != nil {
			//use default 5 seconds
			config.Discord.RefreshRate = 5 * time.Second
			err = nil
			return
		}
		//use parsed time
		return
	}
	config.Discord.RefreshRate, _ = time.ParseDuration(fmt.Sprintf("%ds", refreshRate))
	return
}
