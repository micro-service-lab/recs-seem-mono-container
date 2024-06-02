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

	// RequestFormatError リクエストフォーマットエラー
	RequestFormatError = APIResponseType{
		Code:       "225",
		Message:    "request format error",
		StatusCode: http.StatusBadRequest,
	}

	// AttemptOperatePersonalOrganization 個人用オーガナイゼーションの操作を試みる
	AttemptOperatePersonalOrganization = APIResponseType{
		Code:       "226",
		Message:    "attempt operate personal organization error",
		StatusCode: http.StatusForbidden,
	}

	// AttemptOperateWholeOrganization 全体用オーガナイゼーションの操作を試みる
	AttemptOperateWholeOrganization = APIResponseType{
		Code:       "227",
		Message:    "attempt operate whole organization error",
		StatusCode: http.StatusForbidden,
	}

	// AttemptOperateGroupOrganization グループ用オーガナイゼーションの操作を試みる
	AttemptOperateGroupOrganization = APIResponseType{
		Code:       "228",
		Message:    "attempt operate group organization error",
		StatusCode: http.StatusForbidden,
	}

	// AttemptOperateGradeOrganization 学年用オーガナイゼーションの操作を試みる
	AttemptOperateGradeOrganization = APIResponseType{
		Code:       "229",
		Message:    "attempt operate grade organization error",
		StatusCode: http.StatusForbidden,
	}

	// ConflictStorageKey ストレージキーの競合
	ConflictStorageKey = APIResponseType{
		Code:       "230",
		Message:    "conflict storage key error",
		StatusCode: http.StatusConflict,
	}

	// NotFileOwner ファイルの所有者でない
	NotFileOwner = APIResponseType{
		Code:       "231",
		Message:    "not file owner error",
		StatusCode: http.StatusForbidden,
	}

	// OnlyProfessorAction 教授のみの操作
	OnlyProfessorAction = APIResponseType{
		Code:       "232",
		Message:    "only professor action error",
		StatusCode: http.StatusForbidden,
	}

	// InvalidLoginIDOrPassword ログインIDまたはパスワードが無効
	InvalidLoginIDOrPassword = APIResponseType{
		Code:       "233",
		Message:    "invalid login id or password error",
		StatusCode: http.StatusUnauthorized,
	}

	// InvalidRefreshToken リフレッシュトークンが無効
	InvalidRefreshToken = APIResponseType{
		Code:       "234",
		Message:    "invalid refresh token error",
		StatusCode: http.StatusUnauthorized,
	}

	// ExpireAccessToken アクセストークンの有効期限切れ
	ExpireAccessToken = APIResponseType{
		Code:       "235",
		Message:    "expire access token error",
		StatusCode: http.StatusUnauthorized,
	}

	// ExpireRefreshToken リフレッシュトークンの有効期限切れ
	ExpireRefreshToken = APIResponseType{
		Code:       "236",
		Message:    "expire refresh token error",
		StatusCode: http.StatusUnauthorized,
	}

	// CannotDeleteOrganizationChatRoom 組織チャットルームを削除できない
	CannotDeleteOrganizationChatRoom = APIResponseType{
		Code:       "237",
		Message:    "cannot delete organization chat room error",
		StatusCode: http.StatusForbidden,
	}

	// CannotDeletePrivateChatRoom プライベートチャットルームを削除できない
	CannotDeletePrivateChatRoom = APIResponseType{
		Code:       "238",
		Message:    "cannot delete private chat room error",
		StatusCode: http.StatusForbidden,
	}

	// CannotUpdatePrivateChatRoom プライベートチャットルームを更新できない
	CannotUpdatePrivateChatRoom = APIResponseType{
		Code:       "239",
		Message:    "cannot update private chat room error",
		StatusCode: http.StatusForbidden,
	}

	// CannotAddMemberToPrivateChatRoom プライベートチャットルームにメンバーを追加できない
	CannotAddMemberToPrivateChatRoom = APIResponseType{
		Code:       "240",
		Message:    "cannot add member to private chat room error",
		StatusCode: http.StatusForbidden,
	}

	// CannotWithdrawMemberFromPrivateChatRoom プライベートチャットルームからメンバーを退会できない
	CannotWithdrawMemberFromPrivateChatRoom = APIResponseType{
		Code:       "241",
		Message:    "cannot withdraw member from private chat room error",
		StatusCode: http.StatusForbidden,
	}

	// CannotAddMemberToOrganizationChatRoom 組織チャットルームにメンバーを追加できない
	CannotAddMemberToOrganizationChatRoom = APIResponseType{
		Code:       "242",
		Message:    "cannot add member to organization chat room error",
		StatusCode: http.StatusForbidden,
	}

	// CannotWithdrawMemberFromOrganizationChatRoom 組織チャットルームからメンバーを退会できない
	CannotWithdrawMemberFromOrganizationChatRoom = APIResponseType{
		Code:       "243",
		Message:    "cannot withdraw member from organization chat room error",
		StatusCode: http.StatusForbidden,
	}

	// CannotRemoveMemberFromOrganizationChatRoom 組織チャットルームからメンバーを削除できない
	CannotRemoveMemberFromOrganizationChatRoom = APIResponseType{
		Code:       "244",
		Message:    "cannot remove member from organization chat room error",
		StatusCode: http.StatusForbidden,
	}

	// CannotRemoveMemberFromPrivateChatRoom プライベートチャットルームからメンバーを削除できない
	CannotRemoveMemberFromPrivateChatRoom = APIResponseType{
		Code:       "245",
		Message:    "cannot remove member from private chat room error",
		StatusCode: http.StatusForbidden,
	}

	// MultiPartFormParseError マルチパートフォームのパースエラー
	MultiPartFormParseError = APIResponseType{
		Code:       "246",
		Message:    "multi part form parse error",
		StatusCode: http.StatusBadRequest,
	}

	// CannotDeleteSystemFile システムファイルを削除できない
	CannotDeleteSystemFile = APIResponseType{
		Code:       "247",
		Message:    "cannot delete system file error",
		StatusCode: http.StatusForbidden,
	}

	// NotImageFile 画像ファイルでない
	NotImageFile = APIResponseType{
		Code:       "248",
		Message:    "not image file error",
		StatusCode: http.StatusBadRequest,
	}

	// NotMessageOwner メッセージの所有者でない
	NotMessageOwner = APIResponseType{
		Code:       "249",
		Message:    "not message owner error",
		StatusCode: http.StatusForbidden,
	}

	// NotMatchChatRoomMessage チャットルームとメッセージが一致しない
	NotMatchChatRoomMessage = APIResponseType{
		Code:       "250",
		Message:    "not match chat room message error",
		StatusCode: http.StatusForbidden,
	}

	// CannotReadOwnMessage 自分のメッセージを読めない
	CannotReadOwnMessage = APIResponseType{
		Code:       "251",
		Message:    "cannot read own message error",
		StatusCode: http.StatusForbidden,
	}

	// NotCreateMessageToSelf 自分自身にメッセージを作成できない
	NotCreateMessageToSelf = APIResponseType{
		Code:       "252",
		Message:    "not create message to self error",
		StatusCode: http.StatusForbidden,
	}

	// CannotAttachSystemFile システムファイルを添付できない
	CannotAttachSystemFile = APIResponseType{
		Code:       "253",
		Message:    "cannot attach system file error",
		StatusCode: http.StatusForbidden,
	}
)
