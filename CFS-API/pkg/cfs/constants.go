package cfs

const BaseCFSQuery = `
SELECT 
		cfs.id,
		cfs.description,
		cfs.eventTime,
		cfs.location,
		location.latitude,
		location.longitude,
		location.ward,
		location.neighborhood,
		count(*) OVER() as fullCount 
		
		FROM cfs

		INNER JOIN location
		on cfs.location = location.location

		WHERE 1 = 1
`

const BaseCFSQueryOrder = "ORDER BY cfs.eventTime desc"
const defaultOffset = 0
const defaultLimit = 100
