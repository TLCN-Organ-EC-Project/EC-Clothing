package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Renew Access Token
// @ID renewAccessToken
// @Param X-Refresh-Token header string true "X-Refresh-Token"
// @Tags Auth
// @Success 200 {string} succesfully
// @Failure 401 {string} error
// @Failure 403 {string} error
// @Failure 404 {string} error
// @Failure 500 {string} error
// @Router /api/tokens/renew_access [post]
func (server *Server) renewAccessToken(ctx *gin.Context) {
	// Get the refresh token from the header
	refreshToken := ctx.Request.Header.Get("X-Refresh-Token")

	// Validate the refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Get the session from the database
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != refreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Set the access token in the header
	ctx.Header("X-Access-Token", accessToken)
	ctx.Header("X-Access-Token-Expired-At", accessPayload.ExpiredAt.String())
	// Send the response
	ctx.JSON(http.StatusOK, "successfully")
}
