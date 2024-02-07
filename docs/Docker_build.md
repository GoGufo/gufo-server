For build Gufo in Docker

```
docker build --no-cache -t gufo:latest -f Dockerfile .
```

For Run Gufo in Docker

```
docker run --name gufoserver \
--restart=always \
-p 8090:8090 \
-d gufo:latest

```

Config and files directories on Docker image:

Config file (settings.toml): /var/gufo/config
Translations: /var/gufo/lang
Email templates: /var/gufo/templates
Logs: /var/gufo/log
Files: /var/gufo/files
