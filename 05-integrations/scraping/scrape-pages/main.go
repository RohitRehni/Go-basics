// Topic: scrape-pages | Concept: scraping | Stage: 05-integrations
// Run:   go run ./05-integrations/scraping/scrape-pages
//
// Explore:
//   - Scraping = fetch HTML, parse it into a tree, pick out the data you need
//   - golang.org/x/net/html parses real-world (messy) HTML into nodes you can walk
//   - Select by structure and class names (div.product > h2.name), not by line numbers
//   - Follow the "next" link until there is none: that is how pagination works
//   - Resolve relative links against the page URL (url.ResolveReference)
//   - Be polite: read robots.txt, send an honest User-Agent, rate-limit requests,
//     cache what you fetched, and respect the site's terms of use
//   - Prefer an official API or feed when one exists: HTML layouts change without notice
//   - This lesson serves its own fake shop on 127.0.0.1, so it never touches the internet
package main

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// --- The fake website -------------------------------------------------------

var catalogue = [][]struct {
	name  string
	price float64
}{
	{{"Mechanical Keyboard", 89.00}, {"Wireless Mouse", 24.50}, {"USB-C Hub", 39.99}},
	{{"27\" Monitor", 229.00}, {"Laptop Stand", 34.00}, {"Webcam HD", 59.90}},
	{{"Noise-Cancelling Headphones", 149.00}, {"Desk Lamp", 27.75}, {"Cable Organiser", 9.99}},
}

func fakeShop() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "User-agent: *\nDisallow: /admin\nCrawl-delay: 1\n")
	})
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 || page > len(catalogue) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<html><head><title>Shop page %d</title></head><body>\n<h1>Products</h1>\n", page)
		for i, p := range catalogue[page-1] {
			id := (page-1)*3 + i + 1
			fmt.Fprintf(w, `<div class="product card"><h2 class="name">%s</h2>
<span class="price">₹ %.2f</span> <a href="/products/%d">details</a></div>
`, html.EscapeString(p.name), p.price, id)
		}
		if page < len(catalogue) {
			fmt.Fprintf(w, `<nav><a class="next" href="/products?page=%d">Next &raquo;</a></nav>`, page+1)
		}
		fmt.Fprint(w, "</body></html>")
	})
	return httptest.NewServer(mux)
}

// --- The scraper --------------------------------------------------------------

type Product struct {
	Name  string
	Price float64
	Link  string
}

func hasClass(n *html.Node, class string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" && slices.Contains(strings.Fields(a.Val), class) {
			return true
		}
	}
	return false
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func text(n *html.Node) string {
	var b strings.Builder
	for c := range n.Descendants() {
		if c.Type == html.TextNode {
			b.WriteString(c.Data)
		}
	}
	return strings.TrimSpace(b.String())
}

// find returns every element node below root for which match is true.
func find(root *html.Node, match func(*html.Node) bool) []*html.Node {
	var out []*html.Node
	for n := range root.Descendants() {
		if n.Type == html.ElementNode && match(n) {
			out = append(out, n)
		}
	}
	return out
}

func parsePrice(s string) (float64, error) {
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "₹"))
	return strconv.ParseFloat(s, 64)
}

// disallowed reads robots.txt and returns the Disallow paths for every user agent.
func disallowed(client *http.Client, base *url.URL) []string {
	resp, err := client.Get(base.ResolveReference(&url.URL{Path: "/robots.txt"}).String())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var paths []string
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		if path, ok := strings.CutPrefix(sc.Text(), "Disallow:"); ok {
			paths = append(paths, strings.TrimSpace(path))
		}
	}
	return paths
}

func scrape(ctx context.Context, start string) ([]Product, int, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	base, _ := url.Parse(start)
	blocked := disallowed(client, base)
	fmt.Println("robots.txt disallows:", blocked)

	limiter := time.NewTicker(150 * time.Millisecond) // at most ~6 requests per second
	defer limiter.Stop()

	var products []Product
	pages := 0
	next := start
	for next != "" {
		<-limiter.C
		pageURL, _ := url.Parse(next)
		if slices.ContainsFunc(blocked, func(p string) bool { return p != "" && strings.HasPrefix(pageURL.Path, p) }) {
			return products, pages, fmt.Errorf("robots.txt disallows %s", pageURL.Path)
		}
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, next, nil)
		req.Header.Set("User-Agent", "learning-scraper/1.0 (+https://example.com/contact)")
		resp, err := client.Do(req)
		if err != nil {
			return products, pages, err
		}
		doc, err := html.Parse(resp.Body)
		resp.Body.Close()
		if err != nil {
			return products, pages, err
		}
		pages++
		fmt.Printf("  fetched page %d: %s (status %d)\n", pages, next, resp.StatusCode)

		for _, card := range find(doc, func(n *html.Node) bool { return n.Data == "div" && hasClass(n, "product") }) {
			var p Product
			if names := find(card, func(n *html.Node) bool { return hasClass(n, "name") }); len(names) > 0 {
				p.Name = text(names[0])
			}
			if prices := find(card, func(n *html.Node) bool { return hasClass(n, "price") }); len(prices) > 0 {
				p.Price, _ = parsePrice(text(prices[0]))
			}
			if links := find(card, func(n *html.Node) bool { return n.Data == "a" }); len(links) > 0 {
				href, _ := url.Parse(attr(links[0], "href"))
				p.Link = pageURL.ResolveReference(href).String() // make the link absolute
			}
			products = append(products, p)
		}

		next = ""
		if nexts := find(doc, func(n *html.Node) bool { return n.Data == "a" && hasClass(n, "next") }); len(nexts) > 0 {
			href, _ := url.Parse(attr(nexts[0], "href"))
			next = pageURL.ResolveReference(href).String()
		}
	}
	return products, pages, nil
}

func main() {
	fmt.Println("Topic: scraping paginated pages with golang.org/x/net/html")
	shop := fakeShop()
	defer shop.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()
	products, pages, err := scrape(ctx, shop.URL+"/products?page=1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "scrape failed:", err)
		os.Exit(1)
	}

	fmt.Printf("\n  %-30s %10s  %s\n", "NAME", "PRICE", "LINK")
	total := 0.0
	for _, p := range products {
		fmt.Printf("  %-30s %10.2f  %s\n", p.Name, p.Price, strings.TrimPrefix(p.Link, shop.URL))
		total += p.Price
	}
	fmt.Printf("\n%d products from %d pages in %v | total %.2f | average %.2f\n",
		len(products), pages, time.Since(start).Round(10*time.Millisecond), total, total/float64(len(products)))

	want := 0.0
	for _, page := range catalogue {
		for _, p := range page {
			want += p.price
		}
	}
	if len(products) != 9 || pages != 3 || math.Abs(total-want) > 0.001 || products[3].Name != `27" Monitor` {
		fmt.Println("MISMATCH: scraped data differs from the catalogue")
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
