package avatar

func (dto AvatarDTO) IsEmpty() bool {
	return dto.Label == "" && dto.Avatar == nil
}
