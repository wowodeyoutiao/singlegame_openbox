package db

import (
	"context"
	"fmt"

	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

// TxFunc 在事务中执行的业务逻辑函数
type TxFunc func(ctx context.Context, session *xorm.Session) error

// WithTx 在事务中执行给定的函数，自动提交或回滚
func WithTx(ctx context.Context, fn TxFunc) error {
	if fn == nil {
		return fmt.Errorf("transaction function cannot be nil")
	}

	// 检查上下文是否已取消
	// if err := ctx.Err(); err != nil {
	// 	return err
	// }

	session := database.NewSession()
	defer session.Close()

	// 开始事务
	if err := session.Begin(); err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 执行业务逻辑
	err := fn(ctx, session)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.WithError(rollbackErr).Error("failed to rollback transaction")
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rollbackErr)
		}
		return err
	}

	// 提交事务
	if err := session.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TxFuncWithResult 在事务中执行并返回结果的业务逻辑函数
type TxFuncWithResult[T any] func(ctx context.Context, session *xorm.Session) (T, error)

// WithTxResult 在事务中执行给定的函数并返回结果
func WithTxResult[T any](ctx context.Context, fn TxFuncWithResult[T]) (T, error) {
	var result T
	err := WithTx(ctx, func(ctx context.Context, session *xorm.Session) error {
		var err error
		result, err = fn(ctx, session)
		return err
	})
	return result, err
}
