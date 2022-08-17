package main

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	Debug      bool
	EsaApiKey  string
	Team       string
	Query      string
	Template   string
	OutputType string
	Slack      *slackConfig
	Email      *emailConfig
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
	fs.String("config", "", "filepath for configuration yaml")
	fs.String("template", "", "filepath for template of patrolled result")
	pflag.CommandLine.AddGoFlagSet(fs)
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	v.AutomaticEnv()

	// FIXME: workaround in order to overwrite by env vars
	for ek, k := range map[string]string{
		"ESA_API_KEY":   "esaApiKey",
		"OUTPUT_TYPE":   "outputType",
		"SLACK_TOKEN":   "slack.token",
		"SLACK_CHANNEL": "slack.channel",
	} {
		if s := v.GetString(ek); s != "" {
			v.Set(k, s)
		}
	}

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

	var c *config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return c, nil
}
