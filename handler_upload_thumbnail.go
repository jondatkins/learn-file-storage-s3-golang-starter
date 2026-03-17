package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	// "thumbnail" should match the HTML form input name
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	// Get the media type from the form file's Content-Type header
	mediaType := header.Header.Get("Content-Type")
	// Read all the image data into a byte slice using io.ReadAll

	imageData, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading form file", err)
		return
	}
	defer file.Close()
	// Get the video's metadata from the SQLite database. The apiConfig's db has a GetVideo method you can use
	// If the authenticated user is not the video owner, return a http.StatusUnauthorized response
	vidIDString := r.PathValue("videoID")
	videoID, err = uuid.Parse(vidIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid video ID", err)
		return
	}
	vidMetaData, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No Video for id", err)
		return
	}
	if vidMetaData.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "User is not video's owner", err)
		return
	}
	data := []byte(imageData)
	imageString := base64.StdEncoding.EncodeToString(data)
	fmt.Println(imageString)

	// Create a data URL with the media type and base64 encoded image data. The format is:
	// data:<media-type>;base64,<data>
	// Store the URL in the thumbnail_url column in the database.
	dataURL := fmt.Sprintf("data:%s;base64,%s", mediaType, imageString)
	vidMetaData.ThumbnailURL = &dataURL
	err = cfg.db.UpdateVideo(vidMetaData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating video", err)
		return
	}
	video, err := cfg.db.GetVideo(videoID)
	vid := database.Video{
		ID:           video.ID,
		CreatedAt:    video.CreatedAt,
		UpdatedAt:    video.UpdatedAt,
		ThumbnailURL: video.ThumbnailURL,
		VideoURL:     video.VideoURL,
	}
	respondWithJSON(w, http.StatusOK, vid)
}
