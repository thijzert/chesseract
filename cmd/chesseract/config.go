package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

type Config struct {
}

func defaultConfig() Config {
	return Config{}
}

func saveConfig(c Config, fileName string) error {
	fileName = expandFileName(fileName)
	os.MkdirAll(path.Dir(fileName), 0755)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "\t")

	return e.Encode(c)
}

func loadConfig(fileName string) (Config, error) {
	fileName = expandFileName(fileName)
	g, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		} else {
			return Config{}, err
		}
	}
	defer g.Close()

	rv := Config{}

	d := json.NewDecoder(g)
	err = d.Decode(&rv)
	return rv, err
}

func expandFileName(fileName string) string {
	if len(fileName) > 2 && fileName[0:2] == "~/" {
		if usr, err := user.Current(); err == nil {
			fileName = usr.HomeDir + fileName[1:]
		}
	}
	return fileName
}
