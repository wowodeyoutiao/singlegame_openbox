package db

import (
	"time"

	"nanoserver/db/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

const asyncTaskBacklog = 128

var (
	database *xorm.Engine
	logger   *log.Entry
	chWrite  chan interface{} // async write channel
	chUpdate chan interface{} // async update channel
)

type options struct {
	showSQL      bool
	maxOpenConns int
	maxIdleConns int
}

// SetDatabase sets the global database instance
func SetDatabase(db *xorm.Engine) {
	database = db
}

// GetDatabase returns the global database instance
func GetDatabase() *xorm.Engine {
	return database
}

// ModelOption specifies an option for dialing a xordefaultModel.
type ModelOption func(*options)

// MaxIdleConns specifies the max idle connect numbers.
func MaxIdleConns(i int) ModelOption {
	return func(opts *options) {
		opts.maxIdleConns = i
	}
}

// MaxOpenConns specifies the max open connect numbers.
func MaxOpenConns(i int) ModelOption {
	return func(opts *options) {
		opts.maxOpenConns = i
	}
}

// ShowSQL specifies the buffer size.
func ShowSQL(show bool) ModelOption {
	return func(opts *options) {
		opts.showSQL = show
	}
}

func envInit() {
	// async task
	go func() {
		for {
			select {
			case t, ok := <-chWrite:
				if !ok {
					return
				}

				if _, err := database.Insert(t); err != nil {
					logger.Error(err)
				}

			case t, ok := <-chUpdate:
				if !ok {
					return
				}

				if _, err := database.Update(t); err != nil {
					logger.Error(err)
				}
			}
		}
	}()

	// 定时ping数据库, 保持连接池连接
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		for {
			select {
			case <-ticker.C:
				database.Ping()
			}
		}
	}()
}

// New create the database's connection
func MustStartup(dsn string, opts ...ModelOption) func() {
	logger = log.WithField("component", "model")
	settings := &options{
		maxIdleConns: defaultMaxConns,
		maxOpenConns: defaultMaxConns,
		showSQL:      true,
	}

	// options handle
	for _, opt := range opts {
		opt(settings)
	}

	logger.Infof("DSN=%s ShowSQL=%t MaxIdleConn=%v MaxOpenConn=%v", dsn, settings.showSQL, settings.maxIdleConns, settings.maxOpenConns)

	// create database instance
	if db, err := xorm.NewEngine("mysql", dsn); err != nil {
		panic(err)
	} else {
		database = db
	}

	// 设置日志相关
	database.SetLogger(&Logger{Entry: logger.WithField("orm", "xorm")})

	chWrite = make(chan interface{}, asyncTaskBacklog)
	chUpdate = make(chan interface{}, asyncTaskBacklog)

	// options
	database.SetMaxIdleConns(settings.maxIdleConns)
	database.SetMaxOpenConns(settings.maxOpenConns)
	database.ShowSQL(settings.showSQL)

	syncSchema()
	envInit()

	closer := func() {
		// time.AfterFunc(3*time.Second, func() {
		// 	close(chWrite)
		// 	close(chUpdate)
		// 	database.Close()
		// 	logger.Info("stopped")
		// })

		logger.Info("stopped")
	}

	return closer
}

func MustClose() {
	close(chWrite)
	close(chUpdate)
	database.Close()
	logger.Info("db stopped")
}

func syncSchema() {
	if err := database.StoreEngine("InnoDB").Sync2(
		new(model.Login),
		new(model.LoginRole),
		new(model.Online),
		// new(model.Order),
		// new(model.Recharge),
		new(model.Register),
		new(model.ThirdAccount),
		// new(model.Trade),
		new(model.User),
		new(model.Mail),
		new(model.Attachment),
		new(model.Sign),
		new(model.Store),
		new(model.MailTask),
		new(model.AttachmentTask),
		new(model.AttachmentActivation),
		new(model.Activation),
		new(model.Activity),
		new(model.Privilege),
		new(model.RankList),
		new(model.Grade),
		new(model.Turntable),
		new(model.TreasureBox),
		new(model.UserExtractAccount),
		// new(model.ItemControlRefresh),
		// new(model.ItemControlLast),
		new(model.Partner),
		new(model.PartnerSlot),
		new(model.UserBox),
		new(model.Summon),
		new(model.UserPet),
		new(model.InstanceData),
		new(model.Guild),
		new(model.GuildMember),
		new(model.GuildApplication),
		new(model.GuildLog),
		new(model.GuildDonate),
		new(model.GuildStoreBuy),
		new(model.PvpRecord),
		new(model.UserAttribute),
		new(model.War),
		new(model.Lottery),
		new(model.Avatar),
		new(model.Refine),
		new(model.Artifact),
	); err != nil {
		log.Fatal("failed to sync database schema: ", err)
	}
}
