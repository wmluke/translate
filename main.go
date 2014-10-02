// A command line utility to translate [Java ResourceBundle Properties Files](http://docs.oracle.com/javase/tutorial/i18n/resbundle/propfile.html) with Google Translate.
package main

import (
	"os"
	"log"
	"fmt"
	"strings"
	"sort"
	"strconv"
	"errors"
	"github.com/codegangsta/cli"
	"github.com/jmcvetta/napping"
	goproperties "github.com/dmotylev/goproperties"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// Translate the given phrase with Google Translate
func translate(key string, phrase string, source string, target string) (translation string, err error) {
	res := struct {
			Data struct {
				Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
			}
		}{}

	p := napping.Params{
		"key": key,
		"q": phrase,
		"source": source,
		"target": target,
	}

	resp, err := napping.Get("https://www.googleapis.com/language/translate/v2", &p, &res, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	if resp.Status() != 200 {
		err = errors.New("google translate returned "+strconv.Itoa(resp.Status()))
		log.Println("translate failed", source, target, phrase)
		return
	}

	translation = res.Data.Translations[0].TranslatedText

	log.Println("translate ", phrase, translation)

	return
}

// Return a list of property names sorted alphabetically
func keys(p goproperties.Properties) []string {
	keys := make([]string, len(p))
	i := 0
	for k, _ := range p {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func main() {

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options] source_file target_file

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

	app := cli.NewApp()
	app.Name = "translate"
	app.Usage = "translate a Java ResourceBundle Properties file with Google Translate"
	app.Version = "0.1.0"
	app.Author = "Luke Bunselmeyer"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "source, s",
			Usage: "source langauge code",
		},
		cli.StringFlag{
			Name: "target, t",
			Usage: "target langauge code",
		},
		cli.StringFlag{
			Name: "key, k",
			Usage: "Google translate API key",
			EnvVar: "GOOGLE_API_KEY",
		},
	}

	app.Action = func(c *cli.Context) {
		args := c.Args()
		sourceFile := args.Get(0)
		destFile := args.Get(1)

		if sourceFile == "" {
			println("source property files is required")
			return
		}

		if destFile == "" {
			println("destiation property files is required")
			return
		}

		var src = c.String("source")
		if src == "" {
			println("--source is required")
			return

		}

		var trg = c.String("target")
		if trg == "" {
			println("--target is required")
			return
		}

		var apiKey = c.String("key")
		if apiKey == "" {
			println("--key is required")
			return
		}

		props, err := goproperties.Load(sourceFile)
		if err != nil {
			log.Fatal(err)
			return
		}

		out, err := os.Create(destFile)
		if err != nil {
			log.Fatal(err)
			return
		}

		defer out.Close()

		for _, k := range keys(props) {
			v := props[k]
			t, err := translate(apiKey, v, src, trg)
			if err != nil {
				println("Failed to translate " + v)
				continue
			} else {
				if k != "" && t != "" {
					// escape non-ascii unicode characters per http://docs.oracle.com/javase/7/docs/api/java/util/PropertyResourceBundle.html
					te := strings.Trim(fmt.Sprintf("%+q", t), "\"")
					out.WriteString(k + " = " + te + "\n")
				}
			}
		}
	}

	app.Run(os.Args)
}
