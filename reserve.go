package cbesupperappmemberauth2

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/user/core"
)

func (us *UsersService) UpdateProfilePicture(ctx context.Context, userID string, req dto.UpdateProfilePicture) error {
	var success bool
	if _, err := us.userRepo.FindById(ctx, req.UserId); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(constants.Empty, constants.ProfileTemp)
	if err != nil {
		return errors.ErrFailedToCreateTemp
	}

	defer func() {
		if err := tempFile.Close(); err != nil {
			us.logger.Warnf("Failed to close temp file: %v", err)
		}
		if !success {
			if err := os.Remove(tempFile.Name()); err != nil {
				us.logger.Warnf("Failed to remove temp file: %v", err)
			}
		}
	}()

	if _, err := io.Copy(tempFile, req.File); err != nil {
		us.logger.Errorf("Failed to copy file content to temp file: %v", err)
		return errors.ErrFailedToUpload
	}

	objectName := fmt.Sprintf("profile-pictures/%s/%s", id, filepath.Base(req.ProfilePicture.Filename))

	ProfileUrl, err := core.FileBucketUploader(ctx, us.minioServer, constants.BucketUserProfilePicture, objectName, filepath.Base(req.ProfilePicture.Filename))
	if err != nil {
		return err
	}

	if err := os.Remove(tempFile.Name()); err != nil {
		us.logger.Warnf("Failed to remove temp file after successful upload error: %v", err)
	}

	if err := us.userRepo.Update(ctx, req.UserId, &model.User{Avatar: ProfileUrl}); err != nil {
		return errors.ErrProfileSet
	}

	return nil
}
