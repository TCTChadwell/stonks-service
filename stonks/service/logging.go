package service

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

var initOnce sync.Once

func GetLog(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	initOnce.Do(func() {
		InitLogger()
	})

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {

		metaItem, ok := md["request-id"]
		if ok && len(metaItem) > 0 {
			fields["request-id"] = md["request-id"][0]
		}
	}

	return logrus.WithFields(fields)
}

func InitLogger() {

	logrus.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	logrus.SetLevel(logrus.DebugLevel)
}
