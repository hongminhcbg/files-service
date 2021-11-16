package main

import (
	"cloud.google.com/go/storage"
	"context"
	"file-service/conf"
	"file-service/service"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func main() {
	cfg := conf.Load()

	opts, err := redis.ParseURL(cfg.RedisDns)
	if err != nil {
		panic(err)
	}

	redisCli := redis.NewClient(opts)
	if err := redisCli.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	bucket, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer bucket.Close()

	bucketHandle := bucket.Bucket(cfg.GcpBucketName)
	fileService := service.NewFileService(redisCli, bucketHandle)

	http.HandleFunc("/generate_download_id", fileService.GenerateDownloadId)
	http.HandleFunc("/download", fileService.DownloadStreaming)
	http.ListenAndServe("0.0.0.0:8080", nil)
}


