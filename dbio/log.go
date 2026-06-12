//go:build !private_logger

package dbio

import "github.com/kdnetwork/code-snippet/go/log"

func InitLogger() {
	log.InitDefaultLogger()
}
