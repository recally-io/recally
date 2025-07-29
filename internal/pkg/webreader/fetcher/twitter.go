package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// https://cdn.syndication.twimg.com/tweet-result?id=1882066129990135874&token=4k8
type Tweet struct {
	ReplyToScreenName  string `json:"in_reply_to_screen_name"`
	ReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
	ReplyToUserIDStr   string `json:"in_reply_to_user_id_str"`

	Lang              string    `json:"lang"`
	ReplyCount        int       `json:"reply_count"`
	RetweetCount      int       `json:"retweet_count"`
	FavoriteCount     int       `json:"favorite_count"`
	PossiblySensitive bool      `json:"possibly_sensitive"`
	CreatedAt         time.Time `json:"created_at"`
	Entities          Entities  `json:"entities"`
	IDStr             string    `json:"id_str"`
	Text              string    `json:"text"`
	User              User      `json:"user"`
	Medias            []Media   `json:"mediaDetails"`
}

type Media struct {
	Type string `json:"type"`
	URL  string `json:"media_url_https"`
}

type User struct {
	IDStr           string `json:"id_str"`
	Name            string `json:"name"`
	ProfileImageURL string `json:"profile_image_url_https"`
	ScreenName      string `json:"screen_name"`
}

type Entities struct {
	Hashtags []struct {
		Text string `json:"text"`
	} `json:"hashtags"`
	URLs         []EntityURL   `json:"urls"`
	UserMentions []User        `json:"user_mentions"`
	Symbols      []interface{} `json:"symbols"`
	Media        []EntityURL   `json:"media"`
}

type EntityURL struct {
	URL         string `json:"url"`
	OriginalURL string `json:"expanded_url"`
}

const tweetTemplate = `
### [Tweet](https://x.com/{{.User.ScreenName}}/status/{{.IDStr}}) by [{{.User.Name}}](https://x.com/{{.User.ScreenName}})

**Posted:** {{.CreatedAt.Format "2006-01-02 15:04:05"}}

{{.Text}}

{{if .Entities.Hashtags}}**Tags:** {{range .Entities.Hashtags}}#{{.Text}} {{end}}{{end}}
{{if .Entities.UserMentions}}**Mentions:** {{range .Entities.UserMentions}}@[{{.ScreenName}}](https://x.com/{{.ScreenName}}) {{end}}{{end}}

{{if .Entities.URLs}}**Links:**
{{range .Entities.URLs}}- {{.OriginalURL}}
{{end}}{{end}}

{{if .Medias}}**Media:**
{{range .Medias}}- ![]({{.URL}})
{{end}}{{end}}

---
`

var tweetTmpl = template.Must(template.New("tweet").Parse(tweetTemplate))

type TwitterFetcher struct {
	client *http.Client
}

func NewTwitterFetcher() *TwitterFetcher {
	return &TwitterFetcher{
		client: &http.Client{
			Timeout: time.Duration(30) * time.Second,
		},
	}
}

func (f *TwitterFetcher) Fetch(ctx context.Context, uri string) (*webreader.FetchedContent, error) {
	// https://x.com/FDavisv_/status/1881604238356439483
	// twitterId := "1881604238356439483"
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	paths := strings.Split(u.Path, "/")
	if len(paths) < 3 || paths[2] != "status" {
		return nil, fmt.Errorf("invalid url: %s", uri)
	}
	twitterId := paths[3]
	tweet, err := f.fetchTweet(ctx, twitterId)
	if err != nil {
		return nil, fmt.Errorf("error to fetch tweet %s :%w", twitterId, err)
	}

	tweets := []*Tweet{tweet}
	for tweet.ReplyToStatusIDStr != "" {

		tweet, err = f.fetchTweet(ctx, tweet.ReplyToStatusIDStr)
		if err != nil {
			logger.FromContext(ctx).Error("fetch tweet", "err", err, "id", tweet.ReplyToStatusIDStr)
			return nil, fmt.Errorf("fetch tweet: %w", err)
		}
		tweets = append(tweets, tweet)
	}

	sb := &strings.Builder{}

	// from last tweet to first tweet construct the tweet content and join with new line
	for i := len(tweets) - 1; i >= 0; i-- {
		t, err := f.constructTweet(tweets[i])
		if err != nil {
			logger.FromContext(ctx).Error("construct tweet", "err", err, "tweet", tweets[i])
			return nil, fmt.Errorf("construct tweet: %w", err)
		}
		sb.WriteString("\n")
		sb.WriteString(t)
	}

	fetchedContent := &webreader.FetchedContent{
		Content: webreader.Content{
			Markwdown:     sb.String(),
			URL:           uri,
			Favicon:       tweet.User.ProfileImageURL,
			Title:         tweet.Text[:min(50, len(tweet.Text))],
			SiteName:      "Twitter",
			PublishedTime: &tweet.CreatedAt,
			Author:        tweet.User.Name,
		},
	}
	if tweet.Entities.Media != nil {
		fetchedContent.Image = tweet.Entities.Media[0].OriginalURL
	}

	return fetchedContent, nil
}

func (f *TwitterFetcher) fetchTweet(ctx context.Context, tweetId string) (*Tweet, error) {
	token := f.getToken(tweetId)
	cdnUrl := fmt.Sprintf("https://cdn.syndication.twimg.com/tweet-result?id=%s&token=%s", tweetId, token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cdnUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch tweet: %w", err)
	}
	defer resp.Body.Close()
	var tweet Tweet
	if err := json.NewDecoder(resp.Body).Decode(&tweet); err != nil {
		return nil, fmt.Errorf("decode tweet: %w", err)
	}
	return &tweet, nil
}

func (f *TwitterFetcher) getToken(id string) string {
	idNum, _ := strconv.ParseFloat(id, 64)
	token := (idNum / 1e15) * math.Pi
	base36 := strconv.FormatInt(int64(token), 36)
	re := regexp.MustCompile(`(0+|\.)`)
	return re.ReplaceAllString(base36, "")
}

func (f *TwitterFetcher) constructTweet(tweet *Tweet) (string, error) {
	sb := &strings.Builder{}
	if err := tweetTmpl.Execute(sb, *tweet); err != nil {
		return "", fmt.Errorf("execute tweet template: %w", err)
	}
	return sb.String(), nil
}

func (f *TwitterFetcher) Close() error {
	return nil
}
