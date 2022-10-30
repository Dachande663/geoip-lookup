# GeoIP Lookup

A single-binary GeoIP Lookup API powered by MaxMind. Can be run as a standalone binary or as a docker container.

NB: This binary does not handling fetching or updating MaxMind databases. The standard [geoipupdate](https://github.com/maxmind/geoipupdate) and cron does a far better job then we can.


## API Reference

#### Lookup IP Address

```http
  GET /lookup/{ip}
```

Return information about an IPv4 or IPv6 address, or an error if not found.

#### Status Information

```http
  GET /status
```

Return information about the current system status including database date and queries made.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file.

`HTTP_ADDR`

e.g. ":5225".

The host and port to bind to.

`HTTP_AUTH`

If set, an Authorization: Bearer ${HTTP_AUTH} header is required on all requests.

`MAXMIND_CITY_FILE`

e.g. "GeoLite2-City.mmdb"

The City mmdb database file to read from. Can be downloaded from [MaxMind.com](https://www.maxmind.com) with a free account or see snippet below.


## Fetching Databases

The following code fetches the latest GeoLite2-City database and extracts it to the current directory. This script can be added to a cron job to fetch and restart geoip-lookup.

````
curl "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=YOUR_LICENSE_KEY&suffix=tar.gz" -o GeoLite2-City.tar.gz \
  && tar -xzvf GeoLite2-City.tar.gz \
  && mv GeoLite2-Country_*/GeoLite2-City.mmdb ./GeoLite2-City.mmdb
````


## Authors

- [@Dachande663](https://www.github.com/Dachande663)


## Acknowledgements

 - [MaxMind](https://www.maxmind.com)
 
