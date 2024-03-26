package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var (
	helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
)

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = helpTemplate
	app.Name = "Files Indexer"
	app.Usage = "Specify the target file to search. If no parameter is provider, only indexing is performed"
	app.Flags = []cli.Flag{}
	app.Version = "v1.0.0"
	app.Authors = []cli.Author{
		{
			Name:  "DrDelphi",
			Email: "drdelphi@gmail.com",
		},
	}
	app.Action = func(c *cli.Context) error {
		return startApp(c)
	}

	err := loadConfig()
	if err != nil {
		fmt.Printf("error loading config file: %s\n\r", err)
		os.Exit(1)
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}

	if len(os.Args) > 1 {
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

func startApp(_ *cli.Context) error {
	targetFile := ""
	if len(os.Args) > 1 {
		targetFile = os.Args[1]
	}
	loadFileStrings()
	index()

	if targetFile != "" {
		fileStrings, err := readStrings(targetFile)
		if err != nil {
			return err
		}

		if len(fileStrings) == 0 {
			return errors.New("no strings found in file")
		}

		res, matchPercent := searchStrings(fileStrings)
		if len(res) == 0 {
			fmt.Println("Sorry. No matches found")
		} else {
			bytes, _ := os.ReadFile(targetFile)
			fmt.Printf("%v file(s) found with %.2f%% strings match:\n\r", len(res), matchPercent)
			for _, f := range res {
				analyseFile(f, bytes)
			}
		}
	}

	return nil
}

func searchStrings(targetS []string) ([]string, float64) {
	maxMatched := 0
	matchPercent := float64(0)
	res := make([]string, 0)
	for fileName, fileStrings := range FileStrings {
		matched := 0
		for _, ts := range targetS {
			for _, fs := range fileStrings {
				if ts == fs {
					matched++
				}
			}
		}
		if matched == 0 || matched < maxMatched {
			continue
		}

		if matched > maxMatched {
			res = make([]string, 0)
			maxMatched = matched
			matchPercent = float64(maxMatched*100) / float64(len(targetS))
		}

		if matchPercent >= appConfig.MinMatchPercent {
			res = append(res, fileName)
		}
	}

	return res, matchPercent
}

func analyseFile(fileName string, target []byte) {
	fmt.Printf("%s - ", fileName)
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("file not found !")
		return
	}

	if len(bytes) != len(target) {
		fmt.Printf("different sizes (target = %v, matched = %v)\n\r", len(target), len(bytes))
		return
	}

	different := 0
	for i := 0; i < len(bytes); i++ {
		if bytes[i] != target[i] {
			different++
		}
	}

	if different == 0 {
		fmt.Println("identical file !")
	} else {
		fmt.Printf("%v different byte(s)\n\r", different)
	}
}
