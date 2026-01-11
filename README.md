# geolite2fountries
Scripts for filtering selected countries.
## the use case
maxminddb-golang seems to be quite slow when running on a full version of GeoLite2-Country.mmdb. The use case I have here is only checking if a user is from Vietnam. The service responds much faster with only records for Vietnam.

# Install
```bash
go install -mod=readonly github.com/vanng822/geolite2fountries@latest
```

# Usage
## download
```bash
curl -L -o GeoLite2-Country.mmdb https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-Country.mmdb
```
## run filter
`countries` should be a comma-separated list

```bash
geolite2fountries --input ./GeoLite2-Country.mmdb --output ./vietnam_only.mmdb --countries VN
```