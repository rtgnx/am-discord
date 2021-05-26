# Alert Manager Discord Plugin

[![Build Status](https://ci.revlabs.xyz/api/badges/rtgnx/am-discord/status.svg)](https://ci.revlabs.xyz/rtgnx/am-discord)

```YML
version: '3.8'
services:
  image: 'rtgnx9/am-discord'
  environment:
    DISCORD_WEBHOOK: ${DISCORD_WEBHOOK}
```


### Alert Manager Config

```YML
  global:
    # The smarthost and SMTP sender used for mail notifications.
    smtp_smarthost: 'localhost:25'
    smtp_from: ''
    smtp_auth_username: 'alertmanager'
    smtp_auth_password: 'password'

  # The directory from which notification templates are read.
  templates: 
  - '/etc/alertmanager/template/*.tmpl'

  # The root route on which each incoming alert enters.
  route:
    group_by: ['alertname']
    group_wait: 20s
    group_interval: 5m
    repeat_interval: 3h 
    receiver: discord_webhook

  receivers:
  - name: 'discord_webhook'
    webhook_configs:
    - url: 'http://am-discord:9094'
```