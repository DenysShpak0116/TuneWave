package comment

import (
	"encoding/json"
	"net/http"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
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
// @Tags Comments
// @Security     BearerAuth
// @Accept  json
// @Produce  json
// @Param comment body CreateCommentRequest true "Comment creation data"
// @Router /comments [post]
func (ch *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}

	songUUID, err := uuid.Parse(req.SongID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song ID", err)
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	comment := &models.Comment{
		SongID:  songUUID,
		UserID:  userUUID,
		Content: req.Content,
	}

	if err := ch.CommentService.Create(r.Context(), comment); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to create comment", err)
		return
	}

	commentWithPreload, err := ch.CommentService.GetByID(r.Context(), comment.ID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get comment", err)
		return
	}

	commentDTO := &dtos.CommentDTO{
		ID:         commentWithPreload.ID,
		AuthorID:   commentWithPreload.User.ID,
		AuthorName: commentWithPreload.User.Username,
		Content:    commentWithPreload.Content,
		CreatedAt:  commentWithPreload.CreatedAt,
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, commentDTO)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Deletes a comment by its ID. Returns no content on success.
// @Tags Comments
// @Security     BearerAuth
// @Param id path string true "Comment ID"
// @Router /comments/{id} [delete]
func (ch *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	commentUUID, err := uuid.Parse(commentID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid comment ID", err)
		return
	}

	if err := ch.CommentService.Delete(r.Context(), commentUUID); err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to delete comment", err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}
