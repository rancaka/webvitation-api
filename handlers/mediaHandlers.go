package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo"
	"github.com/rancaka/webvitation-api/db"
)

func HandleUploadMedia(c echo.Context) error {

	auth := GetAuth(c)
	if auth == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	eventID := c.Param("eventID")
	event, err := db.GetEventByEventID(c.Request().Context(), eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if event.CreatorUID != auth.UID {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	reader, err := c.Request().MultipartReader()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var part *multipart.Part
	var mediaType db.MediaType
	var mediaPath string
	var videoPosterPath string

	for err == nil {
		part, err = reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		fieldName := part.FormName()
		if fieldName == "media" {

			b, err := ioutil.ReadAll(part)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			mimeType := http.DetectContentType(b[:512])
			mimeTypeSplit := strings.Split(mimeType, "/")
			mediaType = db.MediaType(mimeTypeSplit[0])
			mediaPath = "events/" + eventID + "/" + fmt.Sprint(time.Now().Unix()) + "-" + part.FileName()

			if mediaType == "image" {
				err = UploadMedia(c.Request().Context(), bytes.NewReader(b), mediaPath)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			} else if mediaType == "video" {
				// TODO: handle video
				videoPosterPath = ""
			} else if mediaType == "audio" {
				// TODO: handle audio

			} else {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("file is not supported: %v", mimeType))
			}

			break
		}
	}

	err = db.InsertMedia(c.Request().Context(), eventID, mediaPath, videoPosterPath, mediaType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"MediaPath": mediaPath,
	})
}

func UploadMedia(ctx context.Context, imageReader io.Reader, path string) error {

	image, err := imaging.Decode(imageReader)
	if err != nil {
		return err
	}

	maxWidth := 1080
	maxHeight := 1350
	width := image.Bounds().Max.X
	height := image.Bounds().Max.Y
	aspectRatio := float32(width) / float32(height)
	if aspectRatio > 1 {
		// handle landscape
		if width > maxWidth {
			image = imaging.Resize(image, maxWidth, 0, imaging.Lanczos)
		}
	} else if aspectRatio < 1 {
		// handle potrait
		if height > maxHeight {
			image = imaging.Resize(image, 0, maxHeight, imaging.Lanczos)
		}
	} else {
		// handle square
		if width > maxWidth {
			image = imaging.Resize(image, maxWidth, maxWidth, imaging.Lanczos)
		}
	}

	bucket, err := fb.Storage.Bucket(config.Storage.Bucket.Name)
	if err != nil {
		return err
	}

	format, err := imaging.FormatFromFilename(path)
	if err != nil {
		return err
	}

	object := bucket.Object(path)
	wc := object.NewWriter(ctx)
	err = imaging.Encode(wc, image, format)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}
