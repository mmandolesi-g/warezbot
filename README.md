# Warezbot

Slackbot to interact with Emby, Radarr and Sonarr.

## Setup

Point the `configFile` flag to a file containing your config information. For example:
```
{
  "loglevel": "debug",
  "tlsconfig": {
    "tlsca": "-----BEGIN CERTIFICATE-----XXX-----END CERTIFICATE-----\n",
    "tlscert": "-----BEGIN CERTIFICATE-----XXX-----END CERTIFICATE-----\n",
    "tlskey": "-----BEGIN RSA PRIVATE KEY-----XXX-----END RSA PRIVATE KEY-----\n"
  },
  "slack": {
    "bottoken": "xoxb-xxx",
    "botid": "xxx",
    "channelid": "xxx"
  },
  "emby": {
    "adminid": "xxx",
    "path": "https://emby.example.com",
    "token": "xxx"
  },
  "radarr": {
    "path": "https://radarr.example.com",
    "apikey": "xxx"
  }
}
```

## Authors

* **Marcelo Mandolesi**

## License

This project is licensed under the MIT License.

