package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"smesearch/indexer"
	"smesearch/searcher"
	"smesearch/service"
	"strings"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
	hugoDir     string
	indexDir    string
	indexName   string

	rootCmd = &cobra.Command{
		Use:   "smesearch",
		Short: "A search library for kb",
		Long:  `A search library for kb`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&indexDir, "indexdir", "~/.config", "where indexer will be used or created")
	rootCmd.PersistentFlags().StringVar(&hugoDir, "hugodir", ".", "Hugo path to index")
	rootCmd.PersistentFlags().StringVar(&indexName, "indexname", "smesearch", "indexdirectory name, a diretory with this name will be created at this location ")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(indexCmd)
	rootCmd.AddCommand(serveCmd)
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "Samuel MEYNARD samuel@meyn.fr")
	viper.SetDefault("license", "apache")
	if hugoDir == "." {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("Error :", err)
			os.Exit(1)
		}
		hugoDir = filepath.Dir(ex)
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve Daemon for Hugo",
	Run: func(cmd *cobra.Command, args []string) {
		service.Serve(indexDir, indexName)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search something in Corpus",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lindex := indexer.GetIndex(indexDir, indexName)
		search := searcher.Search(lindex, args)
		for _, hit := range search.Hits {
			fmt.Println(hit.Score, "https://kb.local.meyn.fr/"+strings.TrimSuffix(hit.ID, ".md")+"/")
		}
	},
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index Corpus",
	Run: func(cmd *cobra.Command, args []string) {
		lindex := indexer.GetIndex(indexDir, indexName)
		// Check in indexBase exist
		info, err := os.Stat(lindex)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Indexing ...: " + lindex)
			indexer.IndexNew(lindex, hugoDir)
			os.Exit(0)
		}
		if !info.IsDir() {
			fmt.Println(lindex + ", is not a directory\n")
			os.Exit(1)
		}
		fmt.Println("Update index ...: " + lindex)
		indexer.IndexExist(lindex, hugoDir)
	},
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
