package db

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试用的User模型
type User struct {
	ID   int64  `xorm:"pk autoincr"`
	Name string `xorm:"varchar(50)"`
	Age  int    `xorm:"int"`
}

// TestMain 设置测试环境
func TestMain(m *testing.M) {
	// 初始化数据库连接
	database, err := xorm.NewEngine("mysql", "root:mysql@/test?charset=utf8")
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// 设置全局数据库实例用于事务函数
	SetDatabase(database)

	// 创建测试表
	if err := database.Sync2(new(User)); err != nil {
		panic(fmt.Sprintf("failed to sync table: %v", err))
	}

	// 运行测试
	m.Run()
}

// 清理测试数据
func cleanupUsers(t *testing.T) {
	_, err := database.Where("1=1").Delete(&User{})
	require.NoError(t, err)
}

func TestWithTx_NilFunction(t *testing.T) {
	ctx := context.Background()
	err := WithTx(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction function cannot be nil")
}

func TestWithTx_Success(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	// 成功的事务
	err := WithTx(ctx, func(ctx context.Context, session *xorm.Session) error {
		user := &User{Name: "Alice", Age: 25}
		_, err := session.Insert(user)
		return err
	})

	assert.NoError(t, err)

	// 验证数据已插入
	var count int64
	count, err = database.Count(&User{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestWithTx_RollbackOnError(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	// 插入一条数据
	user := &User{Name: "Bob", Age: 30}
	_, err := database.Insert(user)
	require.NoError(t, err)

	// 事务中发生错误，应该回滚
	err = WithTx(ctx, func(ctx context.Context, session *xorm.Session) error {
		// 先插入一条数据
		newUser := &User{Name: "Charlie", Age: 35}
		_, err := session.Insert(newUser)
		if err != nil {
			return err
		}
		// 返回错误，触发回滚
		return errors.New("business logic error")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "business logic error")

	// 验证数据没有被插入（事务回滚）
	var count int64
	count, err = database.Count(&User{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count) // 只有第一条数据
}

func TestWithTx_CancelledContext(t *testing.T) {
	cleanupUsers(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消上下文

	err := WithTx(ctx, func(ctx context.Context, session *xorm.Session) error {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			user := &User{Name: "David", Age: 40}
			_, err := session.Insert(user)
			return err
		}
	})

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestWithTx_MultipleOperations(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	err := WithTx(ctx, func(ctx context.Context, session *xorm.Session) error {
		// 插入多条数据
		users := []User{
			{Name: "Frank", Age: 25},
			{Name: "Grace", Age: 30},
			{Name: "Henry", Age: 35},
		}

		for _, user := range users {
			_, err := session.Insert(&user)
			if err != nil {
				return err
			}
		}

		// 更新数据
		_, err := session.Where("name = ?", "Frank").Update(&User{Age: 26})
		return err
	})

	assert.NoError(t, err)

	// 验证数据
	var count int64
	count, err = database.Count(&User{})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 验证更新
	var user User
	has, err := database.Where("name = ?", "Frank").Get(&user)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.Equal(t, 26, user.Age)
}

func TestWithTxResult_Success(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	// 测试返回结果的事务
	result, err := WithTxResult(ctx, func(ctx context.Context, session *xorm.Session) (int64, error) {
		user := &User{Name: "Isabel", Age: 28}
		return session.Insert(user)
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result) // 插入了1条记录

	// 验证数据
	var count int64
	count, err = database.Count(&User{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestWithTxResult_Error(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	// 测试返回错误的事务
	result, err := WithTxResult(ctx, func(ctx context.Context, session *xorm.Session) (string, error) {
		user := &User{Name: "Jack", Age: 32}
		_, insertErr := session.Insert(user)
		if insertErr != nil {
			return "", insertErr
		}
		return "", errors.New("custom error")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "custom error")
	assert.Equal(t, "", result) // 错误时返回零值

	// 验证数据没有被插入（事务回滚）
	var count int64
	count, err = database.Count(&User{})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestWithTxResult_WithStruct(t *testing.T) {
	cleanupUsers(t)
	ctx := context.Background()

	// 测试返回结构体的事务
	result, err := WithTxResult(ctx, func(ctx context.Context, session *xorm.Session) (User, error) {
		user := &User{Name: "Kate", Age: 27}
		_, err := session.Insert(user)
		if err != nil {
			return User{}, err
		}

		// 查询并返回插入的用户
		var insertedUser User
		has, err := session.Where("name = ?", "Kate").Get(&insertedUser)
		if err != nil {
			return User{}, err
		}
		if !has {
			return User{}, errors.New("user not found")
		}
		return insertedUser, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "Kate", result.Name)
	assert.Equal(t, 27, result.Age)
	assert.NotZero(t, result.ID) // ID应该已经被赋值
}

// 基准测试
func BenchmarkWithTx(b *testing.B) {
	cleanupUsers(&testing.T{})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := WithTx(ctx, func(ctx context.Context, session *xorm.Session) error {
			user := &User{Name: fmt.Sprintf("User%d", i), Age: 25}
			_, err := session.Insert(user)
			return err
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWithTxResult(b *testing.B) {
	cleanupUsers(&testing.T{})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := WithTxResult(ctx, func(ctx context.Context, session *xorm.Session) (int64, error) {
			user := &User{Name: fmt.Sprintf("User%d", i), Age: 25}
			return session.Insert(user)
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
