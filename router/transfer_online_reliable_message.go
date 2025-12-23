package router

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

// ============== Mock 常量定义 ==============

const (
	RedisFailCode     = -1
	ReliableMsgSeqPre = "msg_seq_"
	RedisInterval     = "_"
	SeqExpireSeconds  = 10
	DefaultLocale     = "en-US"
)

var NoConnectionErr = errors.New("no conn available")

// ============== Protobuf Mock ==============

// Any 模拟 protobuf 的 Any 类型
type Any struct {
	TypeUrl string
	Value   []byte
}

// proto package mock
type protoPackage struct{}

func (protoPackage) Marshal(v interface{}) ([]byte, error) {
	// 使用JSON序列化替代protobuf
	return json.Marshal(v)
}

func (protoPackage) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

var proto = protoPackage{}

// ptypes package mock
type ptypesPackage struct{}

func (ptypesPackage) Is(any *Any, msg interface{}) bool {
	// 简化实现：检查类型URL
	// 在实际使用中，根据具体需要实现
	return false
}

func (ptypesPackage) UnmarshalAny(any *Any, msg interface{}) error {
	if any == nil || len(any.Value) == 0 {
		return errors.New("any is empty")
	}
	return json.Unmarshal(any.Value, msg)
}

func (ptypesPackage) MarshalAny(msg interface{}) (*Any, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &Any{
		TypeUrl: fmt.Sprintf("type.googleapis.com/%T", msg),
		Value:   data,
	}, nil
}

var ptypes = ptypesPackage{}

// ============== Mock 类型定义（模拟外部依赖的proto和接口） ==============

// Config Mock
type config struct {
	Service struct {
		DisableSendReliable bool
		IsNeedTTDB          bool
		IsStoreReliableMsg  bool
		DisablePingProcess  bool
	}
}

var mockCfg = &config{}

// config package mock
func Get() *config {
	return mockCfg
}

func (c *config) AppIndex(appName string) int {
	// Mock实现，返回固定值
	return 0
}

func (c *config) AppExist(appName string) bool {
	return true
}

// applog package mock
type applog struct{}

func (applog) Error(args ...interface{})                 { fmt.Println(args...) }
func (applog) Errorf(format string, args ...interface{}) { fmt.Printf(format+"\n", args...) }
func (applog) Warnf(format string, args ...interface{})  { fmt.Printf(format+"\n", args...) }
func (applog) Debugf(format string, args ...interface{}) { fmt.Printf(format+"\n", args...) }
func (applog) Infof(format string, args ...interface{})  { fmt.Printf(format+"\n", args...) }

var Applog = applog{}

// connector proto types mock
type ClientSourceEnum int32

const (
	CLIENT_SOURCE_ANDROID ClientSourceEnum = 1
	CLIENT_SOURCE_IOS     ClientSourceEnum = 2
)

type PushContent struct {
	Title      *I18N
	Value      *I18N
	Ticker     *I18N
	Message    string
	CreateTime int64
}

func (p *PushContent) GetTitle() *I18N    { return p.Title }
func (p *PushContent) GetValue() *I18N    { return p.Value }
func (p *PushContent) GetTicker() *I18N   { return p.Ticker }
func (p *PushContent) GetMessage() string { return p.Message }

type I18N struct {
	Value   string
	Locales map[string]string
	Params  []string
}

func (i *I18N) GetValue() string    { return i.Value }
func (i *I18N) GetParams() []string { return i.Params }

type ChatMsg struct {
	Ticker  string
	Message string
}

type UserAgent struct {
	Source       ClientSourceEnum
	AppVersion   string
	AppUIVersion string
}

type TransmitMessageRequest struct {
	UserId          string
	MsgId           string
	MsgType         int32
	MsgData         *Any
	Push            *PushContent
	MsgTypeName     string
	AppName         string
	DeviceIdentifer string
}

// proto_router proto types mock
type TransferMessageRequest struct {
	ReceiverId      string
	MsgId           string
	MsgType         int32
	MsgData         *Any
	Push            *PushContent
	DeviceIdPushes  []*DeviceIdPush
	MsgTypeName     string
	AppName         string
	DeviceIdentifer string
	Filters         map[string]string
	LimitVersion    *LimitVersion
	ForceLangs      []string
}

func (r *TransferMessageRequest) GetReceiverId() string              { return r.ReceiverId }
func (r *TransferMessageRequest) GetMsgId() string                   { return r.MsgId }
func (r *TransferMessageRequest) GetMsgType() int32                  { return r.MsgType }
func (r *TransferMessageRequest) GetMsgData() *Any                   { return r.MsgData }
func (r *TransferMessageRequest) GetPush() *PushContent              { return r.Push }
func (r *TransferMessageRequest) GetDeviceIdPushes() []*DeviceIdPush { return r.DeviceIdPushes }
func (r *TransferMessageRequest) GetFilters() map[string]string      { return r.Filters }

type DeviceIdPush struct {
	DeviceIds []string
	Push      *PushContent
}

func (d *DeviceIdPush) GetDeviceIds() []string { return d.DeviceIds }
func (d *DeviceIdPush) GetPush() *PushContent  { return d.Push }

type LimitVersion struct {
	MinAndroidVersion string
	MaxAndroidVersion string
	MinIosVersion     string
	MaxIosVersion     string
	MinUIVersion      string
}

type TransferPushMessageReply struct {
	IsUserOnline      bool
	DeviceIdentifiers []*DeviceIdentifier
}

type DeviceIdentifier struct {
	Identifer string
	IsOnline  bool
}

// tracing mock
type tracing struct{}

func (tracing) PropagateContextWithServiceContext(ctx context.Context) context.Context {
	return ctx
}

var Tracing = tracing{}

// router types mock
type ConnectorClientWrapper struct {
	DeviceID  string
	Locale    string
	Source    string
	UA        *UserAgent
	Connector ConnectorClient
}

type ConnectorClient interface {
	TransmitMessage(ctx context.Context, req *TransmitMessageRequest) error
}

// Router interface
type Router interface {
	PickConnectors(ctx context.Context, appName, userID, deviceIdentifier string, filters map[string]string) []*ConnectorClientWrapper
}

// Default Router implementation
type DefaultRouter struct{}

func (r *DefaultRouter) PickConnectors(ctx context.Context, appName, userID, deviceIdentifier string, filters map[string]string) []*ConnectorClientWrapper {
	// Mock实现，返回默认连接器
	return []*ConnectorClientWrapper{
		{
			DeviceID: "mock-device-001",
			Locale:   "en-US",
			Source:   "client",
			UA: &UserAgent{
				Source:     CLIENT_SOURCE_IOS,
				AppVersion: "1.0.0",
			},
			Connector: &DefaultConnector{},
		},
	}
}

type DefaultConnector struct{}

func (c *DefaultConnector) TransmitMessage(ctx context.Context, req *TransmitMessageRequest) error {
	return nil
}

// Redis client interface
type RouterRedisClient interface {
	GenSequenceID(ctx context.Context, key string, expireSeconds int) (int64, error)
	HCAD(ctx context.Context, appID, userId, deviceID, source, addr string) (int64, error)
	HCADSR(ctx context.Context, appID, userId, deviceID, source, addr string) (int64, error)
}

// Default Redis client implementation
type DefaultRouterRedisClient struct{}

func (r *DefaultRouterRedisClient) GenSequenceID(ctx context.Context, key string, expireSeconds int) (int64, error) {
	return 1, nil
}

func (r *DefaultRouterRedisClient) HCAD(ctx context.Context, appID, userId, deviceID, source, addr string) (int64, error) {
	return 1, nil
}

func (r *DefaultRouterRedisClient) HCADSR(ctx context.Context, appID, userId, deviceID, source, addr string) (int64, error) {
	return 1, nil
}

// MsgDB mock
type ReliableMsg interface {
	InsertMsg(ctx context.Context, appID, userID int, seq int64, deviceIdentifier, msgID, msgData string) error
}

type DefaultReliableMsg struct{}

func (d *DefaultReliableMsg) InsertMsg(ctx context.Context, appID, userID int, seq int64, deviceIdentifier, msgID, msgData string) error {
	return nil
}

// util mock
type util struct{}

func (util) ContainsString(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

var Util = util{}

// common_util mock
type common_util struct{}

func (common_util) VersionGreaterThanOrEqualTo(v1, v2 string) bool {
	// 简单的版本比较mock
	return v1 >= v2
}

var CommonUtil = common_util{}

// connector error check mock
func IsErrUserNotExist(err error) bool {
	return err != nil && err.Error() == "user not exist"
}

// ============== RouterServer 完全照抄原始实现 ==============

type RouterServer struct {
	Store  RouterRedisClient
	router Router
	MsgDB  ReliableMsg
}

func NewRouterServer(redisClient RouterRedisClient, msgDB ReliableMsg, router Router) *RouterServer {
	return &RouterServer{
		Store:  redisClient,
		router: router,
		MsgDB:  msgDB,
	}
}

// ============== 以下代码完全照抄原始实现 ==============

// TransferOnlineReliableMessage 完全照抄原始实现
func (s *RouterServer) TransferOnlineReliableMessage(ctx context.Context, in *TransferMessageRequest) (*TransferPushMessageReply, error) {
	cfg := Get()
	if cfg.Service.DisableSendReliable {
		return nil, nil
	}
	var err error
	var appIDInt, userIdInt int
	var rpl = &TransferPushMessageReply{}
	if len(in.ReceiverId) == 0 {
		err := errors.New("receiver id is empty when transfer message")
		Applog.Error(err)
		return nil, err
	}
	userIdInt, err = strconv.Atoi(in.ReceiverId)
	if err != nil {
		err := errors.New("receiver id is not valid")
		Applog.Error(err)
		return nil, err
	}
	if len(in.GetMsgId()) == 0 {
		err := fmt.Errorf("msgId is empty when transfer msg, uid: %+v", in.ReceiverId)
		Applog.Error(err)
		return nil, err
	}

	now := time.Now().UnixNano() / 1000000
	if in.GetPush() != nil {
		in.Push.CreateTime = now
	}
	if len(in.GetDeviceIdPushes()) > 0 {
		for _, deviceIdPush := range in.GetDeviceIdPushes() {
			if deviceIdPush.GetPush() != nil {
				deviceIdPush.Push.CreateTime = now
			}
		}
	}

	tm := time.Now()

	appIDInt = cfg.AppIndex(in.AppName)
	seq, err := s.genTTDBSeq(ctx, in.AppName, in.ReceiverId, tm)
	if err != nil {
		Applog.Error(err)
		return nil, err
	}

	connectorWrappers := s.router.PickConnectors(ctx, in.AppName, in.ReceiverId, in.DeviceIdentifer, in.GetFilters())
	if len(connectorWrappers) == 0 {
		return rpl, nil
	}
	for _, w := range connectorWrappers {
		if w.DeviceID != "" {
			rpl.DeviceIdentifiers = append(rpl.DeviceIdentifiers, &DeviceIdentifier{
				Identifer: w.DeviceID,
				IsOnline:  true,
			})
		}
	}
	rpl.IsUserOnline = true
	ctx = Tracing.PropagateContextWithServiceContext(ctx)
	//TODO 一期不做消息存储
	if cfg.Service.IsStoreReliableMsg {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					buf := make([]byte, 1<<15)
					n := runtime.Stack(buf, false)
					err := fmt.Errorf("%v, STACK: %s", r, buf[0:n])
					Applog.Errorf("TransferReliableMessage  panic :%+v", err)
				}
			}()
			raw, err := proto.Marshal(in)
			if err != nil {
				Applog.Errorf("proto.Marshal err:%+v msg is :%+v appID is :%d userID is :%d, seq is :%d", err, *in, appIDInt, userIdInt, seq)
				return
			}
			err = s.MsgDB.InsertMsg(ctx, appIDInt, userIdInt, seq, in.DeviceIdentifer, in.MsgId, base64.StdEncoding.EncodeToString(raw))
			if err != nil {
				Applog.Errorf("insert msgdb err:%+v msg is :%+v appID is :%d userID is :%d, seq is :%d", err, *in, appIDInt, userIdInt, seq)
			}
		}()
	}
	go func() {
		for _, wrapper := range connectorWrappers {
			if in.LimitVersion != nil && s.isLimitVersion(wrapper, in.LimitVersion) {
				Applog.Debugf("isLimitVersion msg is :%+v", *in)
				continue
			}
			if s.isNotForcedLangs(wrapper.Locale, in.ForceLangs) {
				continue
			}
			push := PushContent{}
			originPush := s.getOriginPush(in, wrapper.DeviceID)
			if originPush != nil {
				push = processPush(*originPush, wrapper.Locale)
			}

			if ptypes.Is(in.MsgData, &ChatMsg{}) {
				data, err := processChatMsg(in.GetMsgData(), push)
				if err != nil {
					Applog.Error(err)
					return
				}
				in.MsgData = data
			}
			err := wrapper.Connector.TransmitMessage(ctx, &TransmitMessageRequest{
				UserId:          in.ReceiverId,
				MsgId:           in.GetMsgId(),
				MsgType:         in.GetMsgType(),
				MsgData:         in.GetMsgData(),
				Push:            &push,
				MsgTypeName:     in.MsgTypeName,
				AppName:         in.AppName,
				DeviceIdentifer: wrapper.DeviceID,
			})
			if err != nil {
				s.handleError(ctx, err, in.AppName, in.ReceiverId, wrapper.DeviceID, wrapper.Source)
			}
		}
	}()
	return rpl, nil
}

func (s *RouterServer) getOriginPush(req *TransferMessageRequest, deviceId string) *PushContent {
	if len(req.GetDeviceIdPushes()) > 0 {
		for _, deviceIdPush := range req.GetDeviceIdPushes() {
			for _, tmpDeviceId := range deviceIdPush.GetDeviceIds() {
				if tmpDeviceId == deviceId {
					return deviceIdPush.GetPush()
				}
			}
		}
	}

	return req.GetPush()
}

/*
	生成消息序列号：
		key: appID+userID+当前秒（eg: msg_seq_0_602_20200114144545）
		seq = int64(timestamp+incr(key))
*/

func (s *RouterServer) genTTDBSeq(ctx context.Context, appID, userID string, tm time.Time) (int64, error) {
	secTmStr := tm.Format("20060102150405")
	seqKey := ReliableMsgSeqPre + appID + RedisInterval + userID + secTmStr
	seq, err := s.Store.GenSequenceID(ctx, seqKey, SeqExpireSeconds)
	if err != nil {
		return 0, err
	}
	val := fmt.Sprintf("%d", tm.Unix()) + fmt.Sprintf("%04d", seq)
	seq, _ = strconv.ParseInt(val, 10, 64)
	return seq, nil
}

func processChatMsg(msgData *Any, push PushContent) (*Any, error) {
	var chatMsg ChatMsg
	err := ptypes.UnmarshalAny(msgData, &chatMsg)
	if err != nil {
		Applog.Error(err)
		return nil, err
	}

	chatMsg.Ticker = push.GetTicker().GetValue()
	chatMsg.Message = push.GetMessage()

	return ptypes.MarshalAny(&chatMsg)
}

func processPush(push PushContent, locale string) PushContent {
	if push.GetTitle() != nil {
		push.Title = parseI18n(*push.Title, locale)
	}
	if push.GetValue() != nil {
		push.Value = parseI18n(*push.Value, locale)
	}
	if push.GetTicker() != nil {
		push.Ticker = parseI18n(*push.Ticker, locale)
	}
	return push
}

func parseI18n(i18n I18N, locale string) *I18N {
	var localeStr string
	if s, ok := i18n.Locales[locale]; ok {
		localeStr = s
	} else if len(i18n.Value) > 0 {
		localeStr = i18n.Value
	} else if s, ok := i18n.Locales[DefaultLocale]; ok {
		localeStr = s
	} else if len(i18n.Locales) > 0 {
		err := fmt.Errorf("no locale found in i18n field %v, locale %v", i18n, locale)
		Applog.Error(err)
		return nil
	}
	if len(i18n.GetParams()) > 0 {
		s := make([]interface{}, len(i18n.GetParams()))
		for i, v := range i18n.GetParams() {
			s[i] = v
		}
		i18n.Value = fmt.Sprintf(localeStr, s...)
	} else {
		i18n.Value = localeStr
	}
	i18n.Locales = nil
	i18n.Params = nil
	return &i18n
}

func (s *RouterServer) isLimitVersion(wrapper *ConnectorClientWrapper, limit *LimitVersion) bool {
	var min, max string
	ua := wrapper.UA
	switch ua.Source {
	case CLIENT_SOURCE_ANDROID:
		min = limit.MinAndroidVersion
		max = limit.MaxAndroidVersion
	case CLIENT_SOURCE_IOS:
		min = limit.MinIosVersion
		max = limit.MaxIosVersion
	default:
		return true
	}
	if s.versionNotInRange(min, max, ua.AppVersion) {
		return true
	}
	if len(ua.AppUIVersion) != 0 {
		return s.versionNotInRange(limit.MinUIVersion, "", ua.AppUIVersion)
	}
	return false
}

func (s *RouterServer) versionNotInRange(min string, max string, v string) bool {
	if min != "" && max == "" {
		if CommonUtil.VersionGreaterThanOrEqualTo(v, min) {
			return false
		} else {
			return true
		}
	}
	if min == "" && max == "" {
		return false
	}
	if min != "" && max != "" {
		if CommonUtil.VersionGreaterThanOrEqualTo(max, v) && CommonUtil.VersionGreaterThanOrEqualTo(v, min) {
			return false
		} else {
			return true
		}
	}
	if min == "" && max != "" {
		if CommonUtil.VersionGreaterThanOrEqualTo(max, v) {
			return false
		} else {
			return true
		}
	}
	return false
}

func (p *RouterServer) isNotForcedLangs(lang string, forcedLangs []string) bool {
	// force device language
	if len(forcedLangs) > 0 {
		if !Util.ContainsString(lang, forcedLangs) {
			Applog.Warnf("language %v not in force languages %v", lang, forcedLangs)
			return true
		}
	}
	return false
}

func (s *RouterServer) handleError(ctx context.Context, err error, appID, userId, deviceID, source string) {
	if IsErrUserNotExist(err) { //如果connector返回该错误，需要主动删除相应的路由信息
		if userId == "0" { // ANONYMOUS_USER_ID_STRING mock
			s.Store.HCAD(ctx, appID, userId, deviceID, source, "")
		} else {
			s.Store.HCADSR(ctx, appID, userId, deviceID, source, "")
		}
		Applog.Infof("delete router info, uid: %v, deviceID %v, source: %v by connector", userId, deviceID, source)
	} else {
		Applog.Error(err)
	}
}
