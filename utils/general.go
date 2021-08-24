package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func ReadConfigFile(cfg interface{}, path string) {
	f, err := os.Open(path)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func ArrayContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
