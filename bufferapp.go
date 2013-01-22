package bufferapi

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Client struct {
	AccessToken string
	transport   *oauth.Transport
}

type Profile struct {
	Id                string         `json:"id"`
	UserId            string         `json:"user_id"`
	Avatar            string         `json:"avatar"`
	CreatedAt         int            `json:"created_at"`
	Default           bool           `json:"default"`
	FormattedUsername string         `json:"formatted_username"`
	Schedules         []Schedule     `json:"schedules"`
	Service           string         `json:service"`
	ServiceId         string         `json:"service_id"`
	ServiceUsername   string         `json:"service_username"`
	Statistics        map[string]int `json:"statistics"`
	TeamMembers       []string       `json:"team_members"`
	Timezone          string         `json:"timezone"`
}

type Schedule struct {
	Days  []string `json:"days"`
	Times []string `json:"times"`
}

type NewUpdate struct {
	Text       string            `json:"text"`
	ProfileIds []string          `json:"profile_ids"`
	Shorten    bool              `json:"shorten"`
	Now        bool              `json:"now"`
	Media      map[string]string `json:"media"`
}

type Update struct {
	Id             string            `json:"id"`
	CreatedAt      int               `json:"created_at"`
	Day            string            `json:"day"`
	DueAt          int               `json:"due_at"`
	DueTime        string            `json:"due_time"`
	media          map[string]string `json:"media"`
	ProfileId      string            `json:"profile_id"`
	ProfileService string            `json:"profile_service"`
	Status         string            `json:"status"`
	Text           string            `json:"text"`
	TextFormatted  string            `json:"text_formatted"`
	UserId         string            `json:"user_id"`
	Via            string            `json:"via"`
}

type UpdateResponse struct {
	Success          bool     `json:"success"`
	BufferCount      int      `json:"buffer_count"`
	BufferPercentage int      `json:"buffer_percentage"`
	Updates          []Update `json:"updates"`
}

func ClientFactory(token, clientId, clientSecret, scope, authUrl, tokenUrl, cacheFile string) *Client {
	config := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        scope,
		AuthURL:      authUrl,
		TokenURL:     tokenUrl,
		TokenCache:   oauth.CacheFile(cacheFile),
	}
	transport := &oauth.Transport{Config: config}
	t := &oauth.Token{AccessToken: token}
	transport.Token = t
	c := Client{AccessToken: token, transport: transport}
	return &c
}

func (c *Client) API(method, url string, data interface{}) (respBody []byte, err error) {
	jsonPattern, _ := regexp.Compile(`\.json$`)
	if !jsonPattern.Match([]byte(url)) {
		url += ".json"
	}
	url = "https://api.bufferapp.com/1/" + url

	jsonBody, _ := json.Marshal(data)
	b := bytes.NewBuffer(jsonBody)

	var resp *http.Response
	switch method {
	case "get":
		resp, err = c.transport.Client().Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	case "post":
		resp, err := c.transport.Client().Post(url, "application/json", b)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	default:
		return nil, errors.New("Not a valid request type")
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(respBody) < 2 {
		return nil, errors.New("Malformed JSON response")
	}
	return respBody, nil
}

func (c *Client) Get(url string) (resp []byte, err error) {
	return c.API("get", url, nil)
}

func (c *Client) Post(url string, params interface{}) (resp []byte, err error) {
	return c.API("post", url, params)
}

func (c *Client) Profiles() (profiles *[]Profile, err error) {
	body, err := c.Get("profiles.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, profiles)
	if err != nil {
		return nil, err
	}
	return
}

func (c *Client) Update(update *NewUpdate) (resp *UpdateResponse, err error) {
	respBody, err := c.Post("updates/create.json", update)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBody, resp)
	if err != nil {
		return nil, err
	} else if !resp.Success {
		return nil, errors.New(respBody)
	}
	return resp, nil
}
