package ws

// RequestType はリクエストの種類を表す型。
type RequestType string

const (
	// RequestTypeOnlineMembers オンラインメンバーを取得するリクエスト。
	RequestTypeOnlineMembers RequestType = "online_members"
)

// Request リクエストを表す構造体。
type Request struct {
	RequestType RequestType `json:"request_type"`
	Data        any         `json:"data"`
}
