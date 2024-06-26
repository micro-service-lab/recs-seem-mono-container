package response

import "net/http"

// Code レスポンスコード
type Code string

// APIResponseType API レスポンスの種別
type APIResponseType struct {
	Code       Code
	Message    string
	StatusCode int
}

// String Code を文字列に変換して返す。
func (c Code) String() string {
	return string(c)
}

// String APIResponseType を文字列に変換して返す。
func (r APIResponseType) String() string {
	return r.Code.String() + ": " + r.Message
}

var (
	// Success 成功
	Success = APIResponseType{
		Code:       "000",
		Message:    "success",
		StatusCode: http.StatusOK,
	}
	// System システムエラー
	System = APIResponseType{
		Code:       "100",
		Message:    "system error",
		StatusCode: http.StatusInternalServerError,
	}

	// Validation バリデーションエラー
	Validation = APIResponseType{
		Code:       "200",
		Message:    "validation error",
		StatusCode: http.StatusBadRequest,
	}

	// Permission 権限エラー
	Permission = APIResponseType{
		Code:       "201",
		Message:    "permission error",
		StatusCode: http.StatusForbidden,
	}

	// Unauthorized 承認エラー
	Unauthorized = APIResponseType{
		Code:       "202",
		Message:    "unauthorized error",
		StatusCode: http.StatusUnauthorized,
	}

	// NotFound リソースが見つからない
	NotFound = APIResponseType{
		Code:       "203",
		Message:    "not found error",
		StatusCode: http.StatusNotFound,
	}

	// Unauthenticated 認証エラー
	Unauthenticated = APIResponseType{
		Code:       "204",
		Message:    "unauthenticated error",
		StatusCode: http.StatusUnauthorized,
	}

	// PostTooLarge リクエストが大きすぎる
	PostTooLarge = APIResponseType{
		Code:       "205",
		Message:    "post too large error",
		StatusCode: http.StatusRequestEntityTooLarge,
	}

	// ThrottleRequests リクエストが多すぎる
	ThrottleRequests = APIResponseType{
		Code:       "206",
		Message:    "throttle requests error",
		StatusCode: http.StatusTooManyRequests,
	}

	// InvalidSignature 無効な署名
	InvalidSignature = APIResponseType{
		Code:       "207",
		Message:    "invalid signature error",
		StatusCode: http.StatusUnauthorized,
	}

	// StreamedResponse ストリームレスポンス
	StreamedResponse = APIResponseType{
		Code:       "208",
		Message:    "streamed response error",
		StatusCode: http.StatusInternalServerError,
	}

	// TokenMismatch トークン不一致
	TokenMismatch = APIResponseType{
		Code:       "209",
		Message:    "token mismatch error",
		StatusCode: http.StatusForbidden,
	}

	// MethodNotAllowed 許可されていないメソッド
	MethodNotAllowed = APIResponseType{
		Code:       "210",
		Message:    "method not allowed error",
		StatusCode: http.StatusMethodNotAllowed,
	}

	// NotFoundModel モデルが見つからない
	NotFoundModel = APIResponseType{
		Code:       "211",
		Message:    "not found model error",
		StatusCode: http.StatusNotFound,
	}

	// TokenBlacklisted トークンがブラックリストに登録されている
	TokenBlacklisted = APIResponseType{
		Code:       "212",
		Message:    "token blacklisted error",
		StatusCode: http.StatusUnauthorized,
	}

	// ModelConflict モデルの競合
	ModelConflict = APIResponseType{
		Code:       "213",
		Message:    "model conflict error",
		StatusCode: http.StatusConflict,
	}

	// GuestGuard ゲストガードエラー
	GuestGuard = APIResponseType{
		Code:       "214",
		Message:    "guest guard error",
		StatusCode: http.StatusForbidden,
	}

	// UserOnly ユーザーのみ
	UserOnly = APIResponseType{
		Code:       "215",
		Message:    "user only error",
		StatusCode: http.StatusForbidden,
	}

	// ThrottleLoginRequests ログインリクエストが多すぎる
	ThrottleLoginRequests = APIResponseType{
		Code:       "216",
		Message:    "throttle login requests error",
		StatusCode: http.StatusTooManyRequests,
	}

	// FailedUpload アップロードに失敗
	FailedUpload = APIResponseType{
		Code:       "217",
		Message:    "failed upload error",
		StatusCode: http.StatusInternalServerError,
	}

	// AuthNotFound 認証が見つからない
	AuthNotFound = APIResponseType{
		Code:       "218",
		Message:    "auth not found error",
		StatusCode: http.StatusNotFound,
	}

	// RefreshTokenExpired リフレッシュトークンの有効期限切れ
	RefreshTokenExpired = APIResponseType{
		Code:       "219",
		Message:    "refresh token expired error",
		StatusCode: http.StatusUnauthorized,
	}

	// AlreadyLogout 既にログアウト済み
	AlreadyLogout = APIResponseType{
		Code:       "220",
		Message:    "already logout error",
		StatusCode: http.StatusUnauthorized,
	}

	// SQLQueryError SQL クエリエラー
	SQLQueryError = APIResponseType{
		Code:       "221",
		Message:    "sql query error",
		StatusCode: http.StatusInternalServerError,
	}

	// AlreadyNotExist 既に存在しない
	AlreadyNotExist = APIResponseType{
		Code:       "222",
		Message:    "already not exist error",
		StatusCode: http.StatusNotFound,
	}

	// NotMatchKey キーが一致しない
	NotMatchKey = APIResponseType{
		Code:       "223",
		Message:    "not match key error",
		StatusCode: http.StatusForbidden,
	}

	// UnsupportedMediaType サポートされていないメディアタイプ
	UnsupportedMediaType = APIResponseType{
		Code:       "224",
		Message:    "unsupported media type error",
		StatusCode: http.StatusUnsupportedMediaType,
	}
)
