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

	fs := flag.NewFlagSet("esa-freshness-patroller", flag.ExitOnError)
	fs.String("query", "", "scan by query")
	fs.String("config", "", "filename for configuration yaml")
	pflag.CommandLine.AddGoFlagSet(fs)
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigName("config")
	if cfgPath := v.GetString("config"); cfgPath != "" {
		v.SetConfigFile(cfgPath)
	}
	if err := v.ReadInConfig(); err != nil {
		if e, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, e
		}
	}
	v.AutomaticEnv()

	// FIXME: workaround in order to overwrite by env vars
	for ek, k := range map[string]string{
		"ESA_API_KEY":       "esaApiKey",
		"NOTIFICATION_TYPE": "notificationType",
		"SLACK_TOKEN":       "slack.token",
		"SLACK_CHANNEL":     "slack.channel",
	} {
		if s := v.GetString(ek); s != "" {
			v.Set(k, s)
		}
	}

	var c *config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return c, nil
}
