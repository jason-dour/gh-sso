// cmd/root
//
// Root command for gh-sso.

package cmd

import (
	"context"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	// Version is the current version of gh-sso.
	VERSION = "0.0.1"
)

// Config is the configuration struct.
type Config struct {
	// GitHub API token.
	Token         string   `mapstructure:"token"`
	Organizations []string `mapstructure:"organizations,flow"`
}

type ListFlags struct {
	Enterprise   string
	Organization string
}

type User struct {
	Login      string
	SamlNameId string
}

type Repository struct {
	Name string
}

type Organization struct {
	Name         string
	Repositories map[string]Repository
	Users        map[string]User
}

type Enterprise struct {
	Name          string
	Organizations map[string]Organization
	Users         map[string]User
}

var (
	// Client is the githubv4 client.
	Client *githubv4.Client

	// Ctx is the context.
	Ctx context.Context

	// Struct to hold the configuration.
	Cfg Config

	// Config file given on command line.
	cfgFile string

	// Root command.
	rootCmd = &cobra.Command{
		Version: VERSION,
		Use:     "gh-sso",
		Short:   "Organization & Enterprise SSO User Management for Github.com",
		Long:    `Organization & Enterprise SSO User Management for Github.com`,
	}
)

// Execute is the entry point for the root command.
func Execute() {
	panicOnError(rootCmd.Execute())
}

// panicOnError panics if err is not nil.
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Load config file.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find the home directory.
		home, err := homedir.Dir()
		panicOnError(err)

		// Look in the home directory and the current working directory for the config file.
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

		// Define the config file name.
		viper.SetConfigName(".gh-sso")
	}

	// Read in environment variables that match.
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; ignore error
	} else {
		// Config file was found but another error was produced
		panicOnError(err)
	}

	// Unmarshal the config file into the config struct.
	err = viper.Unmarshal(&Cfg)
	panicOnError(err)

	// Initialize the http client.
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Cfg.Token},
	)
	httpClient := oauth2.NewClient(Ctx, src)

	// Initialize the githubv4 client.
	Client = githubv4.NewClient(httpClient)
}

// Initialization.
func init() {
	// Initialize config.
	cobra.OnInitialize(initConfig)

	// Define the command line flags.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-sso.[yaml|toml|json|hcl])")

	// Define the context.
	Ctx = context.Background()
}
