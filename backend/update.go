package main

import (
	. "backend/common"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/xor-gate/goexif2/exif"
	"github.com/yosssi/gohtml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/update", editHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
			return
		}

		fmt.Println("Request body: ", string(body))
		procesUpdateRequest(r.Context(), string(body))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func procesUpdateRequest(ctx context.Context, body string) (int, string) {
	var reqData updateRequestBody
	err := json.Unmarshal([]byte(body), &reqData)
	if err != nil {
		return http.StatusBadRequest, "invalid request body"
	}

	if reqData.Post.Current.PrimaryTag.Slug != Photography {
		return http.StatusNoContent, "not a photography post"
	}

	div, err := addExifDiv(reqData.Post.Current.Html)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, "something went wrong"
	}

	fmt.Println(gohtml.Format(div))

	return http.StatusOK, "successfully added exif tags"
}

func addExifDiv(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("html parse error: %w", err)
	}

	var addedExif = false
	doc.Find(".kg-image-card img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			div, err := getExifDiv(src)
			if err == nil {
				fmt.Println(src, " :: ", gohtml.Format(div))
				figure := s.Closest(".kg-image-card")
				figure.AfterHtml(div)
				addedExif = true
			} else {
				log.Println("exif add error: ", err)
			}
		}
	})

	if addedExif {
		return doc.Html()
	}
	return "", errors.New("didn't add any exif tag. See log for the reason")
}

func getExifDiv(url string) (string, error) {
	f, err := fileFromUrl(url)
	if err != nil {
		return "", err
	}
	x, err := exif.Decode(f)
	if err != nil {
		return "", fmt.Errorf("unable to get image exif: %w", err)
	}

	var sb strings.Builder
	sb.WriteString(`<div class="image-exif">`)
	validExifFound := false

	model, err := x.Get(exif.Model)
	if err == nil {
		sb.WriteString(fmt.Sprintf(`<span id="camera">%v</span>`, string(model.Val)))
		validExifFound = true
	}
	fStop, err := x.Get(exif.FNumber)
	if err == nil {
		sb.WriteString(fmt.Sprintf(`<span id="aperture">f/%v</span>`, string(fStop.Val)))
		validExifFound = true
	}
	exposure, err := x.Get(exif.ExposureTime)
	if err == nil {
		sb.WriteString(fmt.Sprintf(`<span id="shutter">1/%v</span>`, 1/binary.BigEndian.Uint64(exposure.Val)))
		validExifFound = true
	}
	iso, err := x.Get(exif.ISOSpeedRatings)
	if err == nil {
		sb.WriteString(fmt.Sprintf(`<span id="iso">%v</span>`, string(iso.Val)))
		validExifFound = true
	}
	focal, err := x.Get(exif.FocalLength)
	if err == nil {
		sb.WriteString(fmt.Sprintf(`<span id="focal">%vmm</span>`, string(focal.Val)))
		validExifFound = true
	}

	sb.WriteString(`</div>`)

	if validExifFound {
		return sb.String(), nil
	}

	return "", errors.New("all exif data are invalid")
}

func fileFromUrl(url string) (*os.File, error) {
	resp, err := HttpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to download image: %w", err)
	}

	defer resp.Body.Close()
	file, err := os.Create(fmt.Sprintf("temp/exif/%v", filepath.Base(path.Base(url))))
	if err != nil {
		return nil, fmt.Errorf("unable to create image: %w", err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("unable to save image: %w", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("unable to seek image file: %w", err)
	}

	return file, nil
}

type updateRequestBody struct {
	Post struct {
		Current struct {
			ID          string `json:"id"`
			Html        string `json:"html"`
			UpdatedAt   string `json:"updated_at"`
			PublishedAt string `json:"published_at"`
			PrimaryTag  struct {
				Slug string `json:"slug"`
			} `json:"primary_tag"`
		} `json:"current"`
	} `json:"post"`
}
