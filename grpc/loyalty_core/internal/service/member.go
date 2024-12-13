package service

import (
	"context"
	"loyalty_core/api"
)

func (s *Service) CreateMember(ctx context.Context, req *api.CreateMemberRequest) (*api.CreateMemberResponse, error) {
	s.log.Info("CreateMember")
	err, id := s.biz.MemberBusiness.ProcessCreateMember(ctx, req)
	if err != nil {
		return &api.CreateMemberResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	return &api.CreateMemberResponse{
		Success: true,
		Message: "Member created successfully",
		Data: &api.CreateMemberResponse_Data{
			MemberId: id,
		},
	}, nil
}
