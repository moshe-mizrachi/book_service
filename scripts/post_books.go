package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Book struct {
	Title          string  `json:"title"`
	Price          float64 `json:"price"`
	AuthorName     string  `json:"author_name"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
}

func main() {
	data := `[
	  {
	    "title": "velit sunt nulla tempor magna ullamco",
	    "price": 193.8985,
	    "author_name": "Vonda Mcgee",
	    "ebook_available": false,
	    "publish_date": "2005-12-13T00:00:00Z"
	  },
	  {
	    "title": "nulla est labore dolore fugiat sit",
	    "price": 28.1156,
	    "author_name": "Nanette Bentley",
	    "ebook_available": false,
	    "publish_date": "1976-12-06T00:00:00Z"
	  },
	  {
	    "title": "laborum cupidatat reprehenderit exercitation enim laboris",
	    "price": 54.5411,
	    "author_name": "Noemi Blanchard",
	    "ebook_available": false,
	    "publish_date": "2003-07-12T00:00:00Z"
	  },
	  {
	    "title": "id eiusmod aute voluptate sint exercitation",
	    "price": 91.7133,
	    "author_name": "Miller Blackburn",
	    "ebook_available": true,
	    "publish_date": "2018-10-29T00:00:00Z"
	  },
	  {
	    "title": "excepteur ex ea reprehenderit aliqua dolore",
	    "price": 104.4118,
	    "author_name": "Reva Silva",
	    "ebook_available": false,
	    "publish_date": "1979-02-26T00:00:00Z"
	  },
	  {
	    "title": "ullamco nisi elit laboris aliqua do",
	    "price": 98.079,
	    "author_name": "Rebecca Blake",
	    "ebook_available": false,
	    "publish_date": "2017-08-05T00:00:00Z"
	  },
	  {
	    "title": "occaecat qui ullamco anim est incididunt",
	    "price": 71.669,
	    "author_name": "Kaufman Chavez",
	    "ebook_available": true,
	    "publish_date": "1992-01-14T00:00:00Z"
	  },
	  {
	    "title": "est proident incididunt et eu est",
	    "price": 81.6193,
	    "author_name": "Hilary Pena",
	    "ebook_available": true,
	    "publish_date": "2009-12-18T00:00:00Z"
	  },
	  {
	    "title": "non eiusmod mollit reprehenderit velit culpa",
	    "price": 124.1798,
	    "author_name": "Byers Medina",
	    "ebook_available": false,
	    "publish_date": "1974-02-25T00:00:00Z"
	  },
	  {
	    "title": "adipisicing magna reprehenderit consectetur quis mollit",
	    "price": 126.2604,
	    "author_name": "Gomez Powell",
	    "ebook_available": false,
	    "publish_date": "2016-10-07T00:00:00Z"
	  }
	]`

	var books []Book
	if err := json.Unmarshal([]byte(data), &books); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	url := "http://localhost:1234/api/v1/books"
	method := "POST"

	var uuids []string
	for i, book := range books {
		fmt.Printf("Processing book %d: %s\n", i+1, book.Title)

		payloadBytes, err := json.Marshal(book)
		if err != nil {
			fmt.Printf("Error marshaling book %d: %v\n", i+1, err)
			continue
		}
		payload := strings.NewReader(string(payloadBytes))

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			fmt.Printf("Error creating request for book %d: %v\n", i+1, err)
			continue
		}
		req.Header.Add("X-Username", "FIVERR")
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request for book %d: %v\n", i+1, err)
			continue
		}
		defer res.Body.Close()

		fmt.Printf("Book %d response status: %d\n", i+1, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response body for book %d: %v\n", i+1, err)
			continue
		}

		var response map[string]string
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Printf("Error unmarshaling response for book %d: %v\n", i+1, err)
			continue
		}
		if id, ok := response["id"]; ok {
			uuids = append(uuids, id)
		} else {
			fmt.Printf("No ID found in response for book %d\n", i+1)
		}
	}

	fmt.Println("UUIDs:")
	for _, uuid := range uuids {
		fmt.Println(uuid)
	}
}
