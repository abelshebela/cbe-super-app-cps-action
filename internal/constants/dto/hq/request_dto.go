package hqDto

type UpdateBlockTimeRequest struct {
	BlockTime uint32 `json:"block_time"`
}

type UpdateArchiveTimeRequest struct {
	ArchiveTime uint32 `json:"archive_time"`
}

type UpdatePasswordExpiryRequest struct {
	PasswordExpiry uint32 `json:"password_expiry"`
}
