package cmd

import (
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func RunImportFields() {
	// Внутридворовые спортплощадки районов СПб
	url := "https://kfis.gov.spb.ru/infrastruktura/vnutridvorovye-sportploshadki-rajonov-spb/"

	// Загружаем HTML-страницу
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to fetch the page:", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	// Парсим HTML-документ
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Failed to parse the HTML:", err)
	}

	// Находим таблицу (например, по тегу <table> или классу)
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		log.Printf("Table #%d:\n", i+1)

		// Перебираем строки таблицы
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			// Перебираем ячейки в строке
			row.Find("td, th").Each(func(k int, cell *goquery.Selection) {
				log.Printf("Cell %d-%d: %s\n", j+1, k+1, cell.Text())
			})
		})
	})
}