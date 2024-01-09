package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	defaultCSVName    = "example"
	defaultOutputName = fmt.Sprintf("Sample_%d", time.Now().Unix())
	defaultHTMLName   = "basic"
	defaultTitle      = fmt.Sprintf("The Title (%s)", time.Now().Format("2006-01-02"))
)

func main() {

	csvName := flag.String("csv", defaultCSVName, "The name of the csv file")
	outputName := flag.String("output", defaultOutputName, "Name for the pdf results")
	htmlName := flag.String("html", defaultHTMLName, "The name of the HTML template")
	title := flag.String("title", defaultTitle, "Title to be printed in the header")

	flag.Parse()

	data, err := getCSVData(*csvName)
	if err != nil {
		log.Fatal(err)
	}

	htmlString, err := setupHTML(*htmlName, *title, data)
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

func getCSVData(csvName string) (data [][]string, err error) {
	log.Println("Getting data from csv...")

	file, err := os.Open("data/" + csvName + ".csv")
	if err != nil {
		log.Println("Error opening csv file")
		return data, err
	}

	data, err = csv.NewReader(file).ReadAll()
	if err != nil {
		log.Println("Error reading csv file")
		return data, err
	}

	return data, nil
}

func setupHTML(htmlName, title string, data [][]string) (htmlString string, err error) {
	log.Println("Building the html...")

	htmlBytes, err := os.ReadFile("html/" + htmlName + ".html")
	if err != nil {
		log.Println("Error reading html file")
		return "", err
	}

	var rows string

	for _, slice := range data {
		tr := "<tr>"
		for _, s := range slice {
			tr += fmt.Sprintf("<td>%s</td>", s)
		}
		tr += "</tr>"
		rows += tr
	}

	htmlString = strings.Replace(string(htmlBytes), "[%rows%]", rows, 1)
	htmlString = strings.Replace(htmlString, "[%title%]", title, 1)
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
