package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

// Response handles structure of response JSON.
type Response struct {
	Base, Date string
	Rates      map[string]float64
}

// Page :
type Page struct {
	CurrenciesMap  map[string]string
	To, From       string
	ConversionRate float64
}

var page = Page{CurrenciesMap: availableCurrencyMap}

var availableCurrencyMap = map[string]string{
	"USD": "US dollar",
	"JPY": "Japanese yen",
	"BGN": "Bulgarian lev",
	"CZK": "Czech koruna",
	"DKK": "Danish krone",
	"GBP": "Pound sterling",
	"HUF": "Hungarian forint",
	"PLN": "Polish zloty",
	"RON": "Romanian leu",
	"SEK": "Swedish krona",
	"CHF": "Swiss franc",
	"ISK": "Icelandic krona",
	"NOK": "Norwegian krone",
	"HRK": "Croatian kuna",
	"RUB": "Russian rouble",
	"TRY": "Turkish lira",
	"AUD": "Australian dollar",
	"BRL": "Brazilian real",
	"CAD": "Canadian dollar",
	"CNY": "Chinese yuan renminbi",
	"HKD": "Hong Kong dollar",
	"IDR": "Indonesian rupiah",
	"ILS": "Israeli shekel",
	"INR": "Indian rupee",
	"KRW": "South Korean won",
	"MXN": "Mexican peso",
	"MYR": "Malaysian ringgit",
	"NZD": "New Zealand dollar",
	"PHP": "Philippine piso",
	"SGD": "Singapore dollar",
	"THB": "Thai baht",
	"ZAR": "South African rand",
}

func renderTemplate(w http.ResponseWriter, page Page) {
	tmpl, err := template.ParseFiles("./templates/convert.html")

	if err != nil {
		panic(err)
	} else {
		tmpl.Execute(w, page)
	}
}

func main() {
	http.HandleFunc("/", convertHandler)
	http.HandleFunc("/convert", convertHandler)
	fmt.Println("port from main.go: " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	from := r.Form["fromCurrency"]
	to := r.Form["toCurrency"]

	if len(from) == 1 && len(to) == 1 {
		page.ConversionRate = convert(from[0], to[0])
		page.From = from[0]
		page.To = to[0]
	}

	renderTemplate(w, page)
}

func convert(from, to string) float64 {
	url := "https://api.fixer.io/latest?base=" + from

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	body := resp.Body

	defer body.Close()

	jDec := json.NewDecoder(body)

	jResp := &Response{}

	var conversionRate float64

	if err := jDec.Decode(jResp); err == nil {
		djResp := *jResp
		conversionRate = djResp.Rates[to]
	} else {
		panic(err)
	}

	return conversionRate
}
