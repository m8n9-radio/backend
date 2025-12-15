package dto

type CreateTrackRequest struct {
	MD5Sum string `json:"md5sum" validate:"required,len=32,hexadecimal"`
	Title  string `json:"title" validate:"required,min=1,max=500"`
	Cover  string `json:"cover" validate:"required,url"`
}

type TrackResponse struct {
	ID        string `json:"id"`
	MD5Sum    string `json:"md5sum"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Rotate    int    `json:"rotate"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
