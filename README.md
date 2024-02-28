# STOVE

A small application that converts CSV (Comma-Separated Values) data into a PDF document. First, it reads a CSV file, then generates the formatted HTML table, and converts it into a PDF using the [chromedp](https://github.com/chromedp/chromedp) package.

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
- file: Specifies the csv file you want to convert 
- out: Defines the name for the resulting PDF 
- html: The template you want to use for rendering the data
- title: Title to be printed in the content header
- num: Show rows number

example:
```sh
go run main.go -out=employees_2024 -title="My Employees" -num=true
```

(no need to the include the file extension) 

# Sample
<img src="https://github.com/Noblefel/stove/blob/main/sample.PNG">