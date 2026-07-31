package donation_oracle

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

func (r *repository) replaceImages(ctx context.Context, donationID string, images []shared_types.DonationImage) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if _, err := r.db.ExecContext(ctx,
		`DELETE FROM DONATION_IMAGES WHERE DONATION_ID = HEXTORAW(:1)`,
		donationID,
	); err != nil {
		log.Errorf("[DonationOracle][replaceImages] delete failed: %v", err)
		return local_util.HandleDBError(err)
	}

	insertChild := `INSERT INTO DONATION_IMAGES (DONATION_ID, PHOTO_URL) VALUES (HEXTORAW(:1), :2)`
	for _, img := range images {
		if strings.TrimSpace(img.PhotoURL) == "" {
			continue
		}
		if _, err := r.db.ExecContext(ctx, insertChild, donationID, img.PhotoURL); err != nil {
			log.Errorf("[DonationOracle][replaceImages] insert failed: %v", err)
			return local_util.HandleDBError(err)
		}
	}
	return nil
}

func (r *repository) fetchImagesFor(ctx context.Context, donationIDs []string) (map[string][]types.DonationImage, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	out := make(map[string][]types.DonationImage, len(donationIDs))
	if len(donationIDs) == 0 {
		return out, nil
	}

	in := make([]string, len(donationIDs))
	args := make([]interface{}, len(donationIDs))
	for i, id := range donationIDs {
		in[i] = fmt.Sprintf("HEXTORAW(:%d)", i+1)
		args[i] = id
	}
	q := `SELECT RAWTOHEX(ID), RAWTOHEX(DONATION_ID), PHOTO_URL, CREATED_AT
	      FROM DONATION_IMAGES
	      WHERE DONATION_ID IN (` + strings.Join(in, ",") + `)
	      ORDER BY CREATED_AT`

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		log.Errorf("[DonationOracle][fetchImagesFor] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var imgID, donationID, photoURL string
		var createdAt time.Time
		if err := rows.Scan(&imgID, &donationID, &photoURL, &createdAt); err != nil {
			log.Errorf("[DonationOracle][fetchImagesFor] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		out[donationID] = append(out[donationID], types.DonationImage{
			ID:        imgID,
			PhotoURL:  photoURL,
			CreatedAt: createdAt,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}
