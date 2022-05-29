package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	c = &http.Client{}
)

func Tpb(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/tpb", w)
	}
	i := query.Get("i")
	q := query.Get("q")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	url := "https://tpb23.ukpass.co/apibay/q.php" + "?q=" + url.QueryEscape(q)
	resp, err := c.Get(url)
	if !ERR(err, w) {
		return
	}
	WriteJson(w, r, resp.Body, i)
}

func Google(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/google", w)
	}
	q := query.Get("q")
	i := query.Get("i")
	if query.Get("q") == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	URL := "https://www.google.com/search?q=" + url.QueryEscape(q)
	resp, err := c.Get(URL)
	if !ERR(err, w) {
		return
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	var results []GoogleResult
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("ZINbbc luh4tb xpd O9g5cc uUPGi") {
			var result GoogleResult
			s.Find("div").Each(func(i int, s *goquery.Selection) {
				if s.HasClass("egMi0 kCrYT") {
					result.Title = s.Text()
					url := strings.Split(strings.Replace(s.Find("a").AttrOr("href", ""), "/url?q=", "", 1), "&")[0]
					result.Url = url
				} else if s.HasClass("BNeawe s3v9rd AP7Wnd") {
					result.Description = s.Text()
				}
			})
			results = append(results, result)
		}
	})
	if !ERR(err, w) {
		return
	}
	var data string
	if len(results) != 0 {
		data = EncodeJson(results)
	} else {
		data = `[]`
	}
	WriteJson(w, r, data, i)
}

func Youtube(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/youtube", w)
	}
	q := query.Get("q")
	i := query.Get("i")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	URL := "https://www.youtube.com/results?search_query=" + url.QueryEscape(q)
	resp, err := c.Get(URL)
	if !ERR(err, w) {
		return
	}
	var exp, _ = regexp.Compile(`ytInitialData = [\s\S]*]`)
	b, _ := ioutil.ReadAll(resp.Body)
	match := exp.FindStringSubmatch(string(b))
	var d string
	if len(match) != 0 {
		d = match[0]
		d = strings.Replace(d, "ytInitialData = ", "", 1)
		d = strings.Split(d, ";</script>")[0]
	}
	pData := ParseYoutubeRAW(d)
	WriteJson(w, r, pData, i)
}

func ImDB(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/imdb", w)
	}
	q := query.Get("q")
	i := query.Get("i")
	id := query.Get("id")
	if q == "" && id == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	if q != "" {
		firstLetter := strings.ToLower(string(q[0]))
		URL := "https://v2.sg.media-imdb.com/suggestion/titles/" + firstLetter + "/" + url.QueryEscape(q) + ".json"
		resp, err := c.Get(URL)
		if !ERR(err, w) {
			return
		}
		var data ImDBData
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		var results []ImDBResult
		for _, r := range data.D {
			results = append(results, ImDBResult{Title: r.L, Year: fmt.Sprint(r.Y), ID: r.ID, Actors: r.S, Rank: fmt.Sprint(r.Rank), Link: "https://www.imdb.com/title/" + r.ID, Poster: r.I.ImageURL})
		}
		if !ERR(err, w) {
			return
		}
		var result string
		if len(results) != 0 {
			result = EncodeJson(results)
		} else {
			result = `[]`
		}
		WriteJson(w, r, result, i)
	} else if id != "" {
		w.Write([]byte(`{"title":"` + id + `"}`))
	}
}

func ChatBot(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	var API = "https://icap.iconiq.ai/talk?&botkey=icH-VVd4uNBhjUid30-xM9QhnvAaVS3wVKA3L8w2mmspQ-hoUB3ZK153sEG3MX-Z8bKchASVLAo~&channel=7&sessionid=482070240&client_name=uuiprod-un18e6d73c-user-19422&id=true"
	query := r.URL.Query()
	q := query.Get("message")
	if query.Get("help") != "" {
		WriteHelp("/chatbot", w)
	}
	if q == "" {
		http.Error(w, "missing 'message'", http.StatusBadRequest)
		return
	}
	req, err := http.PostForm(API, url.Values{"input": {q}})
	if !ERR(err, w) {
		return
	}
	defer req.Body.Close()
	var resp map[string]interface{}
	json.NewDecoder(req.Body).Decode(&resp)
	msg := resp["responses"].([]interface{})[0].(string)
	d := `{"message": "` + msg + `"}`
	WriteJson(w, r, d, "")

}

func Math(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/math", w)
	}
	q := query.Get("q")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	url := "https://evaluate-expression.p.rapidapi.com"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-host", "evaluate-expression.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "cf9e67ea99mshecc7e1ddb8e93d1p1b9e04jsn3f1bb9103c3f")
	_query := req.URL.Query()
	_query.Add("expression", q)
	req.URL.RawQuery = _query.Encode()
	resp, err := c.Do(req)
	if !ERR(err, w) {
		return
	}
	defer resp.Body.Close()
	var b interface{}
	json.NewDecoder(resp.Body).Decode(&b)
	if b == nil {
		WriteJson(w, r, []byte(` "invalid mathematical expression"`), "")
		return
	}
	WriteJson(w, r, []byte(`{"expression": "`+q+`", "result": "`+fmt.Sprint(b)+`"}`), "")
}

func Pinterest(w http.ResponseWriter, r *http.Request) {
}

func Spotify(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/spotify", w)
	}
	q := query.Get("q")
	i := query.Get("i")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	t := GetSpotifyCred()
	if t == "" {
		http.Error(w, "missing spotify credentials", http.StatusBadRequest)
		return
	}
	s := SearchSptfy(q, t)
	var data = s.Data.SearchV2.Albums.Items
	d, _ := json.Marshal(data)
	WriteJson(w, r, string(d), i)
}

func LyricsA(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/lyrics", w)
	}
	q := query.Get("q")
	i := query.Get("i")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	ly := Ly3(q)
	ly = strings.TrimSpace(ly)
	var ly_list []string
	for _, a := range strings.Split(ly, "\n") {
		if a != "" {
			ly_list = append(ly_list, a)
		}
	}
	_ily := strings.Join(ly_list, "\n")
	WriteJson(w, r, EncodeJson(`"`+_ily+`"`), i)
}

func Stream(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/stream", w)
	}
	q := query.Get("q")
	i := query.Get("i")
	if q == "" {
		http.Error(w, "missing query", http.StatusBadRequest)
		return
	}
	d := StreamSrc(q)
	WriteJson(w, r, EncodeJson(d), i)
}

func YoutubeStream(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/youtube/stream", w)
	}
	id := query.Get("id")
	i := query.Get("i")
	if id == "" {
		WriteError("missing id", w)
		return
	}
	vid_url := "https://www.youtube.com/watch?v=" + id
	var data = strings.NewReader(`{"url":"` + vid_url + `"}`)
	req, _ := http.NewRequest("POST", "https://api.onlinevideoconverter.pro/api/convert", data)
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if !ERR(err, w) {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if !ERR(err, w) {
		return
	}
	WriteJson(w, r, string(_UnescapeUnicodeCharactersInJSON(body)), i)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("index.html"))
	template.Execute(w, nil)
}

func LinkPreview(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/urlpreview", w)
	}
	_url := query.Get("url")
	i := query.Get("i")
	if _url == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}
	req, _ := http.NewRequest("GET", "https://api.labs.cognitive.microsoft.com/urlpreview/v7.0/search"+"?q="+url.QueryEscape(_url), nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", "27b02a2c7d394388a719e0fdad6edb10")
	resp, err := c.Do(req)
	if !ERR(err, w) {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	WriteJson(w, r, string(body), i)
}

func ScreenShot(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/screenshot", w)
	}
	_url := query.Get("url")
	i := query.Get("i")
	image := query.Get("image")
	if _url == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}
	BASEURL := fmt.Sprintf("https://webshot.deam.io/%s?type=jpeg&quality=100&fullPage=false&height=540&width=960", _url)
	resp, err := c.Get(BASEURL)
	if !ERR(err, w) {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), "<h") {
		WriteJson(w, r, string([]byte(`{"error":"`+string(body)+`"}`)), i)
		return
	}
	if image == "true" {
		w.Header().Set("Content-Type", "image/png")
		w.Write(body)
		return
	}
	sEnc := b64.StdEncoding.EncodeToString(body)
	WriteJson(w, r, string([]byte(`{"image":"`+sEnc+`"}`)), i)
}

func OCR(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "POST") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		w.Write([]byte(strings.ReplaceAll(_help_["imdb"], "{}", r.URL.Hostname())))
		return
	}
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := ioutil.ReadAll(file)
	HEADERS := map[string]string{
		"X-Api-Key": "IQcdz030YPMT3zSRrhHzRQ==sNdD9akTySL4WcpS",
	}
	req := newfileUploadRequest("https://api.api-ninjas.com/v1/imagetotext", map[string]string{}, "image", b, HEADERS)
	resp, err := c.Do(req)
	if !ERR(err, w) {
		return
	}
	defer resp.Body.Close()
	bd, _ := ioutil.ReadAll(resp.Body)
	w.Write(bd)
}

func FileInfo(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/fileinfo", w)
	}
	q := query.Get("q")
	i := query.Get("i")
	if q == "" {
		q = query.Get("ext")
	}
	if q == "" {
		WriteError("missing param, 'ext' or 'q'", w)
		return
	}
	URL := "https://fileinfo.com/extension/" + url.QueryEscape(q)
	resp, err := c.Get(URL)
	if !ERR(err, w) {
		return
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	ext := FileExt{
		Ext:         q,
		Title:       doc.Find("title").Text(),
		Description: doc.Find(".infoBox").Text(),
		Url:         URL,
		Icon:        doc.Find(".entryIcon").AttrOr("data-bg-lg", ""),
	}
	var Programs []map[string][]map[string]string
	var pt []map[string]string
	doc.Find(".programs").Each(func(i int, s *goquery.Selection) {
		platform := s.AttrOr("data-plat", "")
		if platform != "" {
			s.Find(".appmid").Each(func(i int, s *goquery.Selection) {
				pt = append(pt, map[string]string{
					"name":    s.Find("a").Text(),
					"url":     s.Find("a").AttrOr("href", ""),
					"license": s.Find(".license").Text(),
				})
			})
			Programs = append(Programs, map[string][]map[string]string{platform: pt})
			pt = []map[string]string{}
		}
	})
	ext.Programs = Programs
	WriteJson(w, r, ext, i)
}

func ImdbTitleInfo(w http.ResponseWriter, r *http.Request) {
	if !blockWrongMethod(w, r, "GET") {
		return
	}
	r.Header.Set("X-Start-Time", fmt.Sprint(time.Now().UnixNano()))
	query := r.URL.Query()
	if query.Get("help") != "" {
		WriteHelp("/imdb/title", w)
	}
	id := query.Get("id")
	if id == "" {
		id = query.Get("title")
	}
	if id == "" {
		WriteError("missing param, 'id' or 'title'", w)
		return
	}
	T := ImdbTtitle(id)
	if T.Title == "" {
		WriteError("not found in ImDB", w)
		return
	}
	WriteJson(w, r, T, query.Get("i"))
}

func init() {
	http.HandleFunc("/tpb", Tpb)
	http.HandleFunc("/google", Google)
	http.HandleFunc("/youtube", Youtube)
	http.HandleFunc("/youtube/stream", YoutubeStream)
	http.HandleFunc("/imdb", ImDB)
	http.HandleFunc("/imdb/title", ImdbTitleInfo)
	http.HandleFunc("/chatbot", ChatBot)
	http.HandleFunc("/lyrics", LyricsA)
	http.HandleFunc("/math", Math)
	http.HandleFunc("/spotify", Spotify)
	http.HandleFunc("/stream", Stream)
	http.HandleFunc("/urlpreview", LinkPreview)
	http.HandleFunc("/screenshot", ScreenShot)
	http.HandleFunc("/ocr", OCR)
	http.HandleFunc("/fileinfo", FileInfo)
	http.HandleFunc("/", HomePage)
}
