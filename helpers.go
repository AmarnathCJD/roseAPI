package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
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

func _UnescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) json.RawMessage {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil
	}
	return []byte(str)
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

func StreamSrc(query string) []StreamS {
	var data = strings.NewReader(`{"operationName":"GetSearchTitles","variables":{"searchTitlesSortBy":"POPULAR","first":5,"sortRandomSeed":0,"searchAfterCursor":"","searchTitlesFilter":{"searchQuery":"` + query + `","personId":null},"language":"en","country":"IN"},"query":"query GetSearchTitles($country: Country!, $searchTitlesFilter: TitleFilter, $searchAfterCursor: String, $searchTitlesSortBy: PopularTitlesSorting! = POPULAR, $first: Int! = 5, $language: Language!, $sortRandomSeed: Int! = 0, $profile: PosterProfile, $backdropProfile: BackdropProfile, $format: ImageFormat) {\n  popularTitles(\n    country: $country\n    filter: $searchTitlesFilter\n    after: $searchAfterCursor\n    sortBy: $searchTitlesSortBy\n    first: $first\n    sortRandomSeed: $sortRandomSeed\n  ) {\n    totalCount\n    pageInfo {\n      startCursor\n      endCursor\n      hasPreviousPage\n      hasNextPage\n      __typename\n    }\n    edges {\n      ...SearchTitleGraphql\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment SearchTitleGraphql on PopularTitlesEdge {\n  cursor\n  node {\n    id\n    objectId\n    objectType\n    content(country: $country, language: $language) {\n      title\n      fullPath\n      originalReleaseYear\n      scoring {\n        imdbScore\n        imdbVotes\n        tmdbScore\n        tmdbPopularity\n        __typename\n      }\n      posterUrl(profile: $profile, format: $format)\n      backdrops(profile: $backdropProfile, format: $format) {\n        backdropUrl\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n"}`)
	req, _ := http.NewRequest("POST", "https://apis.justwatch.com/graphql", data)
	req.Header.Set("Content-Type", "application/json")
	r, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var pid string
	program_id, _, _, _ := jsonparser.Get(body, "data", "popularTitles", "edges", "[0]", "node", "content", "fullPath")
	pid = string(program_id)
	baseUrl := "https://www.justwatch.com" + pid
	rs, err := c.Get(baseUrl)
	if err != nil {
		log.Println(err)
	}
	defer rs.Body.Close()
	var src []StreamS
	doc, _ := goquery.NewDocumentFromReader(rs.Body)
	doc.Find(".price-comparison__grid__row__element").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a")
		img := a.Find("img")
		src = append(src, StreamS{img.AttrOr("alt", ""), strings.TrimSpace(s.Find("span").Text()), a.AttrOr("href", ""), "Stream", strings.TrimSpace(s.Find(".price-comparison__grid__row__price").Text())})
	})
	return src
}

func WriteError(msg string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	_d := []byte(`{"error":"` + msg + `", "status":400}`)
	w.Write(_d)
}
