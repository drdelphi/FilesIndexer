package main

import (
	"encoding/json"
	"os"
)

const (
	fileStringsFilename = "strings.json"
	configFilename      = "config.json"
)

type AppConfig struct {
	MinStringLen    int     `json:"minStringLen"`
	MinMatchPercent float64 `json:"minMatchPercent"`
}

var (
	FileStrings = make(map[string][]string)
	appConfig   AppConfig
)

func loadConfig() error {
	bytes, err := os.ReadFile(configFilename)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &appConfig)
}

func loadFileStrings() error {
	bytes, err := os.ReadFile(fileStringsFilename)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &FileStrings)
}

func saveFileStrings() error {
	bytes, err := json.MarshalIndent(FileStrings, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileStringsFilename, bytes, 0666)
}

func readStrings(fileName string) ([]string, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	n := len(bytes)
	for i := 0; i < n-appConfig.MinStringLen; i++ {
		isAscii := true
		for j := 0; j < appConfig.MinStringLen; j++ {
			if !isDigitOrLetter(bytes[i+j]) {
				isAscii = false
				break
			}
		}
		if !isAscii {
			continue
		}

		sb := bytes[i : i+appConfig.MinStringLen]
		i += appConfig.MinStringLen
		for i < n && isDigitOrLetter(bytes[i]) {
			sb = append(sb, bytes[i])
			i++
		}

		s := string(sb)
		found := false
		for _, oldS := range res {
			if s == oldS {
				found = true
				break
			}
		}
		if !found {
			res = append(res, s)
		}
	}

	return res, nil
}

func isDigitOrLetter(b byte) bool {
	isCapitalLetter := b >= 'A' && b <= 'Z'
	isSmallLetter := b >= 'a' && b <= 'z'
	isDigit := b >= '0' && b <= '9'

	return isCapitalLetter || isSmallLetter || isDigit
}
