## Set go compiler for the development
it makes the terminal to call the go compiler per new session.

export the below in `.zshrc` or `.bashrc` :
```shell
export PATH="$PATH:~/go/bin"
```
then:
```shell
source ~/.bashrc
```
or
```shell
source ~/.zshrc
```

Note: as the `go.mod` file is located in the `src` directory, use the `cd src` command while installing packages for the service or executing other commands which need the `main.go` file functionality.

## Init New Service

```shell
export service={service-name} && \
git clone -b main --depth 1 git@jobojet:jobojet/backend/base-code.git $service && \
cd $service && \
git commit --amend -m "[branch][main] Init" && \
make init
```

## Swagger

1. install the module in golang bin dir:
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```
2. install the swagger package in the project:
```shell
cd ./src
go get -u github.com/swaggo/swag/cmd/swag
```
3. init swagger in the project. it will create `docs` directory:
```shell
swag init -g root/main.go -o assets/docs
```
4. annotate the HTTP Handler function:

```go
// @Summary Get user by ID
// @Description Get user details by ID
// @ID get-user-by-id
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure	500 {object} ErrResponse "server status"

// @Router /user/{id} [get]
func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
// Your handler logic here
}
```
5. regenerated swagger according to the project main file:
```shell
swag init -g root/main.go -o assets/docs
```
Call the URL below to get all related documents:
```http request
{base_url}/public/swagger/index.html
```

## Logstash
Save the configuration below to a file named `logstash.conf`, and bind this file to the container as a volume.

```logstash
input {
  tcp {
    port => 5044
  }
}

filter {
  json {
    source => "message"
  }

  if [level] == "info" {
    mutate { add_field => { "custom_index" => "zap-infos" } }
  } else if [level] == "error"{
    mutate { add_field => { "custom_index" => "zap-errors" } }
  } else {
    mutate { add_field => { "custom_index" => "zap-logs" } }
  }
}

output {
  stdout { codec => rubydebug }

  elasticsearch {
    hosts => ["http://elasticsearch:9200"]
    index => "%{custom_index}"
#    index => "hello-logstash"
  }
}
```

Use the following Docker Compose configuration to set up the services.

```text
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.1
    container_name: elasticsearch
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false # Disable security for simplicity in this example
#      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - es-data:/usr/share/elasticsearch/data
    ports:
      - "9200:9200"
      - "9300:9300"
    networks:
      - elk-net

  kibana:
    image: docker.elastic.co/kibana/kibana:8.11.1
    container_name: kibana
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    networks:
      - elk-net

  logstash:
    image: docker.elastic.co/logstash/logstash:8.11.1
    container_name: logstash
    ports:
      - "5044:5044" # Beats input
      - "9600:9600"
    depends_on:
      - elasticsearch
    volumes:
      - type: bind
        source: ./logstash-conf
        target: /usr/share/logstash/pipeline
        read_only: true
    networks:
      - elk-net

volumes:
  es-data:
    driver: local

networks:
  elk-net:
    driver: bridge
```