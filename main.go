package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	csvName     = flag.String("csv", "example", "The name of the csv file")
	outputName  = flag.String("output", fmt.Sprintf("sample_%d", time.Now().Unix()), "Name for the pdf results")
	htmlName    = flag.String("html", "basic", "The name of the HTML template")
	title       = flag.String("title", "Untitled", "Title to be printed in the header")
	showNumbers = flag.Bool("number", false, "Show rows number")
)

func main() {
	flag.Parse()

	log.Println("Getting data from csv...")
	file, err := os.Open("data/" + *csvName + ".csv")
	if err != nil {
		log.Fatal("Error opening csv file")
	}

	rows, err := readCSV(file)
	if err != nil {
		log.Fatal(err)
	}

	htmlString, err := setupHTML(rows)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	err = chromedp.Run(ctx, printToPDF(&buf, htmlString))
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("output/"+*outputName+".pdf", buf, 0644)
	if err != nil {
		log.Fatalln("Error writing pdf ", err)
	}

	log.Println("Success")
}

func readCSV(file *os.File) (string, error) {
	reader := csv.NewReader(file)
	var rows strings.Builder

	i := 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		var row strings.Builder
		row.WriteString("<tr>")

		if *showNumbers {
			if i == 0 {
				row.WriteString("<td>No</td>")
			} else {
				row.WriteString(fmt.Sprintf("<td>%d</td>", i))
			}
		}

		for _, s := range line {
			row.WriteString(fmt.Sprintf("<td>%s</td>", s))
		}

		row.WriteString("</tr>")
		rows.WriteString(row.String())
		i++
	}

	return rows.String(), nil
}

func setupHTML(rows string) (htmlString string, err error) {
	log.Println("Building the html...")

	htmlBytes, err := os.ReadFile("html/" + *htmlName + ".html")
	if err != nil {
		log.Println("Error reading html file")
		return "", err
	}

	htmlString = strings.Replace(string(htmlBytes), "[%title%]", *title, 1)
	htmlString = strings.Replace(htmlString, "[%rows%]", rows, 1)
	return htmlString, nil
}

func printToPDF(res *[]byte, htmlString string) chromedp.Tasks {
	log.Println("Printing html to pdf...")

	var wg sync.WaitGroup
	wg.Add(1)

	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					wg.Done()
					cancel()
				}
			})
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			tree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(tree.Frame.ID, htmlString).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			wg.Wait()
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).WithDisplayHeaderFooter(true).Do(ctx)
			if err != nil {
				return err
			}

			*res = buf
			return nil
		}),
	}
}
