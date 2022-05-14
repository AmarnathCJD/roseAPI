package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

func fetchPort() string {
	var p = os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":3000"
}

func ERR(err error, w http.ResponseWriter) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	return true
}

func blockWrongMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func WriteJson(w http.ResponseWriter, r *http.Request, data interface{}, i string) {
	var datav []byte
	switch data := data.(type) {
	case io.ReadCloser:
		datav, _ = ioutil.ReadAll(data)
	case []byte:
		datav = data
	case string:
		datav = []byte(data)
	default:
		datav, _ = json.Marshal(data)
	}
	var bf bytes.Buffer
	if i == "true" {
		json.Indent(&bf, datav, "", "  ")
	} else {
		bf = *bytes.NewBuffer(datav)
	}
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewBuffer([]byte("{\"data\":")))
	io.Copy(w, &bf)
	startTime, _ := strconv.ParseInt(r.Header.Get("X-Start-Time"), 10, 64)
	CalcPing(w, startTime)
}

func CalcPing(w http.ResponseWriter, startTime int64) {
	w.Header().Set("Content-Type", "application/json")
	totaltime := time.Since(time.Unix(0, startTime))
	io.Copy(w, bytes.NewBuffer([]byte(fmt.Sprintf(`}, {"status": "ok", "ping": "%v"}}`, totaltime))))
}

func EncodeJson(v interface{}) string {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetEscapeHTML(false)
	encoder.Encode(v)
	return b.String()
}

func ParseYoutubeRAW(raw string) []byte {
	by := []byte(raw)
	var Results []YoutubeResult
	a, _, _, _ := jsonparser.Get(by, "contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents")
	jsonparser.ArrayEach(a, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		b, _, _, _ := jsonparser.Get(value, "itemSectionRenderer", "contents")
		jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			c, _, _, _ := jsonparser.Get(value, "videoRenderer")
			d, _, _, _ := jsonparser.Get(c, "title", "runs")
			if d != nil {
				var Result YoutubeResult
				jsonparser.ArrayEach(d, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					text, _, _, _ := jsonparser.Get(value, "text")
					Result.Title = string(text)
				})
				e, _, _, _ := jsonparser.Get(c, "thumbnail", "thumbnails", "[0]", "url")
				if e != nil {
					Result.Thumbnail = string(e)
				}
				f, _, _, _ := jsonparser.Get(c, "videoId")
				if f != nil {
					Result.URL = "https://www.youtube.com/watch?v=" + string(f)
				}
				metadata, _, _, _ := jsonparser.Get(c, "detailedMetadataSnippets")
				if metadata != nil {
					jsonparser.ArrayEach(metadata, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
						g, _, _, _ := jsonparser.Get(value, "snippetText", "runs")
						if g != nil {
							var desc string
							jsonparser.ArrayEach(g, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
								text, _, _, _ := jsonparser.Get(value, "text")
								desc += string(text)
							})
							Result.Description = desc
						}
					})
				}
				ownerText, _, _, _ := jsonparser.Get(c, "ownerText", "runs")
				if ownerText != nil {
					jsonparser.ArrayEach(ownerText, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
						text, _, _, _ := jsonparser.Get(value, "text")
						Result.Channel = string(text)
					})
				}
				videoID, _, _, _ := jsonparser.Get(c, "videoId")
				if videoID != nil {
					Result.ID = string(videoID)
				}
				published, _, _, _ := jsonparser.Get(c, "publishedTimeText", "simpleText")
				if published != nil {
					Result.Published = string(published)
				}
				length, _, _, _ := jsonparser.Get(c, "lengthText", "simpleText")
				if length != nil {
					Result.Duration = string(length)
				}
				views, _, _, _ := jsonparser.Get(c, "viewCountText", "simpleText")
				if views != nil {
					Result.Views = string(views)
				}
				Results = append(Results, Result)
			}
		})
	})
	data, _ := json.Marshal(Results)
	return data
}
