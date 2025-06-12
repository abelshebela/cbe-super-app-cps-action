package ad

type AdvertResponse struct {
	ID          string
	Title       string
	Description string
	BannerImage string
	AdvertFor   AdvertFor
	Date        AdvertDate
}

type CreateAdvertRequest struct {
	Title       string
	Description string
	BannerImage string
	AdvertFor   AdvertFor
	Date        AdvertDate
}

type UpdateAdvertRequest struct {
	Title       string
	Description string
	BannerImage string
	AdvertFor   AdvertFor
	Date        AdvertDate
}
