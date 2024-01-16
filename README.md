# STOVE
### Simple Tabular Output View with Ease

A small application that converts CSV (Comma-Separated Values) data into a PDF document. The application reads a CSV file, creates an HTML representation of the data with a styled table, and then converts this HTML into a PDF using the [chromedp](https://github.com/chromedp/chromedp) Package.

<br>

# Installation 
```bash
git clone https://github.com/Noblefel/stove
``` 

<br>


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