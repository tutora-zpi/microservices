package file

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/metadata"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

type localFileService struct {
	serverPrefixPath string
}

// Save implements interfaces.FileService.
func (l *localFileService) Save(ctx context.Context, file *metadata.FileMetadata) (string, error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("save canceled")
	default:
		generalError := fmt.Errorf("failed to save file")
		if err := os.MkdirAll("./media", os.ModePerm); err != nil {
			log.Printf("Failed to create dir: %v", err)
			return "", generalError
		}

		name := file.GenerateUniqueFileName()

		p := path.Join("./media", name)

		osFile, err := os.Create(p)
		if err != nil {
			log.Printf("Failed to create file with name %s", name)
			return "", generalError
		}

		defer osFile.Close()

		bytes, err := io.ReadAll(file.File)
		if err != nil {
			log.Printf("Failed to read file: %v", err)
			return "", generalError
		}

		writtenBytes, err := osFile.Write(bytes)
		if err != nil {
			log.Printf("Failed to write bytes: %v", err)
			return "", generalError
		}

		log.Printf("Successfully saved new file: %s, written bytes %d", name, writtenBytes)

		url := path.Join(l.serverPrefixPath, name)
		log.Printf("Saved under: %s", url)

		return url, nil
	}
}

func NewLocalFileService(prefixPath string) interfaces.FileService {
	return &localFileService{serverPrefixPath: prefixPath}
}
