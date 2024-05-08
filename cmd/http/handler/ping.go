package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
)

// PingRequest ping リクエストを表す。
type PingRequest struct {
	// Message 任意の文字列
	Message string `json:"message"`
}

// PingResponse ping レスポンスを表す。
type PingResponse struct {
	// Message 受信した文字列
	Message string `json:"message"`
	// ReceivedTime サーバー受信時刻
	ReceivedTime time.Time `json:"receivedTime"`
}

// PingHandler 疎通確認 API の HTTP ハンドラ。
func PingHandler(clk clock.Clock) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 受信時刻を取得する。
		// 現在時刻を参照するときは time.Now() ではなく、s.clk.Now() を使用するようにする。
		// これにより単体テスト等で時刻を偽装することが容易になる。
		receivedTime := clk.Now()

		// リクエストボディをデコードする。
		req := &PingRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// デコードに失敗した場合はログ出力して 400 Bad Request を返す。
			log.Printf("[ERROR] request decoding failed: %+v", err)
			errAtr := response.ApplicationErrorAttributes{
				"error": "invalid json",
			}
			err := response.JSONResponseWriter(r.Context(), w, response.Validation, nil, errAtr)
			if err != nil {
				log.Printf("[ERROR] response writing failed: %+v", err)
			}
			return
		}

		// レスポンスボディを表す構造体を生成する。
		resp := &PingResponse{
			Message:      req.Message,
			ReceivedTime: receivedTime,
		}

		err := response.JSONResponseWriter(r.Context(), w, response.Success, resp, nil)
		if err != nil {
			log.Printf("[ERROR] response writing failed: %+v", err)
		}
	}
}
