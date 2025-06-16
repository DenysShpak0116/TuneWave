package comment

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CommentHandler struct {
	CommentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
	}
}

type CreateCommentRequest struct {
	SongID  string `json:"songId"`
	UserID  string `json:"userId"`
	Content string `json:"content"`
}

// Create godoc
// @Summary Create a new comment
// @Description Creates a new comment for a song. Returns the created comment object.
// @Tags comments
// @Security     BearerAuth
// @Accept  json
// @Produce  json
// @Param comment body CreateCommentRequest true "Comment creation data"
// @Router /comments [post]
func (ch *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) error {
	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid request")
	}

	songUUID, err := uuid.Parse(req.SongID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid song ID")
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid user ID")
	}

	comment := &models.Comment{
		SongID:  songUUID,
		UserID:  userUUID,
		Content: req.Content,
	}

	if err := ch.CommentService.Create(r.Context(), comment); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to create comment")
	}

	commentWithPreload, err := ch.CommentService.GetByID(r.Context(), comment.ID, "User", "User.Followers")
	if err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to get comment")
	}

	dtoBuilder := dto.NewDTOBuilder()
	commentDTO := dtoBuilder.BuildCommentDTO(commentWithPreload)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, commentDTO)
	return nil
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Deletes a comment by its ID. Returns no content on success.
// @Tags comments
// @Security     BearerAuth
// @Param id path string true "Comment ID"
// @Router /comments/{id} [delete]
func (ch *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) error {
	commentID := chi.URLParam(r, "id")
	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		return helpers.NewAPIError(http.StatusBadRequest, "invalid comment ID")
	}

	if err := ch.CommentService.Delete(r.Context(), commentUUID); err != nil {
		return helpers.NewAPIError(http.StatusInternalServerError, "failed to delete comment")
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
	return nil
}
