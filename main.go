package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nlopes/slack"
)

func main() {
	var (
		debug = flag.Bool("debug", false, "enable debug on slack client")
		to    = flag.String("to", "", "time.RFC3339 formatted date")
		token = flag.String("token", "", "https://api.slack.com/custom-integrations/legacy-tokens")
	)
	flag.Parse()

	api := slack.New(*token)
	api.SetDebug(*debug)

	deleted := 0
	currentPage := 1 // sentinel
	params := slack.NewGetFilesParameters()

	// Default to up to two weeks ago.
	if *to == "" {
		*to = time.Now().Add(-14 * 24 * time.Hour).Format("2006-01-02")
	}
	t, err := time.Parse("2006-01-02", *to)
	must(err)

	params.TimestampTo = slack.JSONTime(t.Unix())
	params.Count = 1000
	for {
		params.Page = currentPage
		result, paging, err := api.GetFiles(params)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, file := range result {
			for {
				err := api.DeleteFile(file.ID)
				if err != nil {
					fmt.Println(err.Error())
					time.Sleep(5 * time.Second)
				} else {
					deleted++
					break
				}
			}
		}

		if paging.Page >= paging.Pages {
			break
		}
		currentPage++
	}

	if deleted > 0 {
		fmt.Printf("Deleted: %d", deleted)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
