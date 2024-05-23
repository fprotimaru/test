package grpc_server

import (
	"context"
	"test-task/internal/transport/grpc_contract"
	"test-task/internal/usecase/garantex"
)

var _ grpc_contract.APIServer = (*Server)(nil)

type Server struct {
	uc *garantex.UseCase
}

func New(uc *garantex.UseCase) *Server {
	return &Server{
		uc: uc,
	}
}

func (s *Server) GetRates(ctx context.Context, req *grpc_contract.GetRatesRequest) (*grpc_contract.GetRatesResponse, error) {
	res, err := s.uc.GetRates(ctx, req.GetMarket())
	if err != nil {
		return &grpc_contract.GetRatesResponse{}, err
	}

	var ask grpc_contract.Rate
	if len(res.Asks) >= 1 {
		ask = grpc_contract.Rate{
			Price:  res.Asks[0].Price,
			Volume: res.Asks[0].Volume,
			Amount: res.Asks[0].Amount,
			Factor: res.Asks[0].Factor,
			Type:   string(res.Asks[0].Type),
		}
	}

	var bid grpc_contract.Rate
	if len(res.Bids) >= 1 {
		bid = grpc_contract.Rate{
			Price:  res.Bids[0].Price,
			Volume: res.Bids[0].Volume,
			Amount: res.Bids[0].Amount,
			Factor: res.Bids[0].Factor,
			Type:   string(res.Bids[0].Type),
		}
	}

	return &grpc_contract.GetRatesResponse{
		Timestamp: res.Timestamp,
		Ask:       &ask,
		Bid:       &bid,
	}, nil
}
