package sms

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ihuyiSender struct {
	apiID  string
	apiKey string
	client *http.Client
}

type ihuyiResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	SmsID string `json:"smsid"`
}

const ihuyiAPI = "https://106.ihuyi.com/webservice/sms.php?method=Submit"

func NewIhuyiSender(apiID, apiKey string) Sender {
	return &ihuyiSender{
		apiID:  apiID,
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *ihuyiSender) Send(phone, code string) error {
	content := fmt.Sprintf("您的验证码是:%s。请不要把验证码泄露给其他人。", code)

	// Try plain key first (static password mode)
	if err := s.trySend(phone, content, s.apiKey, ""); err == nil {
		return nil
	}

	// Retry with MD5 dynamic password (account + apikey + mobile + content + time)
	ts := fmt.Sprintf("%d", time.Now().Unix())
	hash := fmt.Sprintf("%x", md5.Sum([]byte(s.apiID+s.apiKey+phone+content+ts)))
	return s.trySend(phone, content, hash, ts)
}

func (s *ihuyiSender) trySend(phone, content, password, ts string) error {
	form := url.Values{
		"account":  {s.apiID},
		"password": {password},
		"mobile":   {phone},
		"content":  {content},
		"format":   {"json"},
	}
	if ts != "" {
		form.Set("time", ts)
	}

	resp, err := s.client.Post(ihuyiAPI, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("ihuyi send failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result ihuyiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("ihuyi decode failed: %s body=%s", err, string(body))
	}
	if result.Code != 2 {
		return fmt.Errorf("ihuyi api error: code=%d msg=%s", result.Code, result.Msg)
	}
	return nil
}
