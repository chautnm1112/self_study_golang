package service

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
)

func (s *Service) CreateMember(ctx context.Context, req *api.CreateMemberRequest) (*api.CreateMemberResponse, error) {
	s.log.Named("SERVICE-MEMBER").Info("CreateMember start")

	err, id := s.biz.MemberBusiness.ProcessCreateMember(ctx, req)

	if err != nil {
		return &api.CreateMemberResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	s.log.Named("SERVICE-MEMBER").Info("CreateMember end")

	return &api.CreateMemberResponse{
		Success: true,
		Message: "Member created successfully",
		Data: &api.CreateMemberResponse_Data{
			MemberId: id,
		},
	}, nil
}
