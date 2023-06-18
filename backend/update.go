package main

import (
	. "backend/common"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/xor-gate/goexif2/exif"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/edit", editHandler)

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

		processUpdateRequest(r.Context(), string(body))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func processUpdateRequest(ctx context.Context, body string) (int, string) {
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

	err = updatePost(reqData.Post.Current.ID, div)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, "something went wrong"
	}

	return http.StatusOK, "successfully added exif tags"
}

func updatePost(id, html string) error {
	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	u := fmt.Sprintf("%v/posts/%v", "http://localhost:2368/ghost/api/admin", id)
	println("url:", u)
	body := map[string]any{"posts": []any{map[string]any{"updated_at": now, "id": id, "html": html}}}
	reqBody := new(bytes.Buffer)

	err := json.NewEncoder(reqBody).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, u, reqBody)
	if err != nil {
		return err
	}

	token, err := GetJWTToken("GHOST_API_KEY")
	if err != nil {
		return err
	}

	println("token: ", fmt.Sprintf("Ghost %s", token))
	req.Header.Set("Authorization", fmt.Sprintf("Ghost %s", token))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	println(res.StatusCode)

	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(s))

	return nil

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
				figure := s.Closest(".kg-image-card")
				figure.AfterHtml(div)
				addedExif = true
			} else {
				log.Println("exif add error: ", err)
			}
		}
	})

	if addedExif {
		return doc.Find("body").Html()
	}
	return "", errors.New("didn't add any exif tag. See log for the reason")
}

func getExifDiv(url string) (string, error) {
	file, err := fileFromUrl(url)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}()

	x, err := exif.Decode(file)
	if err != nil {
		return "", fmt.Errorf("unable to get image exif: %w", err)
	}
	var sb strings.Builder
	sb.WriteString(`<div class="image-exif">`)
	validExifFound := false

	model, err := x.Get(exif.Model)
	if err == nil {
		val, _ := model.StringVal()
		sb.WriteString(fmt.Sprintf(`<span id="camera">%v</span>`, val))
		validExifFound = true
	}
	fStop, err := x.Get(exif.FNumber)
	if err == nil {
		num, dem, err := fStop.Rat2(0)
		if err == nil {
			sb.WriteString(fmt.Sprintf(`<span id="aperture">F/%.1f</span>`, float64(num)/float64(dem)))
			validExifFound = true
		}
	}
	exposure, err := x.Get(exif.ExposureTime)
	if err == nil {
		num, dem, err := exposure.Rat2(0)
		if err == nil {
			sb.WriteString(fmt.Sprintf(`<span id="shutter">%v/%v</span>`, num, dem))
			validExifFound = true
		}
	}
	iso, err := x.Get(exif.ISOSpeedRatings)
	if err == nil {
		sb.WriteString(fmt.Sprintf(`<span id="iso">ISO %v</span>`, iso.String()))
		validExifFound = true
	}
	focal, err := x.Get(exif.FocalLength)
	if err == nil {
		num, dem, err := focal.Rat2(0)
		if err == nil {
			sb.WriteString(fmt.Sprintf(`<span id="focal">%vmm</span>`, num/dem))
			validExifFound = true
		}
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
	file, err := os.CreateTemp("", "*")
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
