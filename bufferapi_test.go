package bufferapi

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"math/rand"
	"strconv"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type BufferSuite struct {
	AuthToken    string
	ClientId     string
	ClientSecret string
	Buffer       *Client
	UpdateIds    []string
}

var _ = Suite(&BufferSuite{})

func loadConfig(config *BufferSuite) error {
	filename := "config.json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &config)
	if err != nil {
		return err
	}
	return nil
}

func (s *BufferSuite) SetUpSuite(c *C) {
	err := loadConfig(s)
	if err != nil {
		c.Fatalf("%v", err)
	}
	config := &oauth.Config{
		ClientId:     s.ClientId,
		ClientSecret: s.ClientSecret,
		Scope:        "",
		AuthURL:      "",
		TokenURL:     "",
		TokenCache:   oauth.CacheFile(""),
	}
	transport := &oauth.Transport{Config: config}
	s.Buffer = ClientFactory(s.AuthToken, transport)
}

func (s *BufferSuite) TestGetProfiles(c *C) {
	_, err := s.Buffer.Profiles()
	c.Assert(err, IsNil)
}

func (s *BufferSuite) TestGetUpdates(c *C) {
	profiles, err := s.Buffer.Profiles()
	_, err = s.Buffer.Get("profiles/" + (*profiles)[0].Id + "/updates/pending.json")
	c.Assert(err, IsNil)
}

func (s *BufferSuite) TestNewUpdate(c *C) {
	profiles, err := s.Buffer.Profiles()
	u := NewUpdate{Text: "Test Update " + strconv.FormatInt(rand.Int63(), 10), Media: map[string]string{"link": "http://todaysvote.ca"}, ProfileIds: []string{(*profiles)[0].Id}}
	resp, err := s.Buffer.Update(&u)
	c.Assert(err, IsNil)
	updateCount := len(resp.Updates)
	id := resp.Updates[updateCount-1].Id
	s.UpdateIds = append(s.UpdateIds, id)
	c.Assert(err, IsNil)
}

func (s *BufferSuite) TearDownSuite(c *C) {
	for _, u := range s.UpdateIds {
		_, err := s.Buffer.Post("updates/"+u+"/destroy.json", nil)
		c.Assert(err, IsNil)
	}
}
