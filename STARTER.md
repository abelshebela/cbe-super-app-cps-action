# Sample Code and Descriptions
## Minio Fileupload
```
package main

import (
	"cbe-super-app-api-starter/pkg/config"
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {

	ctx := context.Background()
	env, err := config.Load()
	if err != nil {
		log.Warn().Msg("Failed to Load .env file!")
		return
	}
	minio, err := config.NewMinioClient(env)
	if err != nil {
		log.Fatal().Msgf("Faild: %v", err)
		return
	}

	file_path := "./cmd/test.txt"
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		fmt.Printf("File does not exist at path: %s\n", filePath)
		return
	}
	ui, err := minio.SaveObject(ctx, config.SaveObjectBody{
		BucketName:  "cbe-user-profile",
		ObjectName:  "test.txt",
		File:        file_path,
		ContentType: config.ContentTypePlainText,
	})

	if err != nil {
		log.Info().Msg("Error uploading file")
		return
	}

	fmt.Println("succesfully file uploaded: ", ui)


}
```

