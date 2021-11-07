package repository

import (
	"database/sql"
	"fmt"
	"strings"
	internalError "yula/internal/error"
	"yula/internal/models"
	imageloader "yula/internal/pkg/image_loader"
	"yula/internal/pkg/search"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) search.SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

func (sr *SearchRepository) SelectWithFilter(search *models.SearchFilter, from, count int64) ([]*models.Advert, error) {
	nums := make([]interface{}, 0)
	vars := make([]interface{}, 0)
	queryStr := `
		SELECT * FROM (
			SELECT a.id "id", a.name "name_", a.description, a.price, a.location, a.latitude, a.longitude, a.published_at, 
				a.date_close, a.is_active, a.views, a.publisher_id, c.name "category", array_agg(ai.img_path), a.amount, a.is_new FROM advert a
			JOIN category c ON a.category_id = c.Id
			LEFT JOIN advert_image ai ON a.id = ai.advert_id 
			GROUP BY a.id, a.name, a.Description,  a.price, a.location, a.latitude, a.longitude, a.published_at, 
				a.date_close, a.is_active, a.views, a.publisher_id, c.name
		) as t
		WHERE plainto_tsquery($%d) @@ (to_tsvector(t.name_) || to_tsvector(t.description)) 
	`
	nums = append(nums, 1+len(nums))
	vars = append(vars, search.Query)

	if search.Category != "" {
		queryStr += " AND t.category LIKE $%d"
		nums = append(nums, 1+len(nums))
		vars = append(vars, search.Category)
	}

	if search.TimeDuration != models.TimeDurationNone {
		queryStr += " AND (SELECT EXTRACT(DAY FROM ($%d - t.published_at))) < $%d"
		nums = append(nums, 1+len(nums), 2+len(nums))
		vars = append(vars, search.Date, search.TimeDuration)
	}

	if search.Radius != models.RadiusNone && search.Latitude != models.LatitudeNone && search.Longitude != models.LongitudeNone {
		queryStr += ` AND ST_DWithin(Geography(ST_SetSRID(ST_POINT(longitude, latitude), 4326)),
									 Geography(ST_SetSRID(ST_POINT($%d, $%d), 4326)),
									 $%d)`
		nums = append(nums, 1+len(nums), 2+len(nums), 3+len(nums))
		vars = append(vars, search.Longitude, search.Latitude, search.Radius)
	}

	if search.SortingName {
		queryStr += " ORDER BY t.name_"
	}

	if search.SortingDate {
		if search.SortingName {
			queryStr += ", t.published_at"
		} else {
			queryStr += " ORDER BY t.published_at"
		}
	}

	queryStr = fmt.Sprintf(queryStr, nums...)
	query, err := sr.db.Query(queryStr, vars...)
	if err != nil {
		return nil, internalError.GenInternalError(err)
	}

	defer query.Close()
	adverts := make([]*models.Advert, 0)
	for query.Next() {
		var advert models.Advert
		var images string

		err = query.Scan(&advert.Id, &advert.Name, &advert.Description, &advert.Price, &advert.Location, &advert.Latitude,
			&advert.Longitude, &advert.PublishedAt, &advert.DateClose, &advert.IsActive, &advert.Views,
			&advert.PublisherId, &advert.Category, &images, &advert.Amount, &advert.IsNew)

		if err != nil {
			return nil, internalError.GenInternalError(err)
		}

		advert.Images = make([]string, 0)
		if images[1:len(images)-1] != "NULL" {
			advert.Images = strings.Split(images, ",")
		}

		if len(advert.Images) == 0 {
			advert.Images = append(advert.Images, imageloader.DefaultAdvertImage)
		}

		adverts = append(adverts, &advert)
	}
	return adverts, nil
}
