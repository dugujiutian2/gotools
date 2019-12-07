package login

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/hero1s/gotools/log"
	"github.com/hero1s/gotools/login/fetch"
	"io"
	"math/rand"
	"sort"
	"time"
)

type WeChatInfo struct {
	OpenID     string   `json:"openid"`     // 用户的唯一标识
	NickName   string   `json:"nickname"`   // 用户昵称
	Sex        int      `json:"sex"`        // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province   string   `json:"province"`   // 用户个人资料填写的省份
	City       string   `json:"city"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country"`    // 国家，如中国为CN
	HeadImgURL string   `json:"headimgurl"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Privilege  []string `json:"privilege"`  // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	UnionID    string   `json:"unionid"`    // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。详见：获取用户个人信息（UnionID机制）
	Language   string   `json:"language"`   // 语言
}

func GetWeChatUserInfo(accessToken, openID string) (*WeChatInfo, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s",
		accessToken,
		openID,
	)
	body, err := fetch.Cmd(fetch.Request{
		Method: "GET",
		URL:    url,
	})
	if err != nil {
		return nil, err
	}
	var result WeChatInfo
	err = json.Unmarshal(body, &result)
	return &result, err
}

//一些有关微信方面的操作
var (
	//h5方面的微信,游戏客户端的不需要下面这两个参数
	WeChatAppID     = "wxbc1ba2a3d45458c3"
	WeChatAppSecret = "15ba90574a2d49069a7261f4f7e5d5ad"
)

type WeChatToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

//通过code来获取aceess_token及open_id
func GetWeChatOpenIdAccessToken(code string) (*WeChatToken, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?"+
		"appid=%v&secret=%v&code=%v&grant_type=authorization_code", WeChatAppID, WeChatAppSecret, code)
	body, err := fetch.Cmd(fetch.Request{
		Method: "GET",
		URL:    url,
	})
	if err != nil {
		return nil, err
	}
	var result WeChatToken
	err = json.Unmarshal(body, &result)
	return &result, err
}

type AccessToken struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type WeChatTicket struct {
	ErrCode   int64  `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

type WeChatSign struct {
	AppId     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

func getAccessToken() (string, error) {
	var result AccessToken
	var accessToken string
	var err error
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		WeChatAppID, WeChatAppSecret)
	var body []byte
	body, err = fetch.Cmd(fetch.Request{
		Method: "GET",
		URL:    url,
	})
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("获取微信access_token失败:%v\n,状态码:%v", result.ErrMsg, result.ErrCode)
		return "", err
	}
	accessToken = result.AccessToken
	return accessToken, err
}
func GetWeChatTicket(uri string) (*WeChatSign, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	return getWeChatTicket(accessToken, uri)
}

func getWeChatTicket(accessToken, uri string) (*WeChatSign, error) {
	var result WeChatTicket
	var ticket string
	var err error

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%v&type=jsapi", accessToken)
	var body []byte
	body, err = fetch.Cmd(fetch.Request{
		Method: "GET",
		URL:    url,
	})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("获取微信ticket:%v\n,状态码:%v", result.ErrMsg, result.ErrCode)
		return nil, err
	}
	ticket = result.Ticket
	log.Debug(fmt.Sprintf("从微信服务器里获取ticket成功:%v,err:%v,expirein:%v", ticket, err, result.ExpiresIn))

	nonceStr := RandomStr(16)
	timestamp := time.Now().Unix()
	ticketStr := ticket
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, uri)
	signStr := Signature(str)

	sign := WeChatSign{
		AppId:     WeChatAppID,
		NonceStr:  nonceStr,
		Timestamp: timestamp,
		Signature: signStr,
	}
	return &sign, err
}

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func Signature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}