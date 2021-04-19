package cfs

import "errors"

func validateLimitOffset(search SearchQuery) (SearchQuery, error) {

	if search.Limit == nil {
		limit := int(defaultLimit)
		search.Limit = &limit
	} else if *search.Limit < 1 || *search.Limit > 500 {
		return search, errors.New("Page limit must be between 1 and 500.")
	}

	if search.Offset == nil {
		offset := int(defaultOffset)
		search.Offset = &offset
	} else if *search.Offset < 0 {
		return search, errors.New("Offset must not be a negative number.")
	}

	return search, nil
}
