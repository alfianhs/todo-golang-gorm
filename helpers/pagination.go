package helpers

import (
	"net/url"
	"strconv"
)

func GetOffsetLimit(query url.Values) (page, offset, limit int) {
	// validate offset & limit
	pageString := query.Get("page")
	limitString := query.Get("limit")
	pageInt, _ := strconv.Atoi(pageString)
	limitInt, _ := strconv.Atoi(limitString)

	// set default page & limit
	if pageInt <= 0 {
		pageInt = 1
	}
	if limitInt <= 0 {
		limitInt = 10
	}

	// set offset
	offset = (pageInt - 1) * limitInt

	return pageInt, offset, limitInt
}
