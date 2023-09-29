package indexer

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/fr"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"log"
	"os"
)

func IndexHugo(index bleve.Index, hugoPath string) {
	osFs := hugofs.Os
	config, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{
		Fs:          osFs,
		Filename:    "config.toml",
		ConfigDir:   hugoPath,
		Environment: "development",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	base := config.GetFirstLanguageConfig().BaseConfig()
	base.WorkingDir = hugoPath
	fmt.Println(base.WorkingDir)

	// Config loaded
	fs := hugofs.NewFrom(osFs, base)
	sites, err := hugolib.NewHugoSites(deps.DepsCfg{
		Fs:      fs,
		Configs: config,
	})
	err = sites.Build(hugolib.BuildCfg{SkipRender: true})
	if err != nil {
		log.Fatal("Could not render site:", err)
	}
	fmt.Println(sites.Site.Title())
	if sites != nil && (sites.Pages() != nil) {
		fmt.Printf("Index pages in site")
		fmt.Println(sites.Pages())
		for _, p := range sites.Pages() {
			if p != nil {
				if p.Kind() == "page" {
					index.Index(p.Path(), p.RawContent())
				}

			}
		}
	}

}

func IndexExist(indexBase string, hugoDir string) {
	index, err := bleve.Open(indexBase)
	if err != nil {
		panic(err)
	}
	//index.Index(message.Id, message)
	//index.Index(message2.Id, message2)
	IndexHugo(index, hugoDir)
}

func IndexNew(indexBase string, hugoDir string) {
	mapping := bleve.NewIndexMapping()
	mapping.AnalyzerNamed(fr.AnalyzerName)
	index, err := bleve.New(indexBase, mapping)
	if err != nil {
		panic(err)
	}
	//index.Index(message.Id, message)
	//index.Index(message2.Id, message2)
	IndexHugo(index, hugoDir)
}
