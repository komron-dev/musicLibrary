package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/komron-dev/musicLibrary/api/models"
	db "github.com/komron-dev/musicLibrary/db/sqlc"
	"github.com/lib/pq"
	"net/http"
	"time"
)

// @Summary Add a new song
// @Description Add a new song to the library by providing its name, release date, text, link, and group name.
// @Accept json
// @Produce json
// @Param request body models.AddSongRequest true "Request body to add a new song"
// @Success 200 {object} db.Song "Successfully added the song"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 403 {object} models.ErrorResponse "Song already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs [post]
func (server *Server) addSong(ctx *gin.Context) {
	var request models.AddSongRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		server.logger.Error("Failed to bind request in addSong(): ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	releaseDate, err := time.Parse("02.01.2006", request.ReleaseDate)
	if err != nil {
		server.logger.Error("Failed to parse time in addSong(): ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response, err := server.store.AddSong(ctx, db.AddSongParams{
		Name:        request.Name,
		ReleaseDate: releaseDate,
		Text:        request.Text,
		Link:        request.Link,
		GroupName:   request.GroupName,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				server.logger.Error("psql error in addSong(): ", err)
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get a song
// @Description Get a specific song by providing the group and song name as query parameters.
// @Produce json
// @Param group query string true "Group name"
// @Param song query string true "Song name"
// @Success 200 {object} db.Song "The requested song"
// @Failure 404 {object} models.ErrorResponse "No such song found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs/info [get]
func (server *Server) getSong(ctx *gin.Context) {
	var request models.GetSongRequest
	request.Song = ctx.Query("song")
	request.Group = ctx.Query("group")

	response, err := server.store.GetSong(ctx, db.GetSongParams{
		GroupName: request.Group,
		Name:      request.Song,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			server.logger.Error("No rows: ", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		server.logger.Error("Failed to add songs: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary List all songs
// @Description Retrieve a list of songs in a paginated manner. Provide page_id and page_size as query parameters.
// @Produce json
// @Param page_id query int false "Page number starting from 1"
// @Param page_size query int false "Number of items per page"
// @Success 200 {array} db.Song "List of songs"
// @Failure 400 {object} models.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs [get]
func (server *Server) listSongs(ctx *gin.Context) {
	var pagination models.Pagination
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		server.logger.Error("Failed to bind query for pagination: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.store.ListSongs(ctx, db.ListSongsParams{
		Limit:  pagination.PageSize,
		Offset: (pagination.PageID - 1) * pagination.PageSize,
	})

	if err != nil {
		server.logger.Error("Failed to list songs: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get song lyrics
// @Description Retrieve the lyrics of a specific song by its UUID. The lyrics are paginated by verses.
// @Produce json
// @Param id path string true "Song UUID"
// @Param page_id query int false "Page number starting from 1"
// @Param page_size query int false "Number of verses per page"
// @Success 200 {array} string "A list of verses"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 404 {object} models.ErrorResponse "No lyrics found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs/{id} [get]
func (server *Server) getSongLyrics(ctx *gin.Context) {
	var request models.GetSongLyricsRequest
	var pagination models.Pagination
	if err := ctx.ShouldBindUri(&request); err != nil {
		server.logger.Error("Failed to bind id from uri in getSongLyrics(): ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		server.logger.Error("Failed to bind query in getSongLyrics(): ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuid, err := uuid.Parse(request.ID)
	if err != nil {
		server.logger.Error("Failed to parse to UUID in getSongLyrics(): ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response, err := server.store.GetSongLyrics(ctx, db.GetSongLyricsParams{
		ID:      uuid,
		Column2: pagination.PageID,
		Column3: pagination.PageSize,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			server.logger.Error("Failed to find rows getSongLyrics(): ", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		server.logger.Error("Failed to get song lyrics: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Update a song
// @Description Update the details of an existing song by providing its ID and new details in the request body.
// @Accept json
// @Produce json
// @Param request body models.UpdateSongRequest true "Request body with updated song details"
// @Success 200 {object} db.Song "Successfully updated the song"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 404 {object} models.ErrorResponse "No such song found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs [put]
func (server *Server) updateSong(ctx *gin.Context) {
	var request models.UpdateSongRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuid, err := uuid.Parse(request.Id)
	if err != nil {
		server.logger.Error("Failed to parse to UUID in getSongLyrics(): ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	releaseDate, err := time.Parse("02.01.2006", request.ReleaseDate)
	response, err := server.store.UpdateSong(ctx, db.UpdateSongParams{
		ID:          uuid,
		Name:        request.Name,
		ReleaseDate: releaseDate,
		Text:        request.Text,
		Link:        request.Link,
		GroupName:   request.GroupName,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Delete a song
// @Description Delete an existing song by providing its ID as a path parameter.
// @Produce json
// @Param id path string true "Song UUID"
// @Success 200 {object} nil "Successfully deleted the song"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 404 {object} models.ErrorResponse "No such song found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs/{id} [delete]
func (server *Server) deleteSong(ctx *gin.Context) {
	var request models.DeleteSongRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	uuid, err := uuid.Parse(request.ID)
	if err != nil {
		server.logger.Error("Failed to parse to UUID in deleteSong(): ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteSong(ctx, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
