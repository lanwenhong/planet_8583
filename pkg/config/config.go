package config

import (
	"time"
	
	"github.com/lanwenhong/planet_8583/pkg/utils"
	
	"github.com/lanwenhong/lgobase/confparse"
	"github.com/lanwenhong/lgobase/gconfig"
)

type Config struct {
	PlanetAddr         string        `confpos:"planet:addr" dtype:"base"`
	PlanetPort         int           `confpos:"planet:port" dtype:"base"`
	PlanetConnTimetout time.Duration `confpos:"planet:conn_timeout" dtype:"base"`
	PlanetReadTimeout  time.Duration `confpos:"planet:read_timeout" dtype:"base"`
	PlanetWriteTimeout time.Duration `confpos:"planet:write_timeout" dtype:"base"`
	PlanetCertPath     string        `confpos:"planet:cert_path" dtype:"base"`
	Env                string        `confpos:"planet:env" dtype:"base"`
	PaymentBusicds     []string      `confpos:"planet:payment_busicds" dtype:"base" item_split:","`
	RefundBusicds      []string      `confpos:"planet:refund_busicds" dtype:"base" item_split:","`
	
	// log
	LogFile    string `confpos:"log:logfile" dtype:"base"`
	LogFileErr string `confpos:"log:logfile_err" dtype:"base"`
	LogDir     string `confpos:"log:logdir" dtype:"base"`
	LogLevel   string `confpos:"log:loglevel" dtype:"base"`
	LogStdOut  bool   `confpos:"log:logstdout" dtype:"base"`
	Colorfull  bool   `confpos:"log:colorfull" dtype:"base"`
	RedisLog   bool   `confpos:"log:redis_log" dtype:"base"`
	
	// DB
	TokenFile string `confpos:"db:token_file" dtype:"base"`
	TradeDB   string `confpos:"db:trade" dtype:"base"`
	
	// Redis
	RedisAddr    string   `confpos:"redis:redis_addr" dtype:"base"`
	RedisAddrs   []string `confpos:"redis:redis_addrs" dtype:"base" item_split:","`
	PoolSize     int      `confpos:"redis:pool_size" dtype:"base"`
	MinIdle      int      `confpos:"redis:min_idle" dtype:"base"`
	ReadTimeout  int      `confpos:"redis:read_timeout" dtype:"base"`
	WriteTimeout int      `confpos:"redis:write_timeout" dtype:"base"`
	ConnTimeout  int      `confpos:"redis:connect_timeout" dtype:"base"`
	Db           int      `confpos:"redis:db" dtype:"base"`
	RedisUser    string   `confpos:"redis:user" dtype:"base"`
	RedisPasswd  string   `confpos:"redis:passwd" dtype:"base"`
	Cluster      bool     `confpos:"redis:cluster" dtype:"base"`
	
	// script
	CaptureChannelID string   `confpos:"script:capture_channel_id" dtype:"base"`
	BatchCaptureDay  int      `confpos:"script:batch_capture_day" dtype:"base"`
	CaptureYesterday bool     `confpos:"script:capture_yesterday" dtype:"base"`
	MailTo           []string `confpos:"script:mail_to" dtype:"base" item_split:","`
	MailUserName     string   `confpos:"script:mail_username" dtype:"base"`
	MailPassword     string   `confpos:"script:mail_password" dtype:"base"`
}

var Conf = new(Config)

func ParseConfig(filename string) {
	cfg := gconfig.NewGconf(filename)
	utils.MustNil(cfg.GconfParse())
	utils.MustNil(confparse.CpaseNew(filename).CparseGo(Conf, cfg))
}
