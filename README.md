esa-freshness-patroller :cop:
===

This is a PoC to maintain the documentation's freshness in esa.io.  
This is inspired by Google's efforts as described in "Software Engineering at Google".

## How to use

1. Adding `Last checked at YYYY/MM/DD by @username` in your documents
2. Run `esa-freshness-patroller` as below

``` sh
$ cat <<EOF > config.yaml
team: kanata2-sandbox
query: 'in:"Users/kanata2" Spec'
EOF

$ ESA_API_KEY=xxx esa-freshness-patroller
```

If you want to notify Slack, then run as below.

``` sh
$ cat <<EOF > config.yaml
team: kanata2-sandbox
query: 'in:"Users/kanata2" Spec'
notificationType: slack
slack:
  channel: CXXXXXXX
EOF

$ ESA_API_KEY=xxx SLACK_TOKEN=yyy esa-freshness-patroller
```

### Configurations

| Name | Required | Type | Environment variable | CLI argument | key for Config file(YAML) |
| ---- | -------- | ---- | -------------------- | ------------ | ----------------- |
| esa.io's API Key | Yes | String | `ESA_API_KEY` | | esaApiKey (not recommended) |
| esa team | Yes | String | `TEAM` | `--team` | team | 
| esa's search query | Yes | String | `QUERY` | `--query` | query |
| config file | No(default: ./config.yaml) | String | `CONFIG` | `--config` | |
| template file | No | String | `TEMPLATE` | `--template` | template |
| debug mode | No | Bool | `DEBUG` | | debug |
| output type | No(default: '') | String | `OUTPUT_TYPE` | | outputType |
| slack's Token | No | String | `SLACK_TOKEN` | | slack.token (not recommended) |
| slack's notification channel | No | String | | | slack.channel |

Priority: CLI arguments > Environment variables > Config file
