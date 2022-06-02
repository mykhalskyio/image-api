package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mykhalskyio/image-api/internal/entity"
)

// image service intrrface
type ImageService interface {
	Upload(imageOriginal *entity.ImageUpload) error
	Download(id int, quality int) (string, error)
	Delete(id int) error
}

// image controller struct
type ImageController struct {
	service ImageService
}

// new image controller
func newImageController(imgService ImageService) *ImageController {
	return &ImageController{imgService}
}

// upload image
func (img *ImageController) upload(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	image, err := file.Open()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer image.Close()
	imageBytes := make([]byte, file.Size)
	_, err = image.Read(imageBytes)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = img.service.Upload(&entity.ImageUpload{Image: imageBytes, ImageName: file.Filename})
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusAccepted)
}

// download image
func (img *ImageController) download(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	quality, err := strconv.Atoi(ctx.Query("quality"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if id < 1 {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id less 1"))
		return
	}
	if quality != 100 && quality != 75 && quality != 50 && quality != 25 {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("there is no such quality"))
		return
	}

	imagePath, err := img.service.Download(id, quality)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.File(imagePath)
	ctx.Status(http.StatusOK)
}

// delete image
func (img *ImageController) delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if id < 1 {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id less 1"))
		return
	}
	err = img.service.Delete(id)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.Status(204)
}
