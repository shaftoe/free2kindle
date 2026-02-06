package http

type ArticleRequest struct {
	URL string `json:"url" binding:"required"`
}

type ArticleResponse struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
