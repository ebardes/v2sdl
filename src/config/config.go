package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

type Service interface {
	Start(*Config) error
	Stop()
	Refresh(*Config) error
	Name() string
}

type Config struct {
	DebugLevel int    `json:"debug"`
	Universe   int    `json:"universe"`
	Address    int    `json:"address"`
	Interface  string `json:"interface"`
	Protocol   string `json:"protocol"`
	Storage    string `json:"storage"`

	location string
	services map[string]Service
}

var Media Content

type Content struct {
	Groups map[int]map[int]Item
	home   string
}

type Item struct {
	File string
	Name string
	Type string
	home string
}

func Load(fn string) (cfg Config, err error) {
	if fn == "" {
		fn = "config.json"
	}
	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()

	j := json.NewDecoder(f)
	err = j.Decode(&cfg)
	if err != nil {
		return
	}

	cfg.location = fn

	Media.home = cfg.Storage
	if Media.home == "" {
		Media.home = os.ExpandEnv("${HOME}/media")
	}
	os.MkdirAll(Media.home, 0777)

	fn = path.Join(Media.home, "media.json")
	m, err := os.Open(fn)
	if err != nil {
		err = nil
		Media.Groups = make(map[int]map[int]Item)

		Media.Save()
		return
	}
	defer m.Close()

	j = json.NewDecoder(m)
	err = j.Decode(&Media.Groups)

	data, err := json.Marshal(Media.Groups)
	log.Debug().Err(err).Msg(fmt.Sprint(string(data)))
	return
}

func (c *Content) Get(group, item int) *Item {
	if g, ok := c.Groups[group]; ok {
		if item, ok := g[item]; ok {
			groupdir := fmt.Sprintf("group_%03d", group)
			item.home = path.Join(c.home, groupdir, item.File)
			return &item
		}
	}
	return nil
}

func (c *Content) Save() {
	fn := path.Join(Media.home, "media.json")
	f, err := os.Create(fn)
	if err != nil {
		return
	}
	defer f.Close()

	j := json.NewEncoder(f)
	j.Encode(c.Groups)
}

func (i *Item) Path() string { return i.home }

func (c *Config) AddAndStartService(s Service) {
	if c.services == nil {
		c.services = make(map[string]Service)
	}

	name := s.Name()
	srv, ok := c.services[name]
	if ok {
		srv.Stop()
	}

	c.services[name] = s
	s.Start(c)
}

func (c *Config) StopAll() {
	for _, srv := range c.services {
		srv.Stop()
	}
	c.services = nil
}
