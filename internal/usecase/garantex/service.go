package garantex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"test-task/internal/entity"
)

type Repository interface {
	RateCreate(ctx context.Context, timestamp int64, ask, bid entity.Rate) error
}

type UseCase struct {
	repo       Repository
	httpClient *http.Client
}

func New(repo Repository) *UseCase {
	return &UseCase{
		repo:       repo,
		httpClient: &http.Client{},
	}
}

func (uc *UseCase) GetRates(ctx context.Context, market string) (*entity.DepthAPIResponse, error) {
	endpoint := fmt.Sprintf("https://garantex.org/api/v2/depth?market=%s", market)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		return nil, fmt.Errorf("http newRequest: %w", err)
	}

	res, err := uc.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http client: %w", err)
	}

	defer res.Body.Close()

	var response entity.DepthAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	if len(response.Asks) >= 1 && len(response.Bids) >= 1 {
		err := uc.repo.RateCreate(ctx, response.Timestamp, response.Asks[0], response.Bids[0])
		if err != nil {
			return nil, err
		}
	}

	return &response, nil
}
