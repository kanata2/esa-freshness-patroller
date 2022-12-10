esa-freshness-patroller :cop:
===

This is a PoC to maintain the documentation's freshness in esa.io.  
This is inspired by Google's efforts as described in "Software Engineering at Google".

## How to use

1. Adding metadata in a fenced code which has `esa-freshness-patroller` info string to each post
2. Run `esa-freshness-patroller`

### Metadata for each post
Metadata can be used for patrolling status and indivisual configurations per post.
Metadata must be written in YAML according to the following schema.

| Key | Required | Type | Description |
| --- | -------- | ---- | ----------- |
| `owners` | Yes | Array of string | documentation owner/reviewer |
| `last_checked_at` | Yes | String(format: `YYYY/MM/DD`) | last reviewed date by owner/reviewers |
| `interval` | No | Numbner | day for patrolling interval. this takes precedence over oeverall configuration |
| `skip` | No | Bool | excluding from patrolling targets |
| `custom` | No | Mapping(key/value is String) | custom metadata that can be added freely |


Example:

````markdown
```esa-freshness-patroller
owners:
  - @kanata2
  - @kanata1
last_checked_at: 2022/12/11
interval: 30
skip: false
custom:
  category: daily
  team: SRE
```
````

### General configurations
You can set general configurations through CLI arguments, environment variables or config file.

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
