package eqemuconfig

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Config struct {
	Shortname string   `xml:"world>shortname"`
	Longame   string   `xml:"world>longname"`
	Database  Database `xml:"database,omitempty"`
	QuestsDir string   `xml:"directories>quests,omitempty"`
	Discord   Discord  `xml:"discord,omitempty"`
}

type Database struct {
	Host     string `xml:"host"`
	Port     string `xml:"port"`
	Username string `xml:"username"`
	Password string `xml:"password"`
	Db       string `xml:"db"`
}

type Discord struct {
	Username  string `xml:"username,omitempty"`
	Password  string `xml:"password,omitempty"`
	ServerID  int64  `xml:"serverid,omitempty"`
	ChannelID int64  `xml:"channelid,omitempty"`
}

func LoadConfig() (config Config, err error) {
	f, err := os.Open("eqemu_config.xml")
	if err != nil {
		err = fmt.Errorf("Error opening config: %s", err.Error())
		return
	}

	dec := xml.NewDecoder(f)
	err = dec.Decode(&config)
	if err != nil {
		err = fmt.Errorf("Error decoding config: %s", err.Error())
		return
	}
	if config.QuestsDir == "" {
		config.QuestsDir = "quests"
	}
	return
}
