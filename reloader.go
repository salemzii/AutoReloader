package autoreloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Extensions []string `json:"extensions"`
	ServerCmd  string   `json:"servercmd"`
	Commands   []string `json:"commands"`
}

func StartAutoReloader() {
	/*path is usually assumed to be the located base directory where all your
	packages are accesible from, usually in the same directory with your main.go
	*/
	path, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	/*
		try to get config file AutoReloaderConfig.json
	*/

	workingDir := filepath.Dir(path)
	config, err := os.Open(workingDir + "/" + "AutoReloaderConfig.json")
	if err != nil {
		log.Println("Starting AutoReloader with config file .......")
	}
	defer config.Close()

	configByte, err := ioutil.ReadAll(config)
	if err != nil {
		log.Println("Error loading AutoReloaderConfig.json")
	}
	var conf Config

	json.Unmarshal(configByte, &conf)
	fmt.Println(conf)
	AutoReloader(conf, workingDir)
}

func AutoReloader(config Config, path string) {

	files_dict := map[string]time.Time{}
	func() {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf(err.Error())
			}
			for _, v := range config.Extensions {

				if strings.HasSuffix(info.Name(), string(v)) {
					files_dict[path] = info.ModTime()
				}
			}
			return nil
		})
		fmt.Println(files_dict)
		for {
			for k, v := range files_dict {
				info, err := os.Stat(k)
				if err != nil {
					fmt.Println(err)
				}
				if info.ModTime() != v {
					log.Printf("File Name: %s\n", info.Name())
					files_dict[k] = info.ModTime()
					if err := config.RunServer(path); err != nil {
						log.Println("Error running server")
						log.Println(err)
					}

				}
			}
		}
	}()

}

func (conf *Config) RunServer(path string) error {

	os.Chdir(path)
	CMD := exec.Command(conf.ServerCmd)
	err := CMD.Run()

	if err != nil {
		return err
	}
	return nil
}
