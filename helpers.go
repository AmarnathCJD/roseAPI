package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	spotifyClient = map[string]string{
		"access_token": "",
		"expires_in":   "",
	}
	Aclient = NewAPIClient()
)

func NewAPIClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
	}
	c := &http.Client{Transport: tr}
	return c
}

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
	io.Copy(w, bytes.NewBuffer([]byte(fmt.Sprintf(`, "status": "ok", "ping": "%v"}`, totaltime))))
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

func GetSpotifyCred() string {
	if spotifyClient["access_token"] == "" {
		token, expire := GenSpotifyToken()
		spotifyClient["access_token"] = token
		spotifyClient["expire"] = expire
		return token
	} else {
		exp, _ := strconv.ParseInt(spotifyClient["expire"], 10, 64)
		if time.Now().Unix() > (exp / 1000) {
			token, expire := GenSpotifyToken()
			spotifyClient["access_token"] = token
			spotifyClient["expire"] = expire
			return token
		} else {
			return spotifyClient["access_token"]
		}
	}
}

func GenSpotifyToken() (string, string) {
	req, _ := http.NewRequest("GET", "https://open.spotify.com/", nil)
	req.Header.Set("cookie", "sp_t=db7961ffe77e0f185d4a6c0d2fa5c47a; sp_landing=https%3A%2F%2Fopen.spotify.com%2F%3Fsp_cid%3Ddb7961ffe77e0f185d4a6c0d2fa5c47a%26device%3Ddesktop; sp_dc=AQAUUpb67K1lWwc1699YYBH19NdNSCbWWjCSWTwnK-gFIy5ik30bSF4caXyUL_ZiwvEDQ8DnfMhMdWWme75KJisRIw08KEI7sJWWrhuuu0rO8EzWyByMqEAd38uZVikquXaOUBKpj2sEWSCE7Es_RcNxtkpJM3ZL; sp_key=52113a31-3729-4bb1-8346-8dfa5ecac746; OptanonAlertBoxClosed=2022-05-17T05:27:23.537Z; OptanonConsent=isIABGlobal=false&datestamp=Tue+May+17+2022+13%3A25%3A58+GMT%2B0530+(India+Standard+Time)&version=6.26.0&hosts=&landingPath=NotLandingPage&groups=s00%3A1%2Cf00%3A1%2Cm00%3A1%2Ct00%3A1%2Ci00%3A1%2Cf02%3A1%2Cm02%3A1%2Ct02%3A1&geolocation=IN%3BKL&AwaitingReconsent=false")
	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	accessTokenReg := regexp.MustCompile(`"accessToken":"(.+?)"`)
	expireTimeReg := regexp.MustCompile(`"accessTokenExpirationTimestampMs":(\d+)`)
	var body []byte
	body, _ = ioutil.ReadAll(resp.Body)
	accessToken := accessTokenReg.FindStringSubmatch(string(body))
	expireTime := expireTimeReg.FindStringSubmatch(string(body))
	if len(accessToken) > 1 {
		return accessToken[1], expireTime[1]
	} else {
		return "", ""
	}
}

func SearchSptfy(query string, accessToken string) SpotifyResult {
	req, _ := http.NewRequest("GET", `https://api-partner.spotify.com/pathfinder/v1/query?operationName=searchDesktop&variables=%7B%22searchTerm%22%3A%22`+url.QueryEscape(query)+`%22%2C%22offset%22%3A0%2C%22limit%22%3A10%2C%22numberOfTopResults%22%3A5%2C%22includeAudiobooks%22%3Afalse%7D&extensions=%7B%22persistedQuery%22%3A%7B%22version%22%3A1%2C%22sha256Hash%22%3A%2219967195df75ab8b51161b5ac4586eab9cf73b51b35a03010073533c33fd11ae%22%7D%7D`, nil)
	req.Header.Set("app-platform", "WebPlayer")
	req.Header.Set("authorization", "Bearer "+accessToken)
	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
	var result SpotifyResult
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func Addiotonal(urI string) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/tracks?ids=14jSfsXObzCumwn8wX9amf,3XYvdqcZrTmRntFDDbJkJd,5raWEZXYAapq6Qw1GIEIkU&market=from_token", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("authorization", "Bearer BQDbVhcmEEHvJXuZI3s98dCZEHdw3bCd6Lwj_fqKSTursZcT60emPmEVQr1uQXIrGEvtF3XkT37LXA79MfN754aExY-CV97Lp7m8ailXNIIqF9mYFotVWzl-hD-ePseJDZQUEHxOCJPSqxNBsg4DgbDBZhTAFx8xksCdildsBv6QEOvmDtG1DDBKtRYE7l9VG2qMae4a5cNuu5qizP31-tBUuOvav2qWOr3-rGareoA9AnqhjPL5VuRpquhyqIOWfWcen2CCj1tSrs-b-XohUEJRo-2nVm89p-UlEMHpyKpuqX2TIuIllWbY")
	req.Header.Set("Referer", "https://open.spotify.com/")
	req.Header.Set("DNT", "1")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	c.Do(req)
}

// spotify lyrics tricky

func FetchLy(query string) string {
	url := "https://api.musixmatch.com/ws/1.1/track.search?q_track=" + url.QueryEscape(query) + "&page_size=3&page=1&s_track_rating=desc&apikey=6efc39fa2ad207c07d3d814749804a9b"
	req, _ := http.NewRequest("GET", url, nil)
	r, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	mzk, _, _, _ := jsonparser.Get(body, "message", "body", "track_list")
	var result []string
	jsonparser.ArrayEach(mzk, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		res, _, _, _ := jsonparser.Get(value, "track", "track_id")
		result = append(result, string(res))
	})
	req_2, _ := http.NewRequest("GET", "https://api.musixmatch.com/ws/1.1/track.lyrics.get?track_id="+result[0]+"&apikey=6efc39fa2ad207c07d3d814749804a9b", nil)
	r_2, err := c.Do(req_2)
	if err != nil {
		log.Println(err)
	}
	defer r_2.Body.Close()
	body_2, _ := ioutil.ReadAll(r_2.Body)
	ly, _, _, _ := jsonparser.Get(body_2, "message", "body", "lyrics", "lyrics_body")
	_ly := string(ly)
	_ly = strings.Replace(_ly, "******* This Lyrics is NOT for Commercial use *******", "", -1)
	return _ly
}

func Ly3(q string) string {
	url := "https://gsearch-prod-cloud.gaana.com/gaanasearch-api/mobilesuggest/autosuggest-lite-vltr-ro?geoLocation=IN&query=" + url.QueryEscape(q) + "&content_filter=2&include=allItems&isRegSrch=0&webVersion=mix&rType=web&usrLang=Hindi,English,Punjabi&isChrome=1"
	resp, err := Aclient.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result []string
	b, _, _, _ := jsonparser.Get(body, "gr")
	jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		gd, _, _, _ := jsonparser.Get(value, "gd")
		jsonparser.ArrayEach(gd, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			seo, _ := jsonparser.GetString(value, "seo")
			result = append(result, seo)
		})
	})
	bse := "https://gaana.com/lyrics/" + result[0]
	r, err := c.Get(bse)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
        var Lyric string 
	doc, _ := goquery.NewDocumentFromReader(r.Body)
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("lyr_data") {
			Lyric = s.Text()
		}
	})
	return Lyric
}

// https://gaana.com/lyrics/coca-cola-38
