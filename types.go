package main

type GoogleResult struct {
	Title       string `json:"title,omitempty"`
	Url         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

type YoutubeResult struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Channel     string `json:"channel,omitempty"`
	Published   string `json:"published,omitempty"`
	Duration    string `json:"duration,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	Views       string `json:"views,omitempty"`
}

type GameResult struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	Poster      string `json:"poster,omitempty"`
	ID          string `json:"id,omitempty"`
}

type ImDBData struct {
	D []struct {
		I struct {
			Height   int    `json:"height"`
			ImageURL string `json:"imageUrl"`
			Width    int    `json:"width"`
		} `json:"i,omitempty"`
		ID   string `json:"id"`
		L    string `json:"l"`
		Q    string `json:"q"`
		Rank int    `json:"rank"`
		S    string `json:"s"`
		Vt   int    `json:"vt,omitempty"`
		Y    int    `json:"y"`
	} `json:"d"`
}

type ImDBResult struct {
	Title  string `json:"title,omitempty"`
	ID     string `json:"id,omitempty"`
	Year   string `json:"year,omitempty"`
	Actors string `json:"actors,omitempty"`
	Rank   string `json:"rank,omitempty"`
	Link   string `json:"link,omitempty"`
	Poster string `json:"poster,omitempty"`
}

type LyricsSearch struct {
	Response struct {
		Sections []struct {
			Hits []struct {
				Result struct {
					Type           string `json:"_type"`
					APIPath        string `json:"api_path"`
					ArtistNames    string `json:"artist_names"`
					FullTitle      string `json:"full_title"`
					HeaderImageURL string `json:"header_image_url"`
					ID             int    `json:"id"`
					Title          string `json:"title"`
					Path           string `json:"path"`
				} `json:"result"`
			} `json:"hits"`
		} `json:"sections"`
	} `json:"response"`
}

type SpotifyResult struct {
	Data struct {
		SearchV2 struct {
			Albums struct {
				TotalCount int `json:"totalCount"`
				Items      []struct {
					Data struct {
						Typename string `json:"__typename"`
						URI      string `json:"uri"`
						Name     string `json:"name"`
						Artists  struct {
							Items []struct {
								URI     string `json:"uri"`
								Profile struct {
									Name string `json:"name"`
								} `json:"profile"`
							} `json:"items"`
						} `json:"artists"`
						CoverArt struct {
							Sources []struct {
								URL string `json:"url"`
							} `json:"sources"`
						} `json:"coverArt"`
						Date struct {
							Year int `json:"year"`
						} `json:"date"`
					} `json:"data"`
				} `json:"items"`
			} `json:"albums"`
		} `json:"searchV2"`
	} `json:"data"`
}

type LyricsR struct {
	Lyrics struct {
		Lines []struct {
			StartTimeMs string `json:"startTimeMs"`
			Words       string `json:"words"`
		} `json:"lines"`
		Language string `json:"language"`
	} `json:"lyrics"`
}

type OTT struct {
	Name    string `json:"name,omitempty"`
	Quality string `json:"quality,omitempty"`
	URL     string `json:"url,omitempty"`
	Type    string `json:"type,omitempty"`
	Price   string `json:"price,omitempty"`
}

type Title struct {
	Title       string   `json:"title,omitempty"`
	Year        string   `json:"year,omitempty"`
	Rating      string   `json:"rating,omitempty"`
	Genre       []string `json:"genre,omitempty"`
	Plot        string   `json:"plot,omitempty"`
	Poster      string   `json:"poster,omitempty"`
	ID          string   `json:"id,omitempty"`
	Stars       []string `json:"stars,omitempty"`
	Directors   []string `json:"directors,omitempty"`
	AKA         string   `json:"aka,omitempty"`
	Production  string   `json:"production,omitempty"`
	Language    string   `json:"language,omitempty"`
	ReleaseDate string   `json:"releaseDate,omitempty"`
}

type FileExt struct {
	Ext         string                           `json:"ext,omitempty"`
	Title       string                           `json:"title,omitempty"`
	Description string                           `json:"description,omitempty"`
	Category    string                           `json:"category,omitempty"`
	Icon        string                           `json:"icon,omitempty"`
	Programs    []map[string][]map[string]string `json:"programs,omitempty"`
	Url         string                           `json:"url,omitempty"`
}

type Endpoint struct {
	Name    string      `json:"name,omitempty"`
	Path    string      `json:"path,omitempty"`
	Info    string      `json:"info,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  []Parameter `json:"params,omitempty"`
	Returns []Parameter `json:"return,omitempty"`
	Usage   string      `json:"usage,omitempty"`
}

type Parameter struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Info string `json:"info,omitempty"`
}

var _help_ = map[string]string{
	"google":         `Search Google.` + "\n" + `Usage: {}/google?q=<query>` + "\n" + `Example: {}/google?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Google Result]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"tpb":            `Search The Pirate Bay.` + "\n" + `Usage: {}/tpb?q=<query>` + "\n" + `Example: {}/tpb?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Torrents]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"imdb":           `Search IMDB.` + "\n" + `Usage: {}/imdb?q=<query>` + "\n" + `Example: /imdb?q=Avengers, {}/imdb?id=tt10191320` + "\n" + `Returns:` + "\n" + `    [Array of ImDB Titles or ImDB Single Result]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    id: The IMDB ID. (optional)` + "\n" + `    i: Indentaion (Bool, optional).`,
	"imdb/title":     `Search IMDB by Title.` + "\n" + `Usage: {}/imdb/title?id=<title>` + "\n" + `Example: {}/imdb/title?title=tt10191320` + "\n" + `Returns:` + "\n" + `    [Array of IMDB Titles]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    id: The imdb title ID.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"youtube":        `Search YouTube.` + "\n" + `Usage: {}/youtube?q=<query>` + "\n" + `Example: {}/youtube?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Youtube Videos]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"chatbot":        `Talk to the Kuki chatbot.` + "\n" + `Usage: {}/chatbot?message=<query>` + "\n" + `Example: {}/chatbot?message=Hello` + "\n" + `Returns:` + "\n" + `    {"message":"response"}` + "\n" + `If internal server error, 502 is returned.` + "\n" + `Parameters:` + "\n" + `    message: The message query.`,
	"spotify":        `Search Spotify.` + "\n" + `Usage: {}/spotify?q=<query>` + "\n" + `Example: {}/spotify?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Spotify Results]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"lyrics":         `Search Lyrics.` + "\n" + `Usage: {}/lyrics?q=<query>` + "\n" + `Example: {}/lyrics?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Lyrics Results]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"stream":         `Search Stream.` + "\n" + `Usage: {}/stream?q=<query>` + "\n" + `Example: {}/stream?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Stream Results]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"youtube/stream": `Get stream url for youtube Vids.` + "\n" + `Usage: {}/youtube/stream?id=<vid_id>` + "\n" + `Example: {}/youtube/stream?id=8FAUEv_E_xQ` + "\n" + `Returns:` + "\n" + `    [Array of Stream URLS with Avaliable Qualities]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    id: The Vid ID.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"fileinfo":       `Get file extension info.` + "\n" + `Usage: {}/fileinfo?ext=<ext>` + "\n" + `Example: {}/fileinfo?ext=eac3` + "\n" + `Returns:` + "\n" + `    [File Info]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    ext/q: The extension.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"screenshot":     `Get screenshot for a url.` + "\n" + `Usage: {}/screenshot?url=<url>` + "\n" + `Example: {}/screenshot?url=https://www.google.com` + "\n" + `Returns:` + "\n" + `    [Screenshot]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    url: The url.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"urlpreview":     `Get url preview.` + "\n" + `Usage: {}/urlpreview?url=<url>` + "\n" + `Example: {}/urlpreview?url=https://www.google.com` + "\n" + `Returns:` + "\n" + `    [Url Preview]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    url: The url.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"math":           `Calculate math.` + "\n" + `Usage: {}/math?q=<query>` + "\n" + `Example: {}/math?q=1+1` + "\n" + `Returns:` + "\n" + `    [Math Result]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The math query.` + "\n" + `    i: Indentaion (Bool, optional).`,
}

var Endpoints = []Endpoint{
	{
		Name:   "IMDB Title",
		Path:   "/imdb/title",
		Method: "GET",
		Info:   "Search the internet movie database by the title ID (IMDB ID).",
		Usage:  `{}/imdb/title?id=tt0111161`,
		Params: []Parameter{
			{Name: "id", Type: "string", Info: "The IMDB ID."},
			{Name: "title", Type: "string", Info: "The Title ID."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "Title", Type: "JSON Result", Info: "The title."},
		},
	},
	{
		Name:   "IMDB Title Search",
		Path:   "/imdb",
		Method: "GET",
		Info:   "Search the internet movie database by the title.",
		Usage:  `{}/imdb?q=avengers`,
		Params: []Parameter{
			{Name: "q", Type: "string", Info: "The title."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "List of Titles", Type: "JSON Result", Info: "The title."},
		},
	},
	{
		Name:   "Google Search",
		Path:   "/google",
		Method: "GET",
		Info:   "Search the internet using Google.",
		Usage:  `{}/google?q=avengers`,
		Params: []Parameter{
			{Name: "q", Type: "string", Info: "The query."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "List of Results", Type: "JSON Result", Info: "results."},
		},
	},
	{
		Name:   "Google Image Search",
		Path:   "/google/image",
		Method: "GET",
		Info:   "Search the internet for images using Google.",
		Usage:  `{}/google/image?q=avengers`,
		Params: []Parameter{
			{Name: "q", Type: "string", Info: "The query."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "List of Results", Type: "JSON Result", Info: "Image results."},
		},
	},
	{
		Name:   "Youtube Search",
		Path:   "/youtube",
		Method: "GET",
		Info:   "Search the internet for videos using Youtube.",
		Usage:  `{}/youtube?q=avengers`,
		Params: []Parameter{
			{Name: "q", Type: "string", Info: "The query."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "List of Results", Type: "JSON Result", Info: "Video results."},
		},
	},
	{
		Name:   "Youtube Video",
		Path:   "/youtube/video",
		Method: "GET",
		Info:   "Get video information using Youtube.",
		Usage:  `{}/youtube/video?id=<id>`,
		Params: []Parameter{
			{Name: "id", Type: "string", Info: "The video id."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "Video", Type: "JSON Result", Info: "Video result."},
		},
	},
	{
		Name:   "Spotify Search",
		Path:   "/spotify",
		Method: "GET",
		Info:   "Search for music using Spotify.",
		Usage:  `{}/spotify?q=avengers`,
		Params: []Parameter{
			{Name: "q", Type: "string", Info: "The query."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "List of Results", Type: "JSON Result", Info: "Spotify results."},
		},
	},
	{
		Name:   "Spotify Track",
		Path:   "/spotify/track",
		Method: "GET",
		Info:   "Get track information using Spotify.",
		Usage:  `{}/spotify/track?id=<id>`,
		Params: []Parameter{
			{Name: "id", Type: "string", Info: "The track id."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "Track", Type: "JSON Result", Info: "Track result."},
		},
	},
	{
		Name:   "Lyric Search",
		Path:   "/lyrics",
		Method: "GET",
		Info:   "Search for lyrics using Lyrics.com.",
		Usage:  `{}/lyrics?q=Vaaste`,
		Params: []Parameter{
			{Name: "q", Type: "string", Info: "The query."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "List of Results", Type: "JSON Result", Info: "Lyrics results."},
		},
	},
	{
		Name:   "File Extension Info",
		Path:   "/fileinfo",
		Method: "GET",
		Info:   "Get information about a file extension.",
		Usage:  `{}/fileinfo?ext=mp3`,
		Params: []Parameter{
			{Name: "ext", Type: "string", Info: "The file extension."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "File Extension Info", Type: "JSON Result", Info: "File extension info."},
		},
	},
	{
		Name:   "KUKI Chatbot",
		Path:   "/chatbot",
		Method: "GET",
		Info:   "Chat with KUKI.",
		Usage:  `{}/chatbot?message=Hello`,
		Params: []Parameter{
			{Name: "message", Type: "string", Info: "The message."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "KUKI Response", Type: "JSON Result", Info: "KUKI response."},
		},
	},
	{
		Name:   "Math",
		Path:   "/math",
		Method: "GET",
		Info:   "Do math.",
		Usage:  `{}/math?expression=1+1`,
		Params: []Parameter{
			{Name: "expression", Type: "string", Info: "The expression."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "Result", Type: "JSON Result", Info: "Result."},
		},
	},
	{
		Name:   "Url Scanner",
		Path:   "/urlpreview",
		Method: "GET",
		Info:   "Get preview information about a URL.",
		Usage:  `{}/urlpreview?url=https://www.google.com`,
		Params: []Parameter{
			{Name: "url", Type: "string", Info: "The URL."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "Preview", Type: "JSON Result", Info: "Preview."},
		},
	},
	{
		Name:   "Url Shortener",
		Path:   "/urlshortener",
		Method: "GET",
		Info:   "Shorten a URL.",
		Usage:  `{}/urlshortener?url=https://www.google.com`,
		Params: []Parameter{
			{Name: "url", Type: "string", Info: "The URL."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "Shortened URL", Type: "JSON Result", Info: "Shortened URL."},
		},
	},
	{
		Name:   "URL Screenshot",
		Path:   "/screenshot",
		Method: "GET",
		Info:   "Get a screenshot of a URL.",
		Usage:  `{}/screenshot?url=https://www.google.com`,
		Params: []Parameter{
			{Name: "url", Type: "string", Info: "The URL."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "Screenshot", Type: "JSON Result", Info: "Screenshot."},
		},
	},
	{
		Name:   "The Pirate Bay",
		Path:   "/tpb",
		Method: "GET",
		Info:   "Search for torrents on The Pirate Bay.",
		Usage:  `{}/tpb?q=avengers`,
		Params: []Parameter{
			{Name: "q", Type: "string", Info: "The query."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "List of Results", Type: "JSON Result", Info: "Torrent results."},
		},
	},
	{
		Name:   "OTT Streaming",
		Path:   "/stream",
		Method: "GET",
		Info:   "Stream a video from OTT.",
		Usage:  `{}/stream?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ`,
		Params: []Parameter{
			{Name: "url", Type: "string", Info: "The URL."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "Stream", Type: "JSON Result", Info: "Stream."},
		},
	},
	{
		Name:   "Youtube Download",
		Path:   "/youtube/download",
		Method: "GET",
		Info:   "Download a video from Youtube.",
		Usage:  `{}/youtube/download?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ, {}/youtube/download?id=dQw4w9WgXcQ&video=true, {}/youtube/download?q=Vaaste`,
		Params: []Parameter{
			{Name: "url", Type: "string", Info: "The URL."},
			{Name: "id", Type: "string", Info: "The video ID."},
			{Name: "video", Type: "bool", Info: "Download video (Bool, optional)., By default, audio is downloaded."},
			{Name: "q", Type: "string", Info: "The query."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
			{Name: "query", Type: "string", Info: "The query."},
			{Name: "download", Type: "bool", Info: "Selection between Download or Stream (Bool, optional)., By default, Streaming is selected."},
		},
		Returns: []Parameter{
			{Name: "video/audio", Type: "Bytes", Info: "Video/Audio."},
		},
	},
	{
		Name:   "Netflix Search",
		Path:   "/netflix/search",
		Method: "GET",
		Info:   "Search for a movie on Netflix.",
		Usage:  `{}/netflix/search?q=The%20Lion%20King`,
		Params: []Parameter{
			{Name: "q", Type: "string", Info: "The query."},
			{Name: "i", Type: "bool", Info: "Indentaion (Bool, optional)."},
		},
		Returns: []Parameter{
			{Name: "List of Results", Type: "JSON Result", Info: "Movie results."},
		},
	},
}

func GetEnpointByPath(path string) *Endpoint {
	for _, endpoint := range Endpoints {
		if endpoint.Path == path {
			return &endpoint
		}
	}
	return nil
}
