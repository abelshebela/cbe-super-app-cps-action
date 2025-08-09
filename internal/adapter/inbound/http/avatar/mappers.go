package avatar

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
)

func ToAvatarRequestFromDTO(dto *AvatarDTO) *avatar.AvatarRequest {
	return &avatar.AvatarRequest{
		Label:  dto.Label,
		Avatar: dto.Avatar,
	}
}

func ToAvatarResponseDTO(domain *avatar.Avatar) *AvatarResponseDTO {
	return &AvatarResponseDTO{
		ID:             domain.ID,
		Avatar:         domain.Avatar,
		Label:          domain.Label,
		Enable:         domain.Enable,
		IsDeleted:      domain.IsDeleted,
		CreatedAt:      domain.CreatedAt,
		LastModifiedAt: domain.LastModifiedAt,
		DeletedAt:      domain.DeletedAt,
	}
}
