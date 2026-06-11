package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var KomgaURL = os.Getenv("KOMGA_URL")
var KomgaUsername = os.Getenv("KOMGA_USERNAME")
var KomgaPassword = os.Getenv("KOMGA_PASSWORD")
var Port = os.Getenv("KOMGA_SSR_PORT")

var AuthHeader = fmt.Sprintf(
	"Basic %s",
	b64.StdEncoding.Strict().EncodeToString(
		[]byte(fmt.Sprintf(
			"%s:%s",
			KomgaUsername,
			KomgaPassword,
		),
		),
	),
)

func layout(title, body string) string {
	return fmt.Sprintf(`
	  <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <meta name="apple-mobile-web-app-capable" content="yes">
      <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
      <title>%s</title>
      <style>
        body { font-family: monospace, sans-serif; background: #fff; color: #000; padding: 5px; margin: 0 }
        h1 { font-size: 1.5rem; border-bottom: 2px solid #000; padding-bottom: 5px; }
        .read { max-width: 800px; margin: 0 auto; text-align: center; }
        .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(130px, 1fr)); gap: 15px; }
        .card { border: 1px solid #ccc; padding: 8px; text-align: center; text-decoration: none; color: #000; }
        .card:active { background: #000; color: #fff; } /* E-ink accessible click feedback */
        .card img { max-width: 100%%; height: auto; display: block; margin-bottom: 5px; }
        .title { font-weight: bold; font-size: 0.85rem; }
        .subtitle { font-size: 0.75rem; margin-top: 4px; color: #555; }
        .btn { font-family: monospace; text-decoration: none; color: #000; border: 1px solid #000; padding: 5px 10px; }
        .page { width: auto; max-width: 100%%; max-height: 90vh; height: auto; border: 1px solid #000; display: block; object-fit: contain }
      </style>

      <script type="text/javascript">
		(function(document,navigator,standalone) {
			if ((standalone in navigator) && navigator[standalone]) {
				var curnode, location = document.location, stop = /^(a|html)$/i;
				document.addEventListener('click', function(e) {
					curnode = e.target;
					while (!(stop).test(curnode.nodeName)) {
						curnode = curnode.parentNode;
					}
					// If a link was clicked, bypass default Safari behavior and update window internally
					if('href' in curnode && ( curnode.href.indexOf('http') || ~curnode.href.indexOf(location.host) ) ) {
						e.preventDefault();
						location.href = curnode.href;
					}
				}, false);
			}
		})(document,window.navigator,'standalone');
	  </script>
    </head>
    <body>
      %s
    </body>
    </html>
	`, title, body)
}

func main() {
	if KomgaURL == "" || KomgaUsername == "" || KomgaPassword == "" {
		fmt.Fprintln(os.Stderr, "Error: KOMGA_URL, KOMGA_USERNAME, and KOMGA_PASSWORD must be set")
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/libraries", KomgaURL), nil)
		if err != nil {
			fmt.Println("Error making request: ", err)
			return
		}

		req.Header.Set("Authorization", AuthHeader)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body: ", err)
			return
		}

		libraries := []Library{}
		err = json.Unmarshal(body, &libraries)
		if err != nil {
			fmt.Println("Error decoding json: ", err)
			return
		}

		var builder strings.Builder

		for _, lib := range libraries {
			fmt.Fprintf(&builder, `
				<a href="/library/%v" class="card">
					<div class="title">%v</div>
				</a>
				`, lib.ID, lib.Name)
		}

		listHtml := builder.String()

		bodyHtml := fmt.Sprintf(`<h1>Simple Komga</h1><div class="grid">%s</div>`, listHtml)

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, layout("Libraries", bodyHtml))
	})

	mux.HandleFunc("GET /library/{id}", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		libraryID := r.PathValue("id")

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/series?libraryId=%s&size=50", KomgaURL, libraryID), nil)
		if err != nil {
			fmt.Println("Error making request: ", err)
			return
		}

		req.Header.Set("Authorization", AuthHeader)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body: ", err)
			return
		}

		seriesList := SeriesPage{}
		err = json.Unmarshal(body, &seriesList)
		if err != nil {
			fmt.Println("Error decoding json: ", err)
			return
		}

		var builder strings.Builder

		for _, series := range seriesList.Content {
			fmt.Fprintf(&builder, `
				<a href="/series/%v" class="card">
					<img src="/api-thumb/series/%v" alt="%v" loading="lazy" />
					<div class="title">%v</div>
				</a>
				`, series.ID, series.ID, series.Metadata.Title, series.Metadata.Title)
		}

		seriesHtml := builder.String()

		bodyHtml := fmt.Sprintf(`
			<a href="/" class="btn"><- Back</a>
			<h1>Series</h1>
			<div class="grid">%s</div>
		`, seriesHtml)

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, layout("Library Series", bodyHtml))
	})

	mux.HandleFunc("GET /api-thumb/series/{id}", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		seriesID := r.PathValue("id")

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/series/%s/thumbnail", KomgaURL, seriesID), nil)
		if err != nil {
			fmt.Println("Error making request: ", err)
			return
		}

		req.Header.Set("Authorization", AuthHeader)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return
		}

		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}

		w.Header().Add("Content-Type", contentType)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return
		}
	})

	mux.HandleFunc("GET /series/{id}", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		seriesID := r.PathValue("id")

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/series/%s/books?sort=metadata.numberSort,asc&size=100", KomgaURL, seriesID), nil)
		if err != nil {
			fmt.Println("Error making request: ", err)
			return
		}

		req.Header.Set("Authorization", AuthHeader)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body: ", err)
			return
		}

		books := BookPage{}
		err = json.Unmarshal(body, &books)
		if err != nil {
			fmt.Println("Error decoding json: ", err)
			return
		}

		detailsReq, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/series/%s", KomgaURL, seriesID), nil)
		if err != nil {
			fmt.Println("Error making request: ", err)
			return
		}

		detailsReq.Header.Set("Authorization", AuthHeader)

		detailsResp, err := client.Do(detailsReq)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return
		}

		defer detailsResp.Body.Close()

		detailsBody, err := io.ReadAll(detailsResp.Body)
		if err != nil {
			fmt.Println("Error reading body: ", err)
			return
		}

		seriesDetails := SeriesDetails{}
		err = json.Unmarshal(detailsBody, &seriesDetails)
		if err != nil {
			fmt.Println("Error decoding json: ", err)
			return
		}

		readingDirection := "rtl"

		if seriesDetails.Metadata.ReadingDirection == "LEFT_TO_RIGHT" {
			readingDirection = "ltr"
		}

		var builder strings.Builder

		for _, book := range books.Content {
			status := "Unread"
			var startPage int64 = 1
			if book.ReadProgress != nil {
				if book.ReadProgress.Completed {
					status = "Read"
					startPage = 2
				} else {
					status = fmt.Sprintf("Page %d", book.ReadProgress.Page+1)
					startPage = book.ReadProgress.Page
				}
			}
			fmt.Fprintf(&builder, `
				<a href="/read/%v?page=%d&dir=%s&back=/series/%s" class="card">
					<img src="/book-thumb/%v" alt="%v" loading="lazy" />
					<div class="title">%v</div>
					<div class="subtitle">%v</div>
				</a>
				`, book.ID, startPage, readingDirection, seriesID, book.ID, book.Metadata.Title, book.Metadata.Title, status)
		}

		bookHtml := builder.String()

		bodyHtml := fmt.Sprintf(`
			<a href="javascript:history.back()" class="btn"><- Back</a>
			<h1>Books & Volumes</h1>
			<div class="grid">%s</div>
		`, bookHtml)

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, layout("Series Volumes", bodyHtml))
	})

	mux.HandleFunc("GET /book-thumb/{id}", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		bookID := r.PathValue("id")

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/books/%s/thumbnail", KomgaURL, bookID), nil)
		if err != nil {
			fmt.Println("Error making request: ", err)
			return
		}

		req.Header.Set("Authorization", AuthHeader)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return
		}

		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg"
		}

		w.Header().Add("Content-Type", contentType)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return
		}
	})

	mux.HandleFunc("GET /read/{id}", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		bookID := r.PathValue("id")
		query := r.URL.Query()
		backURL := query.Get("back")
		if backURL == "" {
			backURL = "/"
		}
		backURLEscaped := url.QueryEscape(backURL)

		var currentPage int64

		if pageStr := query.Get("page"); pageStr != "" {
			var err error
			currentPage, err = strconv.ParseInt(pageStr, 10, 64)
			if err != nil {
				currentPage = 1
			}
		}

		currentDir := query.Get("dir")
		if currentDir != "rtl" {
			currentDir = "ltr"
		}
		isRTL := currentDir == "rtl"

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/books/%s/pages", KomgaURL, bookID), nil)
		if err != nil {
			fmt.Println("Error making request: ", err)
			return
		}

		req.Header.Set("Authorization", AuthHeader)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body: ", err)
			return
		}

		pages := []BookPageItem{}
		err = json.Unmarshal(body, &pages)
		if err != nil {
			fmt.Println("Error decoding json: ", err)
			return
		}

		totalPages := int64(len(pages))

		// progress sync
		data := ReadProgressUpdateRequest{
			Page:      currentPage,
			Completed: currentPage == totalPages-1,
		}
		jsonBytes, err := json.Marshal(data)

		go func(bID string, pData []byte) {
			patchReq, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/books/%s/read-progress", KomgaURL, bID), bytes.NewBuffer(pData))
			if err != nil {
				fmt.Println("Error making request: ", err)
				return
			}

			patchReq.Header.Set("Authorization", AuthHeader)
			patchReq.Header.Set("Content-Type", "application/json")

			patchResp, err := client.Do(patchReq)
			if err != nil {
				fmt.Println("Error sending request: ", err)
				return
			}
			patchResp.Body.Close()
		}(bookID, jsonBytes)

		prevPage := int64(-1)
		nextPage := int64(-1)
		if currentPage > 1 {
			prevPage = currentPage - 1
		}
		if currentPage < totalPages-1 {
			nextPage = currentPage + 1
		}

		imageClickURL := backURL
		if nextPage >= 0 {
			imageClickURL = fmt.Sprintf("/read/%s?page=%d&dir=%s&back=%s", bookID, nextPage, currentDir, backURLEscaped)
		}

		prevButtonText := "<- PREV"
		if isRTL {
			prevButtonText = "PREV ->"
		}

		prevButtonHTML := fmt.Sprintf(`<span style="visibility: hidden;">%s</span>`, prevButtonText)
		if prevPage >= 0 {
			prevButtonHTML = fmt.Sprintf(`<a href="/read/%s?page=%d&dir=%s&back=%s" class="btn">%s</a>`, bookID, prevPage, currentDir, backURLEscaped, prevButtonText)
		}

		nextButtonHTML := fmt.Sprintf(`<a href="%s" class="btn">FINISH</a>`, backURL)
		if nextPage >= 0 {
			nextButtonText := "NEXT ->"
			if isRTL {
				nextButtonText = "<- NEXT"
			}
			nextButtonHTML = fmt.Sprintf(`<a href="/read/%s?page=%d&dir=%s&back=%s" class="btn">%s</a>`, bookID, nextPage, currentDir, backURLEscaped, nextButtonText)
		}

		toggleDir := "rtl"
		if isRTL {
			toggleDir = "ltr"
		}
		toggleLabel := "⇄ Switch to RTL (Manga)"
		if isRTL {
			toggleLabel = "⇄ Switch to LTR"
		}
		toggleURL := fmt.Sprintf("/read/%s?page=%d&dir=%s&back=%s", bookID, currentPage, toggleDir, backURLEscaped)
		flexDirection := "row"
		if isRTL {
			flexDirection = "row-reverse"
		}

		nextPageHTML := ""
		if !isRTL && nextPage >= 0 {
			nextPageHTML = fmt.Sprintf(`<img src="/book-image/%s/%d" style="display: none;" />`, bookID, nextPage)
		}

		prevPageHTML := ""
		if isRTL && prevPage >= 0 {
			prevPageHTML = fmt.Sprintf(`<img src="/book-image/%s/%d" style="display: none;" />`, bookID, prevPage)
		}

		readerHTML := fmt.Sprintf(`
			<div class="read">
				<div style="margin-bottom: 15px; display: flex; justify-content: space-between; align-items: center;">
					<a href="%s" class="btn"><- Back</a>

					<a href="%s" class="btn">
						%s
					</a>

					<span style="font-family: monospace; font-weight: bold;" >Page %d / %d</span>
				</div>

				<a href="%s" style="display: flex; justify-content: center; align-items: center; outline: none; cursor: pointer;">
					<img
						src="/book-image/%s/%d"
						class="page"
						alt="Book Page %d"
					/>
				</a>

				<div style="margin-top: 15px; display: flex; justify-content: space-between; flex-direction: %s;">
					%s
					%s
				</div>

				%s
				%s

			</div>
		`, backURL, toggleURL, toggleLabel, currentPage+1, totalPages, imageClickURL, bookID, currentPage, currentPage+1, flexDirection, prevButtonHTML, nextButtonHTML, nextPageHTML, prevPageHTML)

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, layout(fmt.Sprintf("Reading Page %d", currentPage+1), readerHTML))
	})

	mux.HandleFunc("GET /book-image/{bookID}/{pageNumber}", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}

		bookID := r.PathValue("bookID")
		pageNumber := r.PathValue("pageNumber")

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/books/%s/pages/%s?convert=jpeg", KomgaURL, bookID, pageNumber), nil)
		if err != nil {
			fmt.Println("Error making request: ", err)
			return
		}

		req.Header.Set("Authorization", AuthHeader)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request: ", err)
			return
		}

		defer resp.Body.Close()

		w.Header().Add("Content-Type", "image/jpeg")
		w.Header().Add("Cache-Control", "public, max-age=3600")
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return
		}
	})

	port := "25601"
	if Port != "" {
		port = Port
	}

	fmt.Printf("Listening on http://127.0.0.1:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
