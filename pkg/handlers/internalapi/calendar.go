package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	calendarop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/calendar"
	"github.com/transcom/mymove/pkg/handlers"
	"time"
)

// ShowUnavailableMoveDatesHandler returns the unavailable move dates starting at a given date.
type ShowUnavailableMoveDatesHandler struct {
	handlers.HandlerContext
}

// Handle returns the unavailable move dates.
func (h ShowUnavailableMoveDatesHandler) Handle(params calendarop.ShowUnavailableMoveDatesParams) middleware.Responder {
	var datesPayload []strfmt.Date

	startDate := time.Time(params.StartDate)
	const daysToCheck = 90
	const shortFuseTotalDays = 5
	daysChecked := 0
	shortFuseDaysFound := 0

	for d := startDate; daysChecked < daysToCheck; d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			datesPayload = append(datesPayload, strfmt.Date(d))
		} else if shortFuseDaysFound < shortFuseTotalDays {
			datesPayload = append(datesPayload, strfmt.Date(d))
			shortFuseDaysFound++
		}
		daysChecked++
	}

	return calendarop.NewShowUnavailableMoveDatesOK().WithPayload(datesPayload)
}
