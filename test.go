package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

type Server struct {
	s3Client    *s3.Client
	vaultClient *api.Client
}

type Secrets struct {
	AWS_ACCESS_KEY_ID     string
	AWS_BUCKET_NAME       string
	AWS_BUCKET_URL        string
	AWS_REGION            string
	AWS_SECRET_ACCESS_KEY string
}

var secrets Secrets

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env, reading from OS envs...")
	}

	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = os.Getenv("VAULT_ADDR")
	vaultClientInit, err := api.NewClient(vaultConfig)
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}
	vaultClientInit.SetToken(os.Getenv("VAULT_TOKEN"))

	secrets, err = loadSecretsFromVault(vaultClientInit)
	if err != nil {
		log.Fatalf("Failed to load secrets from Vault: %v", err)
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(secrets.AWS_REGION),
		config.WithBaseEndpoint(secrets.AWS_BUCKET_URL),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				secrets.AWS_ACCESS_KEY_ID,
				secrets.AWS_SECRET_ACCESS_KEY,
				"",
			),
		),
	)

	server := &Server{
		s3Client: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.DisableLogOutputChecksumValidationSkipped = true
		}),
		vaultClient: vaultClientInit,
	}

	r := chi.NewRouter()
	r.Get("/api/v1/cbesuperapp/files/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "200", "message": "File service is active 🚀"})
	})
	r.Get("/api/v1/cbesuperapp/files/*", server.streamHandler)
	// Add POST upload endpoint (supports multipart/form-data with "file" field or raw body upload)
	r.Post("/api/v1/cbesuperapp/files/*", server.uploadHandler)

	srv := &http.Server{
		Addr:         ":" + "8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// Graceful shutdown setup
	idleConnsClosed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("🛑 Shutdown signal received, shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("⚠️ HTTP server Shutdown Error: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("🚀 Server listening on port %s...", "8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("⚠️ HTTP server ListenAndServe Error: %v", err)
	}
	<-idleConnsClosed
	log.Println("✅ Server stopped")
}

func loadSecretsFromVault(vaultClient *api.Client) (Secrets, error) {
	secret, err := vaultClient.Logical().Read(os.Getenv("VAULT_PATH"))
	if err != nil {
		return Secrets{}, err
	}
	if secret == nil || secret.Data == nil {
		return Secrets{}, nil
	}

	data, ok := secret.Data["data"].(map[string]any)
	if !ok {
		return Secrets{}, nil
	}

	return Secrets{
		AWS_ACCESS_KEY_ID:     data["AWS_ACCESS_KEY_ID"].(string),
		AWS_BUCKET_NAME:       data["AWS_BUCKET_NAME"].(string),
		AWS_BUCKET_URL:        data["AWS_BUCKET_URL"].(string),
		AWS_REGION:            data["AWS_REGION"].(string),
		AWS_SECRET_ACCESS_KEY: data["AWS_SECRET_ACCESS_KEY"].(string),
	}, nil
}

func (serv *Server) streamHandler(w http.ResponseWriter, r *http.Request) {
	prefix := "/api/v1/cbesuperapp/files/"
	key := r.URL.Path
	if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
		key = key[len(prefix):]
	}

	// If path is /files/<bucket_name>/<key...>, extract bucket from the path.
	bucketNamePrefix := ""
	if key != "" {
		parts := strings.SplitN(strings.Trim(key, "/"), "/", 2)
		if len(parts) >= 2 {
			bucketNamePrefix = strings.Trim(parts[0], "/")
			key = parts[1]
		}
	}

	// If no bucket in path, fall back to query param or header.
	if bucketNamePrefix == "" {
		bucketNamePrefix = strings.Trim(r.URL.Query().Get("bucket_name"), "/")
		if bucketNamePrefix == "" {
			bucketNamePrefix = strings.Trim(r.Header.Get("X-Bucket-Name"), "/")
		}
	}

	if bucketNamePrefix != "" {
		if key == "" {
			http.Error(w, "missing file key", http.StatusBadRequest)
			return
		}
		key = path.Join(bucketNamePrefix, key)
	}

	resp, err := serv.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(secrets.AWS_BUCKET_NAME),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	// Read up to 512 bytes to sniff content type, then prepend them back to stream.
	sniff := make([]byte, 512)
	n, _ := resp.Body.Read(sniff)
	reader := io.MultiReader(bytes.NewReader(sniff[:n]), resp.Body)

	// Determine Content-Type: prefer object's ContentType if informative, otherwise detect.
	contentType := ""
	if resp.ContentType != nil && *resp.ContentType != "" && *resp.ContentType != "application/octet-stream" {
		contentType = *resp.ContentType
	} else {
		detected := http.DetectContentType(sniff[:n])
		if detected == "application/octet-stream" {
			// fallback to extension-based type
			if ct := mime.TypeByExtension(filepath.Ext(key)); ct != "" {
				detected = ct
			}
		}
		contentType = detected
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}
	w.Header().Set("Content-Type", contentType)

	// Force inline rendering in browsers (use filename from key)
	filename := path.Base(key)
	if filename == "" {
		filename = "file"
	}
	w.Header().Set("Content-Disposition", "inline; filename=\""+filename+"\"")

	if resp.ContentLength != nil {
		w.Header().Set("Content-Length", fmt.Sprint(*resp.ContentLength))
	}
	if resp.ETag != nil {
		w.Header().Set("ETag", *resp.ETag)
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)

	io.Copy(w, reader)
}

func (serv *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	prefix := "/api/v1/cbesuperapp/files/"
	key := r.URL.Path
	if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
		key = key[len(prefix):]
	}

	// If path is /files/<bucket_name>/<key...>, extract bucket from the path.
	bucketNamePrefix := ""
	if key != "" {
		parts := strings.SplitN(strings.Trim(key, "/"), "/", 2)
		if len(parts) >= 2 {
			bucketNamePrefix = strings.Trim(parts[0], "/")
			key = parts[1]
		}
	}

	// If no bucket in path, fall back to query param or header.
	if bucketNamePrefix == "" {
		bucketNamePrefix = strings.Trim(r.URL.Query().Get("bucket_name"), "/")
		if bucketNamePrefix == "" {
			bucketNamePrefix = strings.Trim(r.Header.Get("X-Bucket-Name"), "/")
		}
	}

	contentType := r.Header.Get("Content-Type")
	var bodyReader io.Reader
	var contentLen *int64

	// Support multipart/form-data file uploads (field name "file")
	if strings.HasPrefix(contentType, "multipart/") {
		// limit form size (e.g., 100MB); adjust as needed
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			http.Error(w, "invalid multipart form", http.StatusBadRequest)
			return
		}
		// If client provided bucket_name in the multipart form, prefer it when we don't have one yet
		if bucketNamePrefix == "" {
			if bn := strings.Trim(r.FormValue("bucket_name"), "/"); bn != "" {
				bucketNamePrefix = bn
			}
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "missing form file field 'file'", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// If key not provided in URL, use uploaded filename
		if key == "" {
			key = header.Filename
			if key == "" {
				http.Error(w, "no filename provided", http.StatusBadRequest)
				return
			}
		}

		// Read into memory then create reader so we can set content length
		buf := new(bytes.Buffer)
		n, err := io.Copy(buf, file)
		if err != nil {
			http.Error(w, "failed to read uploaded file", http.StatusInternalServerError)
			return
		}
		bodyReader = bytes.NewReader(buf.Bytes())
		cl := n
		contentLen = &cl

		// Try to get content type from multipart header if present
		if ct := header.Header.Get("Content-Type"); ct != "" {
			contentType = ct
		} else {
			contentType = mime.TypeByExtension(filepath.Ext(key))
		}
	} else {
		// Raw body upload: require key in URL
		if key == "" {
			http.Error(w, "missing file key in URL", http.StatusBadRequest)
			return
		}
		bodyReader = r.Body
		// Content-Length may not be present; leave contentLen nil
		if r.ContentLength > 0 {
			cl := r.ContentLength
			contentLen = &cl
		}
		// If no content type provided, guess from extension
		if contentType == "" || contentType == "application/octet-stream" {
			if ct := mime.TypeByExtension(filepath.Ext(key)); ct != "" {
				contentType = ct
			}
		}
	}

	// If a bucket_name prefix was provided (either via path, multipart form, or fallback), ensure the folder exists and prepend prefix to the key
	if bucketNamePrefix != "" {
		folderKey := bucketNamePrefix + "/"
		// Create a zero-byte object for the "directory" (safe to overwrite)
		_, err := serv.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(secrets.AWS_BUCKET_NAME),
			Key:         aws.String(folderKey),
			Body:        bytes.NewReader([]byte{}),
			ContentType: aws.String("application/x-directory"),
		})
		if err != nil {
			log.Println("failed to create folder marker:", err)
			http.Error(w, "failed to create folder", http.StatusInternalServerError)
			return
		}

		// Prepend prefix to the target key (use posix path join)
		if key == "" {
			// shouldn't happen, but guard
			http.Error(w, "no filename provided", http.StatusBadRequest)
			return
		}
		key = path.Join(bucketNamePrefix, key)
	}

	// Check if the object already exists (to indicate overwrite)
	overwritten := false
	_, headErr := serv.s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(secrets.AWS_BUCKET_NAME),
		Key:    aws.String(key),
	})
	if headErr == nil {
		overwritten = true
	} else {
		// If HeadObject failed, we assume object doesn't exist or cannot be accessed.
		// Log the error for visibility (some S3-compatible services may return different errors).
		log.Printf("head object check: %v (will treat as non-existent)", headErr)
	}

	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(secrets.AWS_BUCKET_NAME),
		Key:         aws.String(key),
		Body:        bodyReader,
		ContentType: aws.String(contentType),
	}
	if contentLen != nil {
		putInput.ContentLength = contentLen
	}

	_, err := serv.s3Client.PutObject(context.TODO(), putInput)
	if err != nil {
		log.Println("upload error:", err)
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	// Build fetch URL: http://<host>/api/v1/cbesuperapp/files/<key>
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}
	fetchPath := "/api/v1/cbesuperapp/files/" + strings.TrimPrefix(key, "/")
	fetchURL := "http://" + host + fetchPath

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message":     "uploaded",
		"key":         key,
		"url":         fetchURL,
		"overwritten": overwritten,
	})
}
