package api

import (
	server "rag-server"
	"rag-server/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// RegisterRoutes returns a server.Registrar that registers all API routes.
// It wires user, node, and knowledge base handlers under /api.
func RegisterRoutes(conn *pgx.Conn, repoProxy string) server.Registrar {
	return func(r *gin.Engine) {
		api := r.Group("/api")

		// Apply internal service authentication
		api.Use(auth.InternalAuthMiddleware())

		registerUserRoutes(api)
		registerNodeRoutes(api.Group("/agent"))
		registerKnowledgeRoutes(api, conn, repoProxy)
		registerRAGRoutes(api)
		registerAskAIRoutes(api)
		registerAdminSettingRoutes(api)
	}
}
