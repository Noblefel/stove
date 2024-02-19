# STOVE
### Simple Tabular Output View

A small application that converts CSV (Comma-Separated Values) data into a PDF document. First, it reads a CSV file, then creates an HTML representation of the data with a styled table, and converts it into a PDF using the [chromedp](https://github.com/chromedp/chromedp) Package.

# Installation 
```bash
git clone https://github.com/Noblefel/stove
``` 

# Usage
### Basic 
To quickly convert a CSV file to PDF using the default settings, simply run:
```sh
go run main.go
```

### Command Flags
- csv: Specifies the csv file you want to convert 
- output: Defines the name for the resulting PDF 
- html: The template you want to use for rendering the data
- title: Title to be printed in the content header
- number: Show rows number

example:
```sh
go run main.go -output=employees_2024 -title="My Employees" -number=true
```

(no need to the include the file extension) 

# Sample
<img src="https://github.com/Noblefel/stove/blob/main/sample.PNG">