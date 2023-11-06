package main

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/sergi/go-diff/diffmatchpatch"
	"os"
	"parsContent/urls"
	"strings"
	"time"
)

func main() {
	fileContent1, err := os.Open("text001.txt")
	if err != nil {
		panic(err)
	}
	fileContent2, err := os.Open("text7.txt")
	if err != nil {
		panic(err)
	}
	scanner1 := bufio.NewScanner(fileContent1)
	scanner2 := bufio.NewScanner(fileContent2)
	var text1, text2 string
	for scanner1.Scan() {
		text1 += scanner1.Text() + "\n"
	}
	for scanner2.Scan() {
		text2 += scanner2.Text() + "\n"
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	html := dmp.DiffPrettyText(diffs)

	fmt.Println(html)
	content := scrapPage(urls.Links)
	createWord("newtext3.txt", content)
}
func scrapPage(urls []string) []string {
	var content []string
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{
		Parallelism: 2,               // Максимальное количество одновременных запросов
		Delay:       2 * time.Second, // Задержка между запросами
	})
	c.OnHTML("div.content.seo", func(e *colly.HTMLElement) {
		seoText := e.Text
		content = append(content, seoText)

	})
	for _, url := range urls {

		c.Visit(url)
	}

	return content
}
func createWord(filename string, content []string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	data := strings.Join(content, "\n========================\n")
	_, err = writer.WriteString(data)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
	writer.Flush()

}
