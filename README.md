# geolite2fountries
Scripts for filtering selected countries

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
```bash
geolite2fountries --input ./GeoLite2-Country.mmdb --output ./vietnam_only.mmdb --countries VN
```