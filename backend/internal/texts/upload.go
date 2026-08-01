package texts

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/NickFinchD/chinese-learning-api/internal/response"
	"github.com/gin-gonic/gin"
)

const (
	uploadDir     = "uploads/texts"
	maxUploadSize = 5 << 20 // 5 MiB
)

var allowedImageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

// AdminUploadImage accepts a single image file and stores it on disk under
// uploadDir, returning a URL the admin UI can use directly as a text's
// image_url (and that's later served back by the static /uploads route
// registered in cmd/server). Upload is a separate step from create/update so
// the image can be picked before the text itself has an ID.
func (h *Handler) AdminUploadImage(c *gin.Context) {

	file, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "no image file provided")
		return
	}

	if file.Size > maxUploadSize {
		response.BadRequest(c, "image is too large (max 5MB)")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedImageExts[ext] {
		response.BadRequest(c, "unsupported image type")
		return
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		response.Internal(c)
		return
	}

	name, err := randomFilename(ext)
	if err != nil {
		response.Internal(c)
		return
	}

	if err := c.SaveUploadedFile(file, filepath.Join(uploadDir, name)); err != nil {
		response.Internal(c)
		return
	}

	response.JSON(c, http.StatusCreated, gin.H{
		"url": "/uploads/texts/" + name,
	})
}

func randomFilename(ext string) (string, error) {
	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf) + ext, nil
}
