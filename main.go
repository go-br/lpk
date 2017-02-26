package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/crgimenes/goConfig"
)

type config struct {
	PackageName string `cfg:"name"`
	GoPath      string
	List        string `cfgDefault:"skipvendor"`
	ListAll     bool   `cfg:"-"`
	SkipVendor  bool   `cfg:"-"`
}

var cfg config

func parseListPar() (err error) {
	v := strings.Split(cfg.List, ",")

	for _, p := range v {
		switch p {
		case "skipvendor":
			cfg.SkipVendor = true
		case "all":
			cfg.ListAll = true
		default:
			return fmt.Errorf("Unknow -list parameter %v", p)
		}
	}
	return
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		return nil
	}

	if cfg.SkipVendor && f.Name() == "vendor" {
		return filepath.SkipDir
	}

	pkgName := strings.ToLower(cfg.PackageName)
	dirName := strings.ToLower(f.Name())

	if pkgName == dirName {
		fmt.Println(path)
		if !cfg.ListAll {
			return io.EOF
		}
	}

	return nil
}

func main() {

	cfg = config{}

	err := goConfig.Parse(&cfg)
	if err != nil {
		println(err.Error())
		return
	}

	if cfg.PackageName == "" {
		lastPar := flag.NArg() - 1
		cfg.PackageName = flag.Arg(lastPar)
		if cfg.PackageName == "" {
			println("Package name not defined.")
			goConfig.Usage()
			return
		}
	}

	if cfg.GoPath == "" {
		cfg.GoPath = build.Default.GOPATH
	}

	root := cfg.GoPath + "/src"

	err = parseListPar()
	if err != nil {
		println(err.Error())
		return
	}

	_, err = os.Stat(root)
	if err != nil {
		println(err.Error())
		return
	}

	err = filepath.Walk(root, visit)
	if err != nil && err != io.EOF {
		println(err.Error())
		return
	}
}
