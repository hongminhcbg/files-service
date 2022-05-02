package main

import (
	"context"
	"file-service/conf"
	"file-service/service"
	"file-service/store"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/go-redis/redis/v8"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db := mustConnectMysql(cfg.MySqlUrl)
	s := store.NewStore(db)
	fileService := service.NewFileService(redisCli, bucketHandle, s)

	http.HandleFunc("/generate_download_id", fileService.GenerateDownloadId)
	http.HandleFunc("/download", fileService.DownloadStreaming)
	http.HandleFunc("/views", fileService.ViewFile)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func mustConnectMysql(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
