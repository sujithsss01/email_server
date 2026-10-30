package emailmodule

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type EmailModule struct {
	handler *emailHandler
}

func (*EmailModule) RegisterRoutes(router *gin.RouterGroup) {

}
func NewEmailModule(session *gocql.Session) *EmailModule {
	service := newEmailService(session)
	handler := newEmailHandler(service)
	module := &EmailModule{
		handler: handler,
	}
	return module
}
