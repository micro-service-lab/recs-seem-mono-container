package entity

// AttachedMessageOnMessage メッセージに添付された添付を表す構造体。
type AttachedMessageOnMessage struct {
	AttachableItem AttachableItemWithMimeType `json:"attachable_item"`
}

// AttachableItemWithOnChatRoom チャットルームに添付された添付を表す構造体。
type AttachableItemWithOnChatRoom struct {
	AttachableItem AttachableItemWithMimeType `json:"attachable_item"`
}
