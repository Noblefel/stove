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

func main() {
	var (
		file  = flag.String("file", "example", "The name of the csv file")
		out   = flag.String("out", fmt.Sprintf("sample_%d", time.Now().Unix()), "Name for the pdf result")
		html  = flag.String("html", "default", "The name of the HTML template")
		title = flag.String("title", "Untitled", "Title to be printed in the header")
		num   = flag.Bool("num", false, "Show rows number")
	)
	flag.Parse()

	f, err := os.Open("data/" + *file + ".csv")
	if err != nil {
		log.Fatal("opening file: ", err)
	}
	defer f.Close()

	rows, err := readCSV(f, *num)
	if err != nil {
		log.Fatal("readCSV: ", err)
	}

	htmlBytes, err := os.ReadFile("html/" + *html + ".html")
	if err != nil {
		log.Fatal("reading html: ", err)
	}

	htmlString := setupHTML(htmlBytes, rows, *title)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	err = chromedp.Run(ctx, printToPDF(&buf, htmlString))
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("output/"+*out+".pdf", buf, 0644)
	if err != nil {
		log.Fatal("writing pdf: ", err)
	}

	log.Println("Done")
}

// readCSV returns the formatted table rows
func readCSV(r io.Reader, num bool) (string, error) {
	reader := csv.NewReader(r)
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

		if num {
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

// setupHTML builds the rest of the html string
func setupHTML(html []byte, rows, title string) string {
	htmlString := strings.Replace(string(html), "[%title%]", title, 1)
	htmlString = strings.Replace(htmlString, "[%rows%]", rows, 1)
	return htmlString
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
