package database

import "time"

// TODO: this will be on a separate package

type Settings struct {
	Name         string        `json:"name" yaml:"name"`
	User         string        `json:"user" yaml:"user"`
	Password     string        `json:"password" yaml:"password"`
	Address      string        `json:"address" yaml:"address"`
	Port         int           `json:"port" yaml:"port"`
	Charset      string        `json:"charset" yaml:"charset"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"readTimeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"writeTimeout"`
}
