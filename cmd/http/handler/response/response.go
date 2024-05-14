// Package response provides a application response.
package response

// ApplicationResponse 返却するレスポンスの型
type ApplicationResponse struct {
	Success         bool                       `json:"success"` // 成功(2xx)ならtrue, 失敗(4xx, 5xx)ならfalse
	Data            any                        `json:"data"`    // 成功ならnull以外、空のデータでも空配列で返却される
	Code            Code                       `json:"code"`
	Message         string                     `json:"message"`          // エラーコードごとのメッセージ
	ErrorAttributes ApplicationErrorAttributes `json:"error_attributes"` // エラー属性の追加付与
}

// ApplicationErrorAttributes エラー属性による追加情報
type ApplicationErrorAttributes any
