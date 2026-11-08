package shushumain

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"nanoserver/db"
	"nanoserver/pkg/shushu/src/thinkingdata"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Closer func()

var (
	te              thinkingdata.TDAnalytics
	DefaultReporter Reporter
)
var logger = log.WithField("component", "shushu")

// GetTE 返回 te 的指针
func GetTE() *thinkingdata.TDAnalytics {
	return &te
}

func Startup() {
	thinkingdata.SetLogLevel(thinkingdata.TDLogLevelError)

	// e.g. init consumer, you can choose different consumer
	//consumer, err := generateDebugConsumer() // TDDebugConsumer
	// consumer, err := generateBatchConsumer() // TDBatchConsumer
	// consumer, err := generateLogConsumer() // TDLogConsumer
	config := thinkingdata.TDBatchConfig{
		ServerUrl: viper.GetString("shushu.shushuUrl"),
		AppId:     viper.GetString("shushu.shushuAppId"),
		AutoFlush: viper.GetBool("shushu.shushuAutoFlush"),
		BatchSize: viper.GetInt("shushu.shushuBatchSize"),
		Interval:  viper.GetInt("shushu.shushuInterval"),
	}
	consumer, err := thinkingdata.NewBatchConsumerWithConfig(config)
	if err != nil {
		// consumer init error
	}
	te = thinkingdata.New(consumer)
	DefaultReporter = *NewReporter(4, 1024)
	DefaultReporter.Start()
}
func Stop() {
	DefaultReporter.Stop()
}

// Reporter 用于管理数据上报的通道和相关逻辑
type Reporter struct {
	reportChan      chan map[string]interface{} // 数据上报通道
	userSetChan     chan map[string]interface{} // UserSet数据上报通道
	tencentChan     chan map[string]interface{} // tencent数据上报通道
	userSetOnceChan chan map[string]interface{} // UserSetOnce数据上报通道
	mailChan        chan protocol.MailAdd       // 邮件的通道
	wg              sync.WaitGroup              // 等待组，用于等待所有处理协程完成
	numWorkers      int                         // 处理协程的数量
}

// NewReporter 创建一个新的 Reporter 实例
func NewReporter(numWorkers int, bufferSize int) *Reporter {
	return &Reporter{
		reportChan:      make(chan map[string]interface{}, bufferSize),
		userSetChan:     make(chan map[string]interface{}, bufferSize),
		userSetOnceChan: make(chan map[string]interface{}, bufferSize),
		tencentChan:     make(chan map[string]interface{}, bufferSize),
		mailChan:        make(chan protocol.MailAdd, bufferSize),
		numWorkers:      numWorkers,
	}
}

// Start 启动 Reporter，创建指定数量的处理协程
func (r *Reporter) Start() {
	for i := 0; i < r.numWorkers; i++ {
		r.wg.Add(1)
		go r.processReports()
		r.wg.Add(1)
		go r.processUserSetChannel()
		r.wg.Add(1)
		go r.processUserSetOnceChannel()
		r.wg.Add(1)
		go r.processTencentChannel()
		r.wg.Add(1)
		go r.processMailChannel()
	}
}

// Stop 停止 Reporter，等待所有处理协程完成并关闭通道
func (r *Reporter) Stop() {
	close(r.reportChan)      // 关闭旧通道
	close(r.userSetChan)     // 关闭UserSet通道
	close(r.userSetOnceChan) // 关闭UserSetOnce通道
	close(r.tencentChan)
	close(r.mailChan)
	r.wg.Wait() // 等待所有处理协程完成
	logger.Info("R关闭了")
	te.Close()
}

func (r *Reporter) reportData(data map[string]interface{}) {
	select {
	case r.reportChan <- data:
	default:
		logger.Info("Report channel is full, dropping data")
	}
}

// reportUserSetData 向UserSet通道发送数据
func (r *Reporter) reportUserSetData(data map[string]interface{}) {
	select {
	case r.userSetChan <- data:
	default:
		logger.Info("UserSet channel is full, dropping data")
	}
}

// reportUserSetOnceData 向UserSetOnce通道发送数据
func (r *Reporter) reportUserSetOnceData(data map[string]interface{}) {
	select {
	case r.userSetOnceChan <- data:
	default:
		logger.Info("UserSetOnce channel is full, dropping data")
	}
}

func (r *Reporter) reportTencentData(data map[string]interface{}) {
	select {
	case r.tencentChan <- data:
	default:
		logger.Info("Tencent channel is full, dropping data")
	}
}
func (r *Reporter) reportMailData(data protocol.MailAdd) {
	r.mailChan <- data
}

// ReportData 调用默认 Reporter 发送数据
func ReportData(data map[string]interface{}) {
	DefaultReporter.reportData(data)
}

// ReportUserSetData 调用默认 Reporter 向UserSet通道发送数据
func ReportUserSetData(data map[string]interface{}) {
	DefaultReporter.reportUserSetData(data)
}

// ReportUserSetOnceData 调用默认 Reporter 向UserSetOnce通道发送数据
func ReportUserSetOnceData(data map[string]interface{}) {
	DefaultReporter.reportUserSetOnceData(data)
}

func ReportTencentData(data map[string]interface{}) {
	DefaultReporter.reportTencentData(data)
}

func ReportMailData(data protocol.MailAdd) {
	DefaultReporter.reportMailData(data)
}

// processUserSetChannel 处理UserSet通道数据的协程函数
func (r *Reporter) processUserSetChannel() {
	defer r.wg.Done()
	for data := range r.userSetChan {
		UploadUserSetDate(GetTE(), data)
	}
}

// processUserSetOnceChannel 处理UserSetOnce通道数据的协程函数
func (r *Reporter) processUserSetOnceChannel() {
	defer r.wg.Done()
	for data := range r.userSetOnceChan {
		UploadUserSetOnceDate(GetTE(), data)
	}
}
func (r *Reporter) processTencentChannel() {
	defer r.wg.Done()
	for data := range r.tencentChan {
		UploadTencentDate(data)
	}
}

func (r *Reporter) processMailChannel() {
	defer r.wg.Done()
	for data := range r.mailChan {
		AddMailDate(data)
	}
}

func UploadUserSetDate(te *thinkingdata.TDAnalytics, properties map[string]interface{}) {
	accountId := ""
	distinctId := ""
	// 优先使用 tmp_role_id
	if roleID, ok := properties["tmp_role_id"].(int64); ok && roleID != 0 {
		distinctId = strconv.FormatInt(roleID, 10)
	}

	if account, ok := properties["tmp_account"].(string); ok && account != "" {
		accountId = account
	}

	for j := 0; j < 1; j++ {
		if accountId != "" && distinctId != "" {
			err := te.UserSet(distinctId, accountId, properties)
			if err != nil {
				logger.Errorf("Track error: %v\n,%#+v", err, properties)
			}
		}
	}
}

func UploadUserSetOnceDate(te *thinkingdata.TDAnalytics, properties map[string]interface{}) {
	accountId := ""
	distinctId := ""
	// 优先使用 tmp_role_id
	if roleID, ok := properties["tmp_role_id"].(string); ok && roleID != "" {
		distinctId = roleID
	}
	if account, ok := properties["tmp_account"].(string); ok && account != "" {
		accountId = account
	}

	for j := 0; j < 1; j++ {
		if accountId != "" && distinctId != "" {
			err := te.UserSetOnce(distinctId, accountId, properties)
			if err != nil {
				logger.Errorf("Track error: %v\n,%#+v", err, properties)
			}
		}
	}
}

// processReports 处理数据上报的协程函数
func (r *Reporter) processReports() {
	defer r.wg.Done()
	for data := range r.reportChan {
		// name, ok := data["_name"].(string)
		// if !ok {
		// 	fmt.Printf("Invalid data format: missing or invalid _name field\n")
		// 	continue
		// }
		// delete(data, "_name")
		UploadDate(GetTE(), data)
	}
}

func UploadDate(te *thinkingdata.TDAnalytics, properties map[string]interface{}) {
	accountId := ""
	distinctId := ""
	name := ""
	// 优先使用 tmp_role_id
	if roleID, ok := properties["tmp_role_id"].(int64); ok && roleID != 0 {
		distinctId = strconv.FormatInt(roleID, 10)
	}

	if account, ok := properties["tmp_account"].(string); ok && account != "" {
		accountId = account
	}
	if nameInfo, ok := properties["name"].(string); ok && nameInfo != "" {
		name = nameInfo
	}

	for j := 0; j < 1; j++ {
		if accountId != "" && distinctId != "" {
			err := te.Track(distinctId, accountId, name, properties)
			if err != nil {
				logger.Errorf("Track error: %v\n,%#+v", err, properties)
			}
		}
	}
}
func mainnn() {
	url := "https://api.e.qq.com/v1.3/user_actions/add?access_token=xxxxxxxx&timestamp=1679379053&nonce=xxxxxxxxxx"
	data := `{
  "account_id": 1111111111, //示例账户id
  "user_action_set_id": 1111111111, // 在DATANEXUS平台创建数据源ID
  "actions": [
      {
      "outer_action_id": "示例唯一行为 id_请指定",
      "action_time": 1492998090,
      "user_id": {
        "hash_imei": "f9efca36a..."
      },
      "action_type": "PURCHASE",
      "action_param": {
        "value": 28,
        "quantity": 5,
        "brand_name": "my_brand",
        "int_example": 123
      },
      "product_inform": {
        "content_type": "EC",
        "category_path": [
          "家用电器/厨房小电/豆浆机",
          "本地生活旅游出行/旅游出行/机票火车票"
        ]
      },
      "channel": "TENCENT"
    }
  ]}`
	payload := bytes.NewReader([]byte(data))
	req, _ := http.NewRequest(http.MethodPost, url, payload)
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle error
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	logger.Info(string(body))
}
func UploadTencentDate(properties map[string]interface{}) {
	claim_type := properties["claim_type"].(string)
	nowTime := time.Now().Unix()
	url := viper.GetString("tencentad.url")
	access_token := viper.GetString("tencentad.access_token")
	userActionSetIDStr := viper.GetString("tencentad.set_id")
	accountIdStr := viper.GetString("tencentad.account_id")
	account_id, err := strconv.ParseInt(accountIdStr, 10, 64)
	if err != nil {
		logger.Info("Invalid account_id:", err)
		return
	}
	userActionSetID, err := strconv.ParseInt(userActionSetIDStr, 10, 64)
	if err != nil {
		logger.Info("Invalid userActionSetID:", err)
		return
	}
	//雪花id
	nonce := snowflake.GenID()
	nonceStr := strconv.FormatInt(nonce, 10)

	urlDetail := url + "?access_token=" + access_token + "&timestamp=" + strconv.FormatInt(nowTime, 10) + "&nonce=" + nonceStr
	reqBody := TencentGuiderRequest{
		AccountID:       account_id,
		UserActionSetID: userActionSetID,
		Actions: []TencentGuiderAction{
			{
				OuterActionID: "1111",
				ActionTime:    nowTime,
				UserID: map[string]string{
					"wechat_openid": "oIJgI7L6gIozhh2WCi7VBzPETWPg",
					"wechat_app_id": "wx9d99386ab1e0721d",
					// "hash_imei": "f9efca36a...",
				},
				ActionType:  claim_type,
				ActionParam: map[string]interface{}{
					// "value":       28,
					// "quantity":    5,
					// "brand_name":  "my_brand",
					// "int_example": 123,
				},
				// ProductInform: struct {
				// 	ContentType  string   `json:"content_type"`
				// 	CategoryPath []string `json:"category_path"`
				// }{
				// 	ContentType: "EC",
				// 	CategoryPath: []string{
				// 		"家用电器/厨房小电/豆浆机",
				// 		"本地生活旅游出行/旅游出行/机票火车票",
				// 	},
				// },
				Channel: "TENCENT",
			},
		},
	}

	reqJson, err := json.Marshal(&reqBody)
	logger.Info("json.Marshal error:", reqJson)
	if err != nil {
		logger.Info("json.Marshal error:", err)
		return
	}

	urlReq, err := http.NewRequest(http.MethodPost, urlDetail, bytes.NewReader(reqJson))
	if err != nil {
		logger.Info("http.NewRequest error:", err)
		return
	}
	res, err := HttpDo(urlReq)
	logger.Info("HttpDo result:", res)
	if err != nil {
		logger.Error(err)
		return
	}
}
func AddMailDate(mailadd protocol.MailAdd) {
	uid := mailadd.Uid
	recipientUid := mailadd.RecipientUid
	desc := mailadd.Desc
	title := mailadd.Title
	channelCode := mailadd.ChannelCode
	attachments := make([]protocol.Attachment, len(mailadd.Attachments))
	for i, v := range mailadd.Attachments {
		attachments[i] = protocol.Attachment{
			UserPropsId: v.UserPropsId,
			Num:         v.Num,
		}
	}
	_, _, err := db.MailAdd(
		uid,
		recipientUid,
		title,
		desc,
		attachments,
		channelCode,
	)
	if err != nil {
		logger.Info("MailAdd error:", err)
	}
}
func HttpDo(req *http.Request) ([]byte, error) {
	// client := &http.Client{}
	req.Header.Add("content-type", "application/json")
	// resp, err := client.Do(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type TencentGuiderRequest2 struct {
	AccountID int64                 `json:"account_id"`
	Type      string                `json:"type"`
	Actions   []TencentGuiderAction `json:"actions"`
}
type TencentGuiderRequest struct {
	AccountID       int64                 `json:"account_id"`
	UserActionSetID int64                 `json:"user_action_set_id"`
	Actions         []TencentGuiderAction `json:"actions"`
}
type TencentGuiderAction struct {
	OuterActionID string                 `json:"outer_action_id"`
	ActionTime    int64                  `json:"action_time"`
	UserID        map[string]string      `json:"user_id"`
	ActionType    string                 `json:"action_type"`
	ActionParam   map[string]interface{} `json:"action_param"`
	ProductInform struct {
		ContentType  string   `json:"content_type"`
		CategoryPath []string `json:"category_path"`
	} `json:"product_inform"`
	Channel string `json:"channel"`
}
