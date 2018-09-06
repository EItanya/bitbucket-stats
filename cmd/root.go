package cmd

import (
	"bitbucket-stats/api"
	"bitbucket-stats/cache"
	"bitbucket-stats/logger"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile, username, password, cacheType, url string
	client                                      *api.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "bitbucket-stats",
	Short:             "A stats aggregator for bitbucket",
	PersistentPreRunE: cobraSetup,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".cobra.yaml", "config file (default is .cobra.yaml)")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "username for bitbucket auth")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "password for bitbucket auth")
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "url of bitbucket instance")
	rootCmd.PersistentFlags().StringVarP(&cacheType, "cache", "c", "redis", "type of cache used for local storage")
	rootCmd.PersistentFlags().BoolP("force", "f", false, "Used to force a given action")

	viper.SetDefault("cache", "redis")
	viper.SetDefault("force", false)
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("cache", rootCmd.PersistentFlags().Lookup("cache"))
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("force", rootCmd.PersistentFlags().Lookup("force"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".cobra-demo" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Log.Info("Using config file:", viper.ConfigFileUsed())
	}
}

func cobraSetup(cmd *cobra.Command, args []string) error {
	userInfo, err := setupUser(cmd)
	if err != nil {
		logger.Log.Warn("Failed to correctly parse user info")
		return err
	}
	clientCache, err := setupCache(cmd)
	if err != nil {
		logger.Log.Warn("Failed to create local cache")
		return err
	}
	client, err = setupClient(cmd, userInfo, clientCache)
	if err != nil {
		logger.Log.Warn("Failed to setup API client")
		return err
	}
	logger.Log.Info("Successfully finished setup, moving on")
	return nil
}

func setupUser(cmd *cobra.Command) (*api.UserInfo, error) {
	userInfo := &api.UserInfo{}
	if username := viper.GetString("username"); username != "" {
		userInfo.Username = username
	} else if username := viper.GetString("UserInfo.username"); username != "" {
		userInfo.Username = username
	} else {
		return nil, errors.New("No username found to authenticate with")
	}

	if password := viper.GetString("password"); password != "" {
		userInfo.Password = password
	} else if password := viper.GetString("UserInfo.password"); password != "" {
		userInfo.Password = password
	} else {
		return nil, errors.New("No password found to authenticate with")
	}
	logger.Log.Info("User info successfully parsed")
	return userInfo, nil
}

func setupCache(cmd *cobra.Command) (cache.Cache, error) {
	cacheType := viper.GetString("cache")
	switch cacheType {
	case "file":
		fileCache, err := cache.NewFileCache(nil)
		if err != nil {
			return nil, err
		}
		logger.Log.Info("Local File cache successfully setup")
		return fileCache, nil
	default:
		redisCache, err := cache.NewRedisCache(nil)
		if err != nil {
			return nil, err
		}
		logger.Log.Info("Local Redis cache successfully setup")
		return redisCache, nil
	}
}

func setupClient(cmd *cobra.Command, user *api.UserInfo, clientCache cache.Cache) (*api.Client, error) {
	if url := viper.GetString("url"); url != "" {
		forceReset := false
		if force := viper.GetBool("force"); force {
			forceReset = forceReset || force
			logger.Log.Info("Force set to true for client, updating cache")
		}
		return api.Initialize(user, clientCache, url, forceReset)
	}
	return nil, errors.New("No url found to connect to bitbucket with")
}
