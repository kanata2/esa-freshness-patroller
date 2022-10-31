package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	Debug             bool
	EsaApiKey         string
	Team              string
	Query             string
	Template          string
	Output            string
	Destination       string
	OutdatedThreshold int
	Slack             *slackConfig
	Email             *emailConfig
	Esa               *esaConfig
}

type slackConfig struct {
	Token   string
	Channel string
}

type emailConfig struct {
	From string
	To   string
}

type esaConfig struct {
	ReportPostNumber int
}

func New(args []string) (*config, error) {
	v := viper.New()

	fs := flag.NewFlagSet("esa-freshness-patroller", flag.ExitOnError)
	fs.String("team", "", "esa.io's team")
	fs.String("query", "", "esa.io's search query for scanning. more details: https://docs.esa.io/posts/104")
	fs.String("output", "", "output type(value: json or go-template)")
	fs.String("destination", "", "destination of output(value: stdout, esa)")
	fs.String("config", "", "filepath for configuration yaml")
	fs.String("template", "", "filepath for template of patrolled result")
	fs.Int("outdated-threshold", 180, "filepath for template of patrolled result")
	pflag.CommandLine.AddGoFlagSet(fs)
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)
	// Alias
	v.BindPFlag("outdatedthreshold", pflag.CommandLine.Lookup("outdated-threshold"))

	v.AutomaticEnv()

	// FIXME: workaround in order to overwrite by env vars
	for ek, k := range map[string]string{
		"ESA_API_KEY":        "esaApiKey",
		"OUTPUT":             "output",
		"OUTDATED_THRESHOLD": "outdatedThreshold",
		"DESTINATION":        "destination",
		"SLACK_TOKEN":        "slack.token",
		"SLACK_CHANNEL":      "slack.channel",
		"ESA_REPORT_POST":    "esa.reportPostNumber",
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
