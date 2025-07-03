package unlink

type UnlinkRepository interface {
	UnlinkDevice(userCode, makerUser string, branchCode []string, homeBranch string) error
	ApproveOrDecline(userCode, decision, reason, checkerUser string) error
}
