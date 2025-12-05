package core

import (
	eventdto "cbe-super-app-cps-action/internal/constants/dto/event"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func ParseEventRequestFromMultipartForm(r *http.Request, isCreate bool) (eventdto.EventRequest, error) {
	var req eventdto.EventRequest

	_, fileHeader, err := utils.ParseMultipartFormFile(r, "cover_image", 15<<20)
	if err != nil {
		if err.Error() != localization.ErrorMissingFile.Code || isCreate {
			return req, errors.New(localization.ErrorInvalidFileUpload.Code)
		}
	} else {
		req.CoverImage = fileHeader
	}

	req.MerchantID = r.FormValue("merchant_id")
	req.EventName = r.FormValue("event_name")
	req.EventVenue = r.FormValue("event_venue")
	req.EventCity = r.FormValue("event_city")
	req.EventDescription = r.FormValue("event_description")

	if start := r.FormValue("start_date"); start != "" {
		t, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return req, errors.New(localization.ErrorInvalidDateFormat.Code)
		}
		req.StartDate = t
	}

	if due := r.FormValue("due_date"); due != "" {
		t, err := time.Parse(time.RFC3339, due)
		if err != nil {
			return req, errors.New(localization.ErrorInvalidDateFormat.Code)
		}
		req.DueDate = t
	}

	if total := r.FormValue("total_ticket_count"); total != "" {
		n, err := strconv.ParseUint(total, 10, 32)
		if err != nil {
			return req, errors.New(localization.ErrorInvalidNumberFormat.Code)
		}
		req.TotalTicketCount = uint(n)
	}

	tickets, err := parseTicketsFromForm(r.Form)
	if err != nil {
		return req, err
	}
	req.Tickets = tickets

	return req, nil
}

func parseTicketsFromForm(form url.Values) ([]types.Ticket, error) {
	var tickets []types.Ticket
	for i := 0; ; i++ {
		name := form.Get("ticket_name[" + strconv.Itoa(i) + "]")
		if name == "" {
			break
		}

		category := form.Get("ticket_category[" + strconv.Itoa(i) + "]")
		tType := form.Get("ticket_type[" + strconv.Itoa(i) + "]")
		priceStr := form.Get("ticket_price[" + strconv.Itoa(i) + "]")
		numStr := form.Get("ticket_number_of_ticker[" + strconv.Itoa(i) + "]")

		price, err := strconv.ParseUint(priceStr, 10, 64)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidNumberFormat.Code)
		}

		number, err := strconv.ParseUint(numStr, 10, 8)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidNumberFormat.Code)
		}

		tickets = append(tickets, types.Ticket{
			Name:           name,
			Category:       category,
			Type:           tType,
			Price:          price,
			NumberOfTicker: uint8(number),
		})
	}
	return tickets, nil
}
