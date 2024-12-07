package models

type AddSongRequest struct {
	Name        string `json:"name" binding:"required"`
	GroupName   string `json:"group_name" binding:"required"`
	ReleaseDate string `json:"release_date" binding:"required"`
	Text        string `json:"text" binding:"required"`
	Link        string `json:"link" binding:"required"`
}

type GetSongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

type Pagination struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

type GetSongLyricsRequest struct {
	ID string `uri:"id" binding:"required"`
}

type DeleteSongRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateSongRequest struct {
	Id          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	GroupName   string `json:"group_name" binding:"required"`
	ReleaseDate string `json:"release_date" binding:"required"`
	Text        string `json:"text" binding:"required"`
	Link        string `json:"link" binding:"required"`
}
