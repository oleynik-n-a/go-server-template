package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	tokensTable = "refresh_tokens"
	toolName = "garbage_collector"

	errMsgRevokedTokens = "Error collecting revoked tokens"
	errMsgExpiredTokens = "Error collecting expired tokens"
)

type GarbageCollector struct {
	db *sqlx.DB
	logger *zap.Logger
	startTime time.Time
	period time.Duration
}

func NewGarbageCollector(db *sqlx.DB, logger *zap.Logger, startTime time.Time, period time.Duration) *GarbageCollector {
	return &GarbageCollector{db: db, logger: logger, startTime: startTime, period: period}
}

func (gc *GarbageCollector) Run() {
	go func() {
		time.Sleep(time.Until(gc.startTime))

		for {
			time.Sleep(gc.period)

			revokedTokensErased, err := gc.eraseRevokedTokens()
			if err != nil {
				gc.logger.Error(errMsgRevokedTokens, zap.Namespace(toolName), zap.Error(err))
				continue
			}
			gc.logger.Debug("Erased revoked tokens", zap.Int64("value", revokedTokensErased))

			expiredTokensErased, err := gc.eraseExpiredTokens()
			if err != nil {
				gc.logger.Error(errMsgExpiredTokens, zap.Namespace(toolName), zap.Error(err))
				continue
			}
			gc.logger.Debug("Erased expired tokens", zap.Int64("value", expiredTokensErased))

			gc.logger.Debug("Garbage collected successfully", zap.Namespace(toolName))
		}
	}()
}

func (gc *GarbageCollector) eraseRevokedTokens() (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE revoked=TRUE;", tokensTable)
	result, err := gc.db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (gc *GarbageCollector) eraseExpiredTokens() (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE expires_at<$1;", tokensTable)
	result, err := gc.db.Exec(query, time.Now())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
