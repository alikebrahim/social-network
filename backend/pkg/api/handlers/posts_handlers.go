package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"socialNetwork/pkg/api/utils"
	"socialNetwork/pkg/domain/notifications"
	"socialNetwork/pkg/domain/posts"
	"socialNetwork/pkg/errors"
	"socialNetwork/pkg/logger"
	"socialNetwork/pkg/storage"
)

type PostsHandler struct {
	store  storage.Storage
	logger logger.Logger
}

func NewPostsHandler(store storage.Storage) *PostsHandler {
	return &PostsHandler{
		store:  store,
		logger: logger.GetLogger("api/handlers/posts"),
	}
}

func (h *PostsHandler) HandlePostCreate(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Parse request body
	var post posts.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}

	// Validate post content
	if post.Content == "" && post.Image == "" {
		log.Warn("Invalid post: empty content and image")
		return errors.ErrInvalidInput
	}
	
	// Validate image format if provided
	if post.Image != "" {
		isValidFormat := false
		validFormats := []string{".jpg", ".jpeg", ".png", ".gif"}
		
		// Simple extension check - in a real app, you would validate the actual file contents
		for _, format := range validFormats {
			if len(post.Image) > len(format) && post.Image[len(post.Image)-len(format):] == format {
				isValidFormat = true
				break
			}
		}
		
		if !isValidFormat {
			log.Warn("Invalid image format", "image", post.Image)
			return utils.WriteJson(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid image format. Supported formats: JPEG, PNG, GIF",
			})
		}
	}
	
	// Validate privacy level
	if post.PrivacyLevel != "" {
		if post.PrivacyLevel != posts.PrivacyPublic && 
		   post.PrivacyLevel != posts.PrivacyAlmostPrivate && 
		   post.PrivacyLevel != posts.PrivacyPrivate {
			log.Warn("Invalid privacy level", "privacy_level", post.PrivacyLevel)
			return utils.WriteJson(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid privacy level. Supported levels: public, almost_private, private",
			})
		}
	} else {
		// Default to public if not specified
		post.PrivacyLevel = posts.PrivacyPublic
	}

	// Get user ID from session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	userID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Set the user ID and creation time
	post.UserID = userID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	// Create the post
	postID, err := h.store.CreatePost(post)
	if err != nil {
		log.Error("Failed to create post", "error", err)
		return err
	}

	log.Info("Post created", "post_id", postID, "user_id", userID)

	// Return the created post ID
	return utils.WriteJson(w, http.StatusCreated, map[string]int64{
		"post_id": postID,
	})
}

func (h *PostsHandler) HandlePostsGet(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get post ID from the URL
	idStr := r.PathValue("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid post ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get user ID from session for visibility checks
	var userID int64
	cookie, err := r.Cookie("session_token")
	if err == nil {
		userID, err = h.store.GetUserIdBySession(cookie.Value)
		if err != nil {
			log.Warn("Failed to get user ID from session", "error", err)
			// Continue with userID = 0 (anonymous)
		}
	}

	// Get the post
	retrievedPosts, err := h.store.GetPostByID(userID, postID)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Warn("Post not found", "post_id", postID)
			return errors.ErrNotFound
		}
		log.Error("Failed to get post", "error", err)
		return err
	}

	if len(retrievedPosts) == 0 {
		log.Warn("Post not found", "post_id", postID)
		return errors.ErrNotFound
	}

	log.Info("Post retrieved", "post_id", postID, "user_id", userID)

	return utils.WriteJson(w, http.StatusOK, retrievedPosts[0])
}

func (h *PostsHandler) HandlePostEdit(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get post ID from the URL
	idStr := r.PathValue("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid post ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Parse request body
	var post posts.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}

	// Validate post content
	if post.Content == "" && post.Image == "" {
		log.Warn("Invalid post update: empty content and image")
		return errors.ErrInvalidInput
	}
	
	// Validate image format if provided
	if post.Image != "" {
		isValidFormat := false
		validFormats := []string{".jpg", ".jpeg", ".png", ".gif"}
		
		// Simple extension check - in a real app, you would validate the actual file contents
		for _, format := range validFormats {
			if len(post.Image) > len(format) && post.Image[len(post.Image)-len(format):] == format {
				isValidFormat = true
				break
			}
		}
		
		if !isValidFormat {
			log.Warn("Invalid image format", "image", post.Image)
			return utils.WriteJson(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid image format. Supported formats: JPEG, PNG, GIF",
			})
		}
	}
	
	// Validate privacy level if provided
	if post.PrivacyLevel != "" {
		if post.PrivacyLevel != posts.PrivacyPublic && 
		   post.PrivacyLevel != posts.PrivacyAlmostPrivate && 
		   post.PrivacyLevel != posts.PrivacyPrivate {
			log.Warn("Invalid privacy level", "privacy_level", post.PrivacyLevel)
			return utils.WriteJson(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid privacy level. Supported levels: public, almost_private, private",
			})
		}
	}

	// Get user ID from session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	userID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Check if user is the post owner
	isOwner, err := h.store.IsPostOwner(userID, postID)
	if err != nil {
		log.Error("Failed to check post ownership", "error", err)
		return err
	}

	if !isOwner {
		log.Warn("Unauthorized post edit attempt", "post_id", postID, "user_id", userID)
		return errors.ErrUnauthorized
	}

	// Update the post
	post.ID = postID
	post.UpdatedAt = time.Now()
	err = h.store.EditPost(post)
	if err != nil {
		log.Error("Failed to update post", "error", err)
		return err
	}

	log.Info("Post updated", "post_id", postID, "user_id", userID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Post updated successfully",
	})
}

func (h *PostsHandler) HandlePostDelete(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get post ID from the URL
	idStr := r.PathValue("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid post ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get user ID from session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	userID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Check if user is the post owner
	isOwner, err := h.store.IsPostOwner(userID, postID)
	if err != nil {
		log.Error("Failed to check post ownership", "error", err)
		return err
	}

	if !isOwner {
		log.Warn("Unauthorized post delete attempt", "post_id", postID, "user_id", userID)
		return errors.ErrUnauthorized
	}

	// Delete the post
	err = h.store.DeletePost(postID)
	if err != nil {
		log.Error("Failed to delete post", "error", err)
		return err
	}

	log.Info("Post deleted", "post_id", postID, "user_id", userID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Post deleted successfully",
	})
}

func (h *PostsHandler) HandlePostComment(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get post ID from the URL
	idStr := r.PathValue("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid post ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Parse request body
	var comment posts.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return errors.ErrBadRequest
	}

	// Validate comment content
	if comment.Content == "" {
		log.Warn("Invalid comment: empty content")
		return errors.ErrInvalidInput
	}

	// Get user ID from session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	userID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Check if user can see the post
	canSee, err := h.store.CanUserSeePost(userID, postID)
	if err != nil {
		log.Error("Failed to check post visibility", "error", err)
		return err
	}

	if !canSee {
		log.Warn("Unauthorized comment attempt", "post_id", postID, "user_id", userID)
		return errors.ErrNotFound
	}

	// Set comment properties
	comment.PostID = postID
	comment.UserID = userID
	comment.CreatedAt = time.Now()

	// Create the comment
	err = h.store.CreateComment(comment)
	if err != nil {
		log.Error("Failed to create comment", "error", err)
		return err
	}

	log.Info("Comment created", "post_id", postID, "user_id", userID)
	
	// Create a notification for the post owner if the commenter is not the owner
	// Get the post owner ID
	var postsList []posts.Post
	postsList, err = h.store.GetPostByID(postID, userID)
	if err == nil && len(postsList) > 0 {
		postOwnerID := postsList[0].UserID
		
		// Don't notify if you're commenting on your own post
		if postOwnerID != userID {
			// Create notification
			notification := &notifications.Notification{
				UserID:    postOwnerID,
				Type:      notifications.TypeComment,
				Content:   "Someone commented on your post",
				RelatedID: postID,
				SenderID:  userID,
				CreatedAt: time.Now(),
			}
			
			err = h.store.CreateNotification(notification)
			if err != nil {
				// Log but don't return the error - comment was still created successfully
				log.Error("Failed to create notification for post comment", "error", err)
			} else {
				log.Debug("Comment notification created", "post_id", postID, "owner_id", postOwnerID)
			}
		}
	}

	return utils.WriteJson(w, http.StatusCreated, map[string]string{
		"message": "Comment added successfully",
	})
}

func (h *PostsHandler) HandlePostLike(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get post ID from the URL
	idStr := r.PathValue("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid post ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get user ID from session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	userID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Check if user can see the post
	canSee, err := h.store.CanUserSeePost(userID, postID)
	if err != nil {
		log.Error("Failed to check post visibility", "error", err)
		return err
	}

	if !canSee {
		log.Warn("Unauthorized like attempt", "post_id", postID, "user_id", userID)
		return errors.ErrNotFound
	}

	// Create the like
	like := posts.Like{
		PostID: postID,
		UserID: userID,
	}

	err = h.store.CreateLike(like)
	if err != nil {
		if errors.IsConflict(err) {
			log.Info("Post already liked", "post_id", postID, "user_id", userID)
			return utils.WriteJson(w, http.StatusOK, map[string]string{
				"message": "Post already liked",
			})
		}
		log.Error("Failed to like post", "error", err)
		return err
	}

	log.Info("Post liked", "post_id", postID, "user_id", userID)
	
	// Create a notification for the post owner if the liker is not the owner
	// Get the post owner ID
	var postsList []posts.Post
	postsList, err = h.store.GetPostByID(postID, userID)
	if err == nil && len(postsList) > 0 {
		postOwnerID := postsList[0].UserID
		
		// Don't notify if you're liking your own post
		if postOwnerID != userID {
			// Create notification
			notification := &notifications.Notification{
				UserID:    postOwnerID,
				Type:      notifications.TypeLike,
				Content:   "Someone liked your post",
				RelatedID: postID,
				SenderID:  userID,
				CreatedAt: time.Now(),
			}
			
			err = h.store.CreateNotification(notification)
			if err != nil {
				// Log but don't return the error - like was still created successfully
				log.Error("Failed to create notification for post like", "error", err)
			} else {
				log.Debug("Like notification created", "post_id", postID, "owner_id", postOwnerID)
			}
		}
	}

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Post liked successfully",
	})
}

func (h *PostsHandler) HandlePostUnlike(w http.ResponseWriter, r *http.Request) error {
	// Get logger from request context
	log := logger.FromRequest(r)

	// Get post ID from the URL
	idStr := r.PathValue("id")
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid post ID", "id", idStr, "error", err)
		return errors.ErrBadRequest
	}

	// Get user ID from session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Error("Failed to get session token", "error", err)
		return errors.ErrUnauthorized
	}

	userID, err := h.store.GetUserIdBySession(cookie.Value)
	if err != nil {
		log.Error("Failed to get user ID from session", "error", err)
		return errors.ErrUnauthorized
	}

	// Remove the like
	like := posts.Like{
		PostID: postID,
		UserID: userID,
	}

	err = h.store.RemoveLikes(like)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Info("Post not liked", "post_id", postID, "user_id", userID)
			return utils.WriteJson(w, http.StatusOK, map[string]string{
				"message": "Post was not liked",
			})
		}
		log.Error("Failed to unlike post", "error", err)
		return err
	}

	log.Info("Post unliked", "post_id", postID, "user_id", userID)

	return utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Post unliked successfully",
	})
}