package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

func GetList(dir string) (map[string]string, error) {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, d := range dirs {
		if strings.HasSuffix(d.Name(), ".so") && !d.IsDir() {
			m[strings.TrimSuffix(d.Name(), ".so")] = filepath.Join(dir, d.Name())
		}
	}

	return m, nil
}

func GetFunc(name string, funcName string) (func([]byte) []byte, error) {
	p, err := plugin.Open(name)
	if err != nil {
		return nil, err
	}

	fn, err := p.Lookup(funcName)
	if err != nil {
		return nil, err
	}

	v, ok := fn.(func([]byte) []byte)
	if !ok {
		return nil, fmt.Errorf("plugin.%s is not a function", funcName)
	}

	return v, nil
}

func GetString(name string, stringName string) (string, error) {
	p, err := plugin.Open(name)
	if err != nil {
		return "", err
	}

	fn, err := p.Lookup(stringName)
	if err != nil {
		return "", err
	}

	v, ok := fn.(string)
	if !ok {
		return "", fmt.Errorf("plugin.%s is not a string", stringName)
	}

	return v, nil
}

func GetInt(name string, intName string) (int, error) {
	p, err := plugin.Open(name)
	if err != nil {
		return 0, err
	}

	fn, err := p.Lookup(intName)
	if err != nil {
		return 0, err
	}

	v, ok := fn.(int)
	if !ok {
		return 0, fmt.Errorf("plugin.%s is not an int", intName)
	}

	return v, nil
}

func GetFloat64(name string, floatName string) (float64, error) {
	p, err := plugin.Open(name)
	if err != nil {
		return 0, err
	}

	fn, err := p.Lookup(floatName)
	if err != nil {
		return 0, err
	}

	v, ok := fn.(float64)
	if !ok {
		return 0, fmt.Errorf("plugin.%s is not a float64", floatName)
	}

	return v, nil
}
