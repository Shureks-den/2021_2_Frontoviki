package models

import (
	"net/url"
	"strconv"
	"time"
	internalError "yula/internal/error"
)

const (
	TimeDurationNone int64   = -1
	LatitudeNone     float64 = -80
	LongitudeNone    float64 = 80
	RadiusNone       int64   = -1
)

type SearchFilter struct {
	Query        string    `valid:"type(string)"`
	Category     string    `valid:"type(string),optional"`
	Date         time.Time `valid:"-"`
	TimeDuration int64     `valid:"in(-1|1|3|7|30),optional"`
	Latitude     float64   `valid:"latitude,optional"`
	Longitude    float64   `valid:"longitude,optional"`
	Radius       int64     `valid:"int,optional"`
	SortingDate  bool      `valid:"optional"`
	SortingName  bool      `valid:"optional"`
}

func NewSearchFilter(values *url.Values) (*SearchFilter, error) {
	sf := &SearchFilter{}
	query := values.Get("query")
	if query == "" {
		return nil, internalError.BadRequest
	}
	sf.Query = query
	sf.Category = values.Get("category")
	sf.Date = time.Now()
	parse := values.Get("time_duration")
	switch parse {
	case "":
		sf.TimeDuration = TimeDurationNone
	default:
		tmp, err := strconv.ParseInt(parse, 10, 64)
		if err != nil {
			return nil, internalError.BadRequest
		}
		sf.TimeDuration = tmp
	}

	parse = values.Get("longitude")
	switch parse {
	case "":
		sf.Longitude = LongitudeNone
	default:
		tmp, err := strconv.ParseFloat(parse, 64)
		if err != nil {
			return nil, internalError.BadRequest
		}
		sf.Longitude = tmp
	}

	parse = values.Get("latitude")
	switch parse {
	case "":
		sf.Latitude = LatitudeNone
	default:
		tmp, err := strconv.ParseFloat(parse, 64)
		if err != nil {
			return nil, internalError.BadRequest
		}
		sf.Latitude = tmp
	}

	parse = values.Get("radius")
	switch parse {
	case "":
		sf.Radius = RadiusNone
	default:
		tmp, err := strconv.ParseInt(parse, 10, 64)
		if err != nil {
			return nil, internalError.BadRequest
		}
		sf.Radius = tmp
	}

	parse = values.Get("sorting_name")
	switch parse {
	case "":
		sf.SortingName = false
	default:
		tmp, err := strconv.ParseBool(parse)
		if err != nil {
			return nil, internalError.BadRequest
		}
		sf.SortingName = tmp
	}

	parse = values.Get("sorting_date")
	switch parse {
	case "":
		sf.SortingDate = false
	default:
		tmp, err := strconv.ParseBool(parse)
		if err != nil {
			return nil, internalError.BadRequest
		}
		sf.SortingDate = tmp
	}

	return sf, nil
}
