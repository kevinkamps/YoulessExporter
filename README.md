# YoulessExporter
A Youless exporter for Prometheus

## Examples

### Docker compose example
```yaml
version: '2'
services:
  youless_exporter:
    image: kevinkamps/youless_exporter:latest
    ports:
     - "80:8080"
    command: -ip=.123.123.123.123 -name="Meter1" -s0name="Meter1s0" -refreshInSeconds=1
```

### Commandline example
```bash
docker run -d \
    --name=youless_exporter \
    kevinkamps/youless_exporter:latest \
    -p 80:8080 \
      -ip=.123.123.123.123 -name="Meter1" -s0name="Meter1s0" -refreshInSeconds=1 
```


## Development notes
GO mods must be enabled for this project before you can build it.
* Terminal: `export GO111MODULE=on`
* Intellij: Settings -> Language & Frameworks -> Go -> Go modules -> Enable Go Modules