package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/emarifer/gocms/internal/admin_app/api/dto"
	"github.com/emarifer/gocms/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a *API) getImageHandler(c *gin.Context) {
	ctx := c.Request.Context()
	imageBinding := &dto.ImageBinding{}

	// localhost:8080/image/{uuid}
	if err := c.ShouldBindUri(imageBinding); err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"could not get image uuid",
		)
		c.Error(customError)

		return
	}

	// Get image metadata by UUID

	imageMetadata, err := a.serv.RecoverImageMetadata(ctx, imageBinding.UUID)
	if err != nil {
		if re.MatchString(err.Error()) {
			customError := NewCustomError(
				http.StatusInternalServerError,
				err.Error(),
				"An unexpected condition was encountered.",
			)
			c.Error(customError)

			return
		}

		if strings.Contains(err.Error(), "no rows in result set") {
			customError := NewCustomError(
				http.StatusNotFound,
				err.Error(),
				"The requested resource could not be found but may be available again in the future.",
			)
			c.Error(customError)

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid": imageMetadata.UUID,
		"name": imageMetadata.Name,
		"alt":  imageMetadata.Alt,
	})
}

func (a *API) addImageHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10*1000000)
	ctx := c.Request.Context()

	mForm, err := c.MultipartForm()
	if err != nil {
		customError := NewCustomError(
			http.StatusInternalServerError,
			err.Error(),
			"could not create multipart form",
		)
		c.Error(customError)

		return
	}

	alt_text_array := mForm.Value["alt"]
	alt_text := "unknown"
	if len(alt_text_array) > 0 {
		alt_text = alt_text_array[0]
	}

	/* addImageBinding := &dto.AddImageRequest{}

	if err := c.ShouldBind(addImageBinding); err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"invalid request body",
		)
		c.Error(customError)

		return
	}

	alt_text := addImageBinding.Alt */

	// Begin saving the file to the filesystem
	file_array := mForm.File["file"]
	if len(file_array) == 0 {
		customError := NewCustomError(
			http.StatusBadRequest,
			"error: could not get the file array",
			"could not get the file array",
		)
		c.Error(customError)

		return
	}
	file := file_array[0]
	if file == nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			"error: could not upload file",
			"could not upload file",
		)
		c.Error(customError)

		return
	}
	allowed_types := []string{"image/jpeg", "image/png", "image/gif"}
	if file_content_type := file.Header.Get("content-type"); !slices.Contains(
		allowed_types, file_content_type,
	) {
		customError := NewCustomError(
			http.StatusBadRequest,
			"error: file type not supported",
			"file type not supported",
		)
		c.Error(customError)

		return
	}
	allowed_extensions := []string{".jpeg", ".jpg", ".png"}
	ext := filepath.Ext(file.Filename)
	// check ext is supported
	if ext != "" && !slices.Contains(allowed_extensions, ext) {
		customError := NewCustomError(
			http.StatusBadRequest,
			"error: file extension is not supported",
			"file extension is not supported",
		)
		c.Error(customError)

		return
	}

	uuid := uuid.New()
	filename := fmt.Sprintf("%s%s", uuid.String(), ext)
	rootPath, _ := os.Getwd()
	image_path := filepath.Join(
		rootPath, "assets/images_files", filename,
	)
	err = c.SaveUploadedFile(file, image_path)
	if err != nil {
		customError := NewCustomError(
			http.StatusInternalServerError,
			err.Error(),
			"could not save file",
		)
		c.Error(customError)

		return
	}
	// End saving to filesystem

	// Save metadata into the DB
	imageData := &entity.Image{
		UUID: uuid,
		Name: filename,
		Alt:  alt_text,
	}

	err = a.serv.CreateImageMetadata(ctx, imageData)
	if err != nil {
		customError := NewCustomError(
			http.StatusInternalServerError,
			err.Error(),
			"could not add image metadata to db",
		)
		c.Error(customError)

		err = os.Remove(image_path)
		if err != nil {
			customError := NewCustomError(
				http.StatusInternalServerError,
				err.Error(),
				"could not remove image",
			)
			c.Error(customError)
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": uuid.String(),
	})
}

func (a *API) deleteImageHandler(c *gin.Context) {
	ctx := c.Request.Context()
	imageBinding := &dto.ImageBinding{}

	// localhost:8080/image/{uuid}
	if err := c.ShouldBindUri(imageBinding); err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"could not get image uuid",
		)
		c.Error(customError)

		return
	}

	// Get image metadata by UUID

	imageMetadata, err := a.serv.RecoverImageMetadata(ctx, imageBinding.UUID)
	if err != nil {
		if re.MatchString(err.Error()) {
			customError := NewCustomError(
				http.StatusInternalServerError,
				err.Error(),
				"An unexpected condition was encountered.",
			)
			c.Error(customError)

			return
		}

		if strings.Contains(err.Error(), "no rows in result set") {
			customError := NewCustomError(
				http.StatusNotFound,
				err.Error(),
				"The requested resource could not be found but may be available again in the future.",
			)
			c.Error(customError)

			return
		}
	}

	// Remove image metadata by uuid

	rowsAffected, err := a.serv.RemoveImage(ctx, imageBinding.UUID)
	if err != nil {
		if re.MatchString(err.Error()) {
			customError := NewCustomError(
				http.StatusInternalServerError,
				err.Error(),
				"An unexpected condition was encountered",
			)
			c.Error(customError)

			return
		}
	}

	if rowsAffected == -1 {
		customError := NewCustomError(
			http.StatusNotFound,
			"no rows in result set",
			"could not delete image metadata",
		)
		c.Error(customError)

		return
	}

	// Delete file system image by its path
	rootPath, _ := os.Getwd()
	image_path := filepath.Join(
		rootPath, "assets/images_files", imageMetadata.Name,
	)
	err = os.Remove(image_path)
	if err != nil {
		customError := NewCustomError(
			http.StatusInternalServerError,
			err.Error(),
			"could not remove image file",
		)
		c.Error(customError)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rowsAffected": rowsAffected,
	})
}

/* REFERENCES:
https://www.freeformatter.com/json-escape.html#before-output
https://www.youtube.com/watch?v=yqDYYjbatfE
https://github.com/sikozonpc/fullstack-go-htmx
https://stackoverflow.com/questions/9722603/storing-image-in-database-directly-or-as-base64-data

https://tutorialedge.net/golang/go-file-upload-tutorial/

https://www.google.com/search?q=golang+file+upload+rest+api&sca_esv=dcb7b8806130e33c&sxsrf=ACQVn08616pDvbQ3cGeZyya8IskU-hxACA%3A1709290929268&ei=sbXhZezxD4yX9u8Ps6CtuAI&oq=golang+file+u&gs_lp=Egxnd3Mtd2l6LXNlcnAiDWdvbGFuZyBmaWxlIHUqAggBMgwQIxiABBiKBRgTGCcyCBAAGIAEGMsBMggQABiABBjLATIIEAAYgAQYywEyCBAAGIAEGMsBMggQABiABBjLATIGEAAYFhgeMgYQABgWGB4yBhAAGBYYHjIGEAAYFhgeSPVIUNkMWKsScAF4AZABAJgBqAGgAcYEqgEDMi4zuAEByAEA-AEBmAIGoAL8BMICChAAGEcY1gQYsAOYAwCIBgGQBgiSBwMyLjQ&sclient=gws-wiz-serp
*/
