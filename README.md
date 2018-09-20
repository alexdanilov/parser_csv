# CSV uploader

CSV uploader contains from two parts:
* service – a server accepts connections and data in JSON format. Stores data in memory as a map.
* parser – a CLI tool to parse CSV files, put data to a queue and send it to a service when possible.

## How to run
`docker-compose up`

### How to test
`cd item; go test`
