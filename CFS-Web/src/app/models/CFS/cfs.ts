import { Time } from "@angular/common";

export interface CFS {
    Id:  string,
    description: string,
    eventTime:  Time,
    location: string,
    latitude: number,
    longitude: number,
    ward: string,
    neighborhood: string
}

export interface CFSResponse {
    cfs: CFS[]
    meta: Meta[]
}

export interface Meta {
    offset: string
    pageNumber: number
    pageSize: number 
    totalRecords: number
    next: string
    prev: string
}

// type Meta struct {
// 	Offset       string `json:"offset"`
// 	PageNumber   int    `json:"pageNumber"`
// 	PageSize     int    `json:"pageSize"`
// 	TotalRecords int    `json:"totalRecords"`
// 	Next         string `json:"next"`
// 	Prev         string `json:"prev"`
// }


// type CallForService struct {
// 	Description  string    `json:"description"`
// 	EventTime    time.Time `json:"eventTime"`
// 	Location     string    `json:"location"`
// 	Lat          float64   `json:"latitude"`
// 	Lng          float64   `json:"longitude"`
// 	Ward         string    `json:"ward"`
// 	Neighborhood string    `json:"neighborhood"`
// }