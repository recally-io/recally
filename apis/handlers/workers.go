package handlers

import (
	"context"
	"vibrain/apis"
	"vibrain/internal/worders"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func (h *Handler) GetApiToolsAlipanSearch(ctx context.Context, request apis.GetApiToolsAlipanSearchRequestObject) (apis.GetApiToolsAlipanSearchResponseObject, error) {
	req := webSearcherRequest{
		Query: request.Params.Query,
	}

	if err := validate.Struct(req); err != nil {
		return nil, newHttpError(400, err.Error())
	}

	content, err := worders.SearchAliPan(req.Query)
	if err != nil {
		return nil, newHttpError(500, err.Error())
	}

	items := make([]apis.PansearchResponseItemData, 0, len(content))
	for _, c := range content {
		items = append(items, apis.PansearchResponseItemData{
			Content: c.Content,
			Id:      c.Id,
			Image:   c.Image,
			Pan:     c.Pan,
			Time:    c.Time,
		})
	}

	return apis.GetApiToolsAlipanSearch200JSONResponse{
		Data:    items,
		Success: true,
	}, nil
}

type webReaderRequest struct {
	Url string `json:"url" validate:"required,url"`
}

func (h *Handler) GetApiToolsWebReader(ctx context.Context, request apis.GetApiToolsWebReaderRequestObject) (apis.GetApiToolsWebReaderResponseObject, error) {
	req := webReaderRequest{
		Url: request.Params.Url,
	}

	if err := validate.Struct(req); err != nil {
		return nil, newHttpError(400, err.Error())
	}

	content, err := worders.WebReader(req.Url)
	if err != nil {
		return nil, newHttpError(500, err.Error())
	}
	return apis.GetApiToolsWebReader200JSONResponse{
		Data: apis.WebReaderResponseData{
			Content:     content.Content,
			Description: &content.Description,
			Title:       content.Title,
			Url:         content.Url,
		},
		Success: true,
	}, nil
}

type webSearcherRequest struct {
	Query string `json:"query" validate:"required,min=1"`
}

func (h *Handler) GetApiToolsWebSearch(ctx context.Context, request apis.GetApiToolsWebSearchRequestObject) (apis.GetApiToolsWebSearchResponseObject, error) {
	req := webSearcherRequest{
		Query: request.Params.Query,
	}

	if err := validate.Struct(req); err != nil {
		return nil, newHttpError(400, err.Error())
	}

	content, err := worders.WebSearcher(req.Query)
	if err != nil {
		return nil, newHttpError(500, err.Error())
	}

	data := make([]apis.WebReaderResponseData, 0, len(content))
	for _, c := range content {
		data = append(data, apis.WebReaderResponseData{
			Content:     c.Content,
			Description: &c.Description,
			Title:       c.Title,
			Url:         c.Url,
		})
	}

	return apis.GetApiToolsWebSearch200JSONResponse{
		Data:    data,
		Success: true,
	}, nil
}

func (h *Handler) GetApiToolsWebSummary(ctx context.Context, request apis.GetApiToolsWebSummaryRequestObject) (apis.GetApiToolsWebSummaryResponseObject, error) {
	req := webReaderRequest{
		Url: request.Params.Url,
	}

	if err := validate.Struct(req); err != nil {
		return nil, newHttpError(400, err.Error())
	}

	content, err := worders.WebSummary(req.Url)
	if err != nil {
		return nil, newHttpError(500, err.Error())
	}

	return apis.GetApiToolsWebSummary200JSONResponse{
		Data: apis.WebSummaryResponseData{
			Data: content,
		},
		Success: true,
	}, nil
}
