package main

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/gocolly/colly"
	"os"
)

func main() {
	collector := colly.NewCollector()

	defaultLink := "https://centraluniversity.ru/"

	visitedURLs, g := collyRun(collector, defaultLink)

	fmt.Println("________________________________")
	fmt.Println("Посещенные сайты: ")
	for _, url := range visitedURLs {
		fmt.Println(url)
	}

	// Сохраняем граф в файл
	file, _ := os.Create("site-graph.gv")
	_ = draw.DOT(g, file)
}

func collyRun(c *colly.Collector, defaultLink string) ([]string, graph.Graph[string, string]) {
	var visitedURLs []string
	visitedMap := make(map[string]bool)
	g := graph.New(graph.StringHash, graph.Directed())

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		parentURL := e.Request.URL.String()
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if len(visitedMap) > 5 {
			return
		}

		// Добавляем вершину для родительского URL, если её ещё нет
		if !visitedMap[parentURL] {
			_ = g.AddVertex(parentURL)
			visitedMap[parentURL] = true
			visitedURLs = append(visitedURLs, parentURL)
		}

		// Проверяем, что link не равен parentURL и не является пустым
		if link != "" && link != parentURL {
			// Добавляем вершину для нового URL, если её ещё нет
			if !visitedMap[link] {
				_ = g.AddVertex(link)
				visitedMap[link] = true
				visitedURLs = append(visitedURLs, link)
				fmt.Printf("Найдена ссылка: %s\n", link)
			}

			// Пытаемся добавить ребро. Если оно уже существует, AddEdge вернет ошибку
			err := g.AddEdge(parentURL, link)
			if err != nil {
				// Ребро уже существует, игнорируем ошибку
				fmt.Printf("Ребро уже существует: %s -> %s\n", parentURL, link)
			}

			c.Visit(link)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Посещаем: %s\n", r.URL.String())
	})

	c.Visit(defaultLink)

	return visitedURLs, g
}
