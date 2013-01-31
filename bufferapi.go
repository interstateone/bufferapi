package bufferapi

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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

type Media map[string]string

type Update struct {
	Id             string `json:"id"`
	CreatedAt      int    `json:"created_at"`
	Day            string `json:"day"`
	DueAt          int    `json:"due_at"`
	DueTime        string `json:"due_time"`
	Media          media  `json:"media"`
	ProfileId      string `json:"profile_id"`
	ProfileService string `json:"profile_service"`
	Status         string `json:"status"`
	Text           string `json:"text"`
	TextFormatted  string `json:"text_formatted"`
	UserId         string `json:"user_id"`
	Via            string `json:"via"`
}

type UpdateResponse struct {
	Success          bool     `json:"success"`
	BufferCount      int      `json:"buffer_count"`
	BufferPercentage int      `json:"buffer_percentage"`
	Updates          []Update `json:"updates"`
}

type Valuer interface {
	UrlValues() url.Values
}

func ClientFactory(token string, transport *oauth.Transport) *Client {
	t := &oauth.Token{AccessToken: token}
	transport.Token = t
	c := Client{AccessToken: token, transport: transport}
	return &c
}

func (c *Client) API(method, uri string, data Valuer) (respBody []byte, err error) {
	jsonPattern, _ := regexp.Compile(`\.json$`)
	if !jsonPattern.Match([]byte(uri)) {
		uri += ".json"
	}
	uri = "https://api.bufferapp.com/1/" + uri

	var resp *http.Response
	switch method {
	case "get":
		resp, err = c.transport.Client().Get(uri)
		if err != nil {
			return nil, err
		}

	case "post":
		var values url.Values
		if data == nil {
			values = make(url.Values)
		} else {
			values = data.UrlValues()
		}
		resp, err = c.transport.Client().PostForm(uri, values)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("Not a valid request type")
	}
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll((*resp).Body)

	if resp.StatusCode >= 400 {
		return nil, errors.New((*resp).Status)
	}

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

func (c *Client) Post(url string, params Valuer) (resp []byte, err error) {
	return c.API("post", url, params)
}

func (c *Client) Profiles() (profiles *[]Profile, err error) {
	body, err := c.Get("profiles.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &profiles)
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
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	} else if !resp.Success {
		return nil, errors.New(string(respBody))
	}
	return resp, nil
}

func (u *NewUpdate) UrlValues() (values url.Values) {
	values = make(url.Values)
	values.Set("text", u.Text)
	for key, value := range u.Media {
		values.Set("media["+key+"]", value)
	}
	for _, profile := range u.ProfileIds {
		values.Set("profile_ids[]", profile)
	}
	values.Set("shorten", strconv.FormatBool(u.Shorten))
	values.Set("now", strconv.FormatBool(u.Now))
	return
}
