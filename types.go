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

type StreamS struct {
	Name    string `json:"name,omitempty"`
	Quality string `json:"quality,omitempty"`
	URL     string `json:"url,omitempty"`
	Type    string `json:"type,omitempty"`
	Price   string `json:"price,omitempty"`
}

var _help_ = map[string]string{
	"google":  `Search Google.` + "\n" + `Usage: {}/google?q=<query>` + "\n" + `Example: {}/google?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Google Result]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"tpb":     `Search The Pirate Bay.` + "\n" + `Usage: {}/tpb?q=<query>` + "\n" + `Example: {}/tpb?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Torrents]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"imdb":    `Search IMDB.` + "\n" + `Usage: {}/imdb?q=<query>` + "\n" + `Example: /imdb?q=Avengers, {}/imdb?id=tt10191320` + "\n" + `Returns:` + "\n" + `    [Array of ImDB Titles or ImDB Single Result]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    id: The IMDB ID. (optional)` + "\n" + `    i: Indentaion (Bool, optional).`,
	"youtube": `Search YouTube.` + "\n" + `Usage: {}/youtube?q=<query>` + "\n" + `Example: {}/youtube?q=Avengers` + "\n" + `Returns:` + "\n" + `    [Array of Youtube Videos]` + "\n" + `If no results are found, an empty array is returned.` + "\n" + `Parameters:` + "\n" + `    q: The search query.` + "\n" + `    i: Indentaion (Bool, optional).`,
	"chatbot": `Talk to the Kuki chatbot.` + "\n" + `Usage: {}/chatbot?message=<query>` + "\n" + `Example: {}/chatbot?message=Hello` + "\n" + `Returns:` + "\n" + `    {"message":"response"}` + "\n" + `If internal server error, 502 is returned.` + "\n" + `Parameters:` + "\n" + `    message: The message query.`,
}
