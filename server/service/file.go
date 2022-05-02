package service

import (
	"context"
	"file-service/store"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"net/http"
)

const filename = "test.tar.gz"

type FileService struct {
	redisCli *redis.Client
	bucket   *storage.BucketHandle
	store    store.IStore
}

func NewFileService(redisCli *redis.Client, bucket *storage.BucketHandle, s store.IStore) *FileService {
	return &FileService{
		redisCli: redisCli,
		bucket:   bucket,
		store:    s,
	}
}

func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func (s *FileService) GenerateDownloadId(w http.ResponseWriter, r *http.Request) {
	cors(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	downloadId := uuid.NewString()
	err := s.redisCli.Set(context.Background(), downloadId, filename, 3*time.Second).Err()
	if err != nil {
		log.Println("internal server error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"download_id\": \"%s\"}", downloadId)))
}

func (s *FileService) DownloadStreaming(w http.ResponseWriter, r *http.Request) {
	cors(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	downloadId := r.URL.Query().Get("download_id")
	err := s.redisCli.Get(context.Background(), downloadId).Err()
	if err != nil {
		log.Println("internal server error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal server error"}`))
		return
	}

	fileReader, err := s.bucket.Object(filename).NewReader(context.Background())
	if err != nil {
		log.Println("internal server error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal server error"}`))
		return
	}

	defer fileReader.Close()

	//attachment file, browser will handler that
	w.Header().Set("Content-Disposition", "attachment; filename=test.tar.gz")
	io.Copy(w, fileReader)
}

func (s *FileService) ViewFile(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	idStr2 := queries.Get("id")
	if len(idStr2) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"not ok"}`))
		return
	}

	idInt, _ := strconv.Atoi(idStr2)
	template, err := s.store.GetTemplateById(context.Background(), int64(idInt))
	if err != nil || template == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"not ok"}`))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(template.Content)
	return
}
