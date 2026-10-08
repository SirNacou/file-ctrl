package handlers

type ListFilesRequest struct{}
type ListFilesResponse struct{}

type ListFilesHandler struct{}

func NewListFilesHandler() *ListFilesHandler {
	return &ListFilesHandler{}
}

func (h *ListFilesHandler) Handle(req *ListFilesRequest) (*ListFilesResponse, error) {
	return &ListFilesResponse{}, nil
}
