package bookmarks

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/philipjkim/goreadability"
)

type urlImage struct {
	URL  string
	Size string
}

type preview struct {
	Title          string
	Description    string
	BannerImageURL string
}

const previewTemplate = `
	<html>
		<head><title>{{.Title}}</title></head>
		<body>
			<div style='width: 400px'>
				<div><img alt='{{.Title}}' width='400px' height='200px' src='{{.BannerImageURL}}' /></div>
				<div style='width: 400px'>Title: {{.Title}}</div>
				<div style='width: 400px'>Description: {{.Description}}</div>
			</div>
		</body>
	</html>
	`

func (b *BookmarkManager) savePreview(url, location string) error {
	opt := readability.NewOption()
	opt.ImageRequestTimeout = 3000
	content, err := readability.Extract(url, opt)
	if err != nil {
		return err
	}

	t, err := template.New("preview").Parse(previewTemplate)
	if err != nil {
		return err
	}

	title := content.Title
	description := content.Description
	var bannerImageURL string

	if len(content.Images) > 0 {
		imageURL := content.Images[0].URL
		bannerImageURL, err = downloadImage(imageURL, location)
		if err != nil {
			fmt.Printf("unable to download image %s \n", imageURL)
		}
	}

	data := preview{
		Title:          title,
		Description:    description,
		BannerImageURL: bannerImageURL,
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(location+"/preview.html", buf.Bytes(), 0644)
}

func downloadImage(url, location string) (string, error) {
	extarr := strings.Split(url, ".")
	if len(extarr) == 0 {
		return "", fmt.Errorf("invalid url")
	}

	ext := extarr[len(extarr)-1]
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	imageLocation := fmt.Sprintf("%s/banner.%s", location, ext)

	file, err := os.Create(imageLocation)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return imageLocation, err
}
