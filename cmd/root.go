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
	// Version - Current version of gh-sso.
	VERSION = "0.0.1"
)

// Configuration object.
type Config struct {
	// GitHub API token.
	Token         string   `mapstructure:"token"`
	Organizations []string `mapstructure:"organizations,flow"`
}

// lstCmd flags.
type ListFlags struct {
	Enterprise   string
	Organization string
}

// User - User details.
type User struct {
	Name       string
	SamlNameId string
}

// Repository - Individual repository in an organization.
type Repository struct {
	Name string
}

// Organization - Top-Level organization; also used for Enterprise member orgs.
type Organization struct {
	ID           string
	Login        string
	Name         string
	Repositories map[string]Repository
	Users        map[string]User
}

// Enterprise - Top-level enterprise.
type Enterprise struct {
	ID            string
	Login         string
	Name          string
	Organizations map[string]Organization
	Users         map[string]User
}

var (
	// Client for use with githubv4.
	Client *githubv4.Client

	// Context.
	Ctx context.Context

	// Configuration.
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
