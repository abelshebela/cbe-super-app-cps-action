package ad

// nonEmptyString returns the new value if non-empty, otherwise the old value
func nonEmptyString(new, old string) string {
	if new != "" {
		return new
	}
	return old
}

// nonEmptyAdvertFor returns the new value if non-empty, otherwise the old value
func nonEmptyAdvertFor(new, old AdvertFor) AdvertFor {
	if new != "" {
		return new
	}
	return old
}

// nonEmptyAdvertDate returns the new date if non-zero, otherwise the old date
func nonEmptyAdvertDate(new, old AdvertDate) AdvertDate {
	result := old
	if !new.StartedAt.IsZero() {
		result.StartedAt = new.StartedAt
	}
	if !new.ExpiredAt.IsZero() {
		result.ExpiredAt = new.ExpiredAt
	}
	return result
}

// generateAdvert creates a copy of an advert
func generateAdvert(advert Advert) *Advert {
	return &Advert{
		ID:            advert.ID,
		Title:         advert.Title,
		Description:   advert.Description,
		BannerImage:   advert.BannerImage,
		AdvertFor:     advert.AdvertFor,
		Date:          advert.Date,
		Enabled:       advert.Enabled,
		IsDeleted:     advert.IsDeleted,
		CreatedAt:     advert.CreatedAt,
		LastUpdatedAt: advert.LastUpdatedAt,
		DeletedAt:     advert.DeletedAt,
	}
}