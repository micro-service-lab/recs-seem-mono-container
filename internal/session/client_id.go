package session

import "strconv"

// ToClientID int32 の ID をクライアント ID として使うための文字列に変換する。
func ToClientID(id int32) string {
	if id == 0 {
		return ""
	}

	return strconv.FormatInt(int64(id), 10)
}

// FromClientID クライアント ID を int32 に変換する。
// 変換できない場合は 0 を返す。
func FromClientID(clientID string) int32 {
	id, err := strconv.ParseInt(clientID, 10, 32)
	if err != nil {
		return 0
	}

	return int32(id)
}
