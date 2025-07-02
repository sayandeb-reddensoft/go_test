package globals

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nelsonin-research-org/cdc-auth/models/appschema"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var RelationalDb *gorm.DB
var RedisClient *redis.Client
var AppKeys appschema.CertificateKeys
var RequestStore appschema.RequestStore
var AWSSesSession *session.Session