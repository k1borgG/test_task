package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	ig "github.com/k1borgG/test_task/internal/grpc"
	"github.com/k1borgG/test_task/internal/repository"
	"github.com/k1borgG/test_task/internal/service"
	"github.com/k1borgG/test_task_grpc"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	ElasticsearchURLs     []string `envconfig:"ELASTIC_URLS" default:"http://elasticsearch:9200"`
	ElasticsearchUsername string   `envconfig:"ELASTIC_USERNAME"`
	ElasticsearchPassword string   `envconfig:"ELASTIC_PASSWORD"`
}

func parsCFG(cfg Config) (EsCfg elasticsearch.Config) {
	EsCfg = elasticsearch.Config{
		Addresses: cfg.ElasticsearchURLs,
		Username:  cfg.ElasticsearchUsername,
		Password:  cfg.ElasticsearchPassword,
	}
	return EsCfg
}

func ConnectToElastic(EsCfg elasticsearch.Config) (client *elasticsearch.Client, err error) {
	client, err = elasticsearch.NewClient(EsCfg)
	if err != nil {
		log.Printf("New client error in Elasticsearch: %v", err)
		return nil, err
	}

	return client, err
}

func main() {
	fmt.Println("Иммитация работы DevOps'a... Pls Wait 1 minute.")
	time.Sleep(1 * time.Minute)

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Error processing environment variables: %s", err)
	}
	EsCfg := parsCFG(cfg)

	es, err := ConnectToElastic(EsCfg)
	if err != nil {
		log.Printf("Failed to connect ES")
	}

	err = RunInitEs(context.Background(), es)
	if err != nil {
		log.Printf("Failed to Run ES , %s", err)
	}

	esRepo := repository.NewElasticsearchRepository(es, "my_custom_index5")

	productService := service.NewProductService(esRepo)

	go func() {
		s := grpc.NewServer()
		srv := ig.NewGRPCServer(productService)
		test_task_grpc.RegisterProductServiceServer(s, srv)

		reflection.Register(s)

		l, err := net.Listen("tcp", "0.0.0.0:8080")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		if err = s.Serve(l); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Get finish signbal: %s. Stop...", sig)
}

func RunInitEs(ctx context.Context, es *elasticsearch.Client) error {
	indexName := "my_custom_index5"

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body: strings.NewReader(`{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1
  },
  "mappings": {
    "properties": {
      "ID": {
        "type": "keyword"
      },
      "Name": {
        "type": "text"
      },
      "Description": {
        "type": "text"
      },
      "Brand": {
        "type": "keyword"
      },
      "Model": {
        "type": "keyword"
      },
      "coordinates": {
        "type": "geo_point"
      },
      "Date": {
        "type": "date",
        "format": "yyyy-MM-dd"
      },
      "Price": {
        "type": "integer"
      }
    }
  }
}`),
	}
	res, err := req.Do(ctx, es)
	if err != nil {
		log.Printf("Request error: %s", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Error creating index: %s", err)
		}
	}(res.Body)

	if res.IsError() {
		var buf bytes.Buffer
		if _, err = buf.ReadFrom(res.Body); err != nil {
			log.Printf("error parsing response body: %s", err)
			return err
		}
		log.Printf("error response from Elasticsearch: %s", buf.String())
		return err
	}

	fmt.Printf("Index %s created successfully\n", indexName)
	return err
}
