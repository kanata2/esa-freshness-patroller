package main

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	Debug            bool
	EsaApiKey        string
	Team             string
	Query            string
	NotificationType string
	Slack            *slackConfig
	Email            *emailConfig
}

type slackConfig struct {
	Token   string
	Channel string
}

type emailConfig struct {
	From string
	To   string
}

func newConfigFrom(args []string) (*config, error) {
	v := viper.New()
	v.AutomaticEnv()

	fs := flag.NewFlagSet("esa-freshness-patroller", flag.ExitOnError)
	fs.String("query", "", "scan by query")
	fs.String("config", "", "filename for configuration yaml")
	pflag.CommandLine.AddGoFlagSet(fs)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigName("config.yaml")
	v.SetConfigName(viper.GetString("config"))
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var c *config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return c, nil
}
