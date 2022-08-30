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

### Configurations

| Name | Required | Type | Environment variable | CLI argument | key for Config file(YAML) |
| ---- | -------- | ---- | -------------------- | ------------ | ----------------- |
| esa.io's API Key | Yes | String | `ESA_API_KEY` | | esaApiKey (not recommended) |
| esa team | Yes | String | `TEAM` | `--team` | team | 
| esa's search query | Yes | String | `QUERY` | `--query` | query |
| config file | No(default: ./config.yaml) | String | `CONFIG` | `--config` | |
| debug mode | No | Bool | `DEBUG` | | debug |
| output type | No(default: 'json') | String(json or go-template) | `OUTPUT_TYPE` | `--output` | outputType |
| Go template file path | No(Yes if output type is go-template) | String | `TEMPLATE` | `--template` | template |
| destination type | No(default: 'stdout') | String(stdout or esa) | `DESTINATION` | `--destination` | destination |
| esa's post number for reporting results | No(Yes if destination type is esa) | Number | | esa.reportPostNumber |

Priority: CLI arguments > Environment variables > Config file
