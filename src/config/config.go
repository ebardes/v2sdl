package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

type Config struct {
	DebugLevel int    `json:"debug"`
	Universe   int    `json:"universe"`
	Address    int    `json:"address"`
	Interface  string `json:"interface"`
	Protocol   string `json:"protocol"`
	Storage    string `json:"storage"`
}

var Media Content

type Content struct {
	Groups map[int]Group
	home   string
}

type Group struct {
	Items map[int]Item
}

type Item struct {
	File string
	Name string
	Type string
	home string
}

func Load() (cfg Config, err error) {
	f, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer f.Close()

	j := json.NewDecoder(f)
	err = j.Decode(&cfg)
	if err != nil {
		return
	}

	Media.home = cfg.Storage
	if Media.home == "" {
		Media.home = os.ExpandEnv("${HOME}/media")
	}
	os.MkdirAll(Media.home, 0777)

	fn := path.Join(Media.home, "media.json")
	m, err := os.Open(fn)
	if err != nil {
		err = nil
		Media.Groups = make(map[int]Group)

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
		if item, ok := g.Items[item]; ok {
			groupdir := fmt.Sprintf("group_%03d", group)
			item.home = path.Join(c.home, groupdir, item.File)
			return &item
		}
	}
	return nil
}

func (c *Content) Put(group, item int, x Item) {
	g, ok := c.Groups[group]
	if !ok {
		g = Group{}
		g.Items = make(map[int]Item)
		c.Groups[group] = g
	}
	g.Items[item] = x

	go c.Save()
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
