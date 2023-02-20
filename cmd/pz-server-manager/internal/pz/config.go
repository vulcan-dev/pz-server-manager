package pz

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Option struct {
	Name         string
	Desc         string
	Value        interface{}
	Type         string
	Min          int  // For int & float
	Max          int  // For int & float
	MinSpecified bool // For int & float
	MaxSpecified bool // For int & float
}

type Options map[string]Option

func Parse(name string, path string) (Options, error) {
	path = filepath.Join(path, "Server/"+name+".ini")
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	options := make(map[string]Option)
	for _, section := range cfg.Sections() {
		for _, key := range section.Keys() {
			// Clean up the comment
			key.Comment = strings.Replace(key.Comment, "#", "", -1)
			if len(key.Comment) > 0 && key.Comment[0] == ' ' {
				key.Comment = key.Comment[1:]
			}

			option := Option{
				Name:  key.Name(),
				Desc:  key.Comment,
				Value: key.Value(),
			}

			// Find min/max values
			if strings.Contains(key.Comment, "Minimum=") {
				min := strings.Split(strings.Split(key.Comment, "Minimum=")[1], " ")[0]
				min = strings.Split(min, ",")[0]
				min = strings.Split(min, ".")[0]
				option.Min, _ = strconv.Atoi(min)
				option.MinSpecified = true
			}

			if strings.Contains(key.Comment, "Maximum=") {
				max := strings.Split(strings.Split(key.Comment, "Maximum=")[1], " ")[0]
				max = strings.Split(max, ",")[0]
				max = strings.Split(max, ".")[0]
				option.Max, _ = strconv.Atoi(max)
				option.MaxSpecified = true
			}

			// Find the type depending on the value
			value := key.Value()
			if value == "true" || value == "false" {
				option.Type = "bool"
				option.Value, _ = strconv.ParseBool(value)
			} else if _, err := strconv.Atoi(value); err == nil {
				option.Type = "int"
				option.Value, _ = strconv.Atoi(value)
			} else if _, err := strconv.ParseFloat(value, 64); err == nil {
				option.Type = "float"
				option.Value, _ = strconv.ParseFloat(value, 64)
			} else {
				option.Type = "string"
			}

			options[key.Name()] = option
		}
	}

	return options, nil
}

func (o *Options) Get(name string) (Option, bool) {
	option, ok := (*o)[name]
	return option, ok
}

func (o *Options) Set(name string, value interface{}) error {
	option, ok := (*o)[name]
	if !ok {
		return errors.New("option not found")
	}

	// Check if they are the same type
	if reflect.TypeOf(value) != reflect.TypeOf(option.Value) {
		return errors.New("value type does not match option type")
	}

	switch option.Type {
	case "int":
		if option.Min > 0 && value.(int) < option.Min {
			return errors.New("value is below minimum")
		}
		if option.Max > 0 && value.(int) > option.Max {
			return errors.New("value is above maximum")
		}
	case "float":
		if option.Min > 0 && value.(float64) < float64(option.Min) {
			return errors.New("value is below minimum")
		}
		if option.Max > 0 && value.(float64) > float64(option.Max) {
			return errors.New("value is above maximum")
		}
	}

	option.Value = value
	return nil
}

func (o *Options) Save(name string, path string) error {
	path = filepath.Join(path, "Server/"+name+".ini")
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	for _, section := range cfg.Sections() {
		for _, key := range section.Keys() {
			option, ok := (*o)[key.Name()]
			if !ok {
				continue
			}

			key.SetValue(fmt.Sprintf("%v", option.Value))
		}
	}

	return cfg.SaveTo(path)
}

func (o *Options) JSON() (map[string]interface{}, error) {
	jsonObj := make(map[string]interface{})
	for _, option := range *o {
		jsonObj[option.Name] = map[string]interface{}{
			"key":   option.Name,
			"desc":  option.Desc,
			"value": option.Value,
			"type":  option.Type,
		}

		if option.MinSpecified {
			jsonObj[option.Name].(map[string]interface{})["min"] = option.Min
		}
		if option.MaxSpecified {
			jsonObj[option.Name].(map[string]interface{})["max"] = option.Max
		}
	}

	return jsonObj, nil
}
