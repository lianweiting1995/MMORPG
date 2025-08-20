package helper

import (
	"MMORPG/internal/conf"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

func Config(flagconf string) (*conf.Bootstrap, error) {
	c := config.New(config.WithSource(file.NewSource(flagconf)))
	defer c.Close()

	if err := c.Load(); err != nil {
		return nil, err
	}
	var bc conf.Bootstrap

	if err := c.Scan(&bc); err != nil {
		return nil, err
	}

	return &bc, nil
}


func rootPath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	p := filepath.Dir(filename)

	return strings.TrimRight(p, "/helper")
}

// RootPath 项目根目录
func RootPath() string {
	return rootPath()
}
