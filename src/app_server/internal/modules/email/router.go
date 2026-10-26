package emailmodule

import (
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/gin-gonic/gin"
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
