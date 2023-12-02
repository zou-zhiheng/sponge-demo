package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"user/internal/model"

	userV1 "user/api/user/v1"
	"user/internal/cache"
	"user/internal/dao"
	"user/internal/ecode"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

func init() {
	registerFns = append(registerFns, func(server *grpc.Server) {
		userV1.RegisterUserServer(server, NewUserServer()) // register service to the rpc service
	})
}

var _ userV1.UserServer = (*user)(nil)

type user struct {
	userV1.UnimplementedUserServer

	iDao dao.UserDao
}

// NewUserServer create a new service
func NewUserServer() userV1.UserServer {
	return &user{
		iDao: dao.NewUserDao(
			model.GetDB(),
			cache.NewUserCache(model.GetCacheType()),
		),
	}
}

// Create a record
func (s *user) Create(ctx context.Context, req *userV1.CreateUserRequest) (*userV1.CreateUserReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record := &model.User{}
	err = copier.Copy(record, req)
	if err != nil {
		return nil, ecode.StatusCreateUser.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	err = s.iDao.Create(ctx, record)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("user", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &userV1.CreateUserReply{Id: record.ID}, nil
}

// DeleteByID delete a record by id
func (s *user) DeleteByID(ctx context.Context, req *userV1.DeleteUserByIDRequest) (*userV1.DeleteUserByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	err = s.iDao.DeleteByID(ctx, req.Id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &userV1.DeleteUserByIDReply{}, nil
}

// DeleteByIDs delete records by batch id
func (s *user) DeleteByIDs(ctx context.Context, req *userV1.DeleteUserByIDsRequest) (*userV1.DeleteUserByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	err = s.iDao.DeleteByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("ids", req.Ids), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &userV1.DeleteUserByIDsReply{}, nil
}

// UpdateByID update a record by id
func (s *user) UpdateByID(ctx context.Context, req *userV1.UpdateUserByIDRequest) (*userV1.UpdateUserByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record := &model.User{}
	err = copier.Copy(record, req)
	if err != nil {
		return nil, ecode.StatusUpdateByIDUser.Err()
	}
	record.ID = req.Id

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	err = s.iDao.UpdateByID(ctx, record)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("user", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &userV1.UpdateUserByIDReply{}, nil
}

// GetByID get a record by id
func (s *user) GetByID(ctx context.Context, req *userV1.GetUserByIDRequest) (*userV1.GetUserByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	fmt.Println("GetById---in")
	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	userMapper := model.GetUserMapper()
	res, err := userMapper.SelectById(req.Id)
	if err != nil {
		logger.Warn("GetByID error: User is not exist", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusNotFound.Err()
	}

	fmt.Println(res)

	//record, err := s.iDao.GetByID(ctx, req.Id)
	//if err != nil {
	//	if errors.Is(err, query.ErrNotFound) {
	//		logger.Warn("GetByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
	//		return nil, ecode.StatusNotFound.Err()
	//	}
	//	logger.Error("GetByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
	//	return nil, ecode.StatusInternalServerError.ToRPCErr()
	//}
	//
	//data, err := convertUser(record)
	//if err != nil {
	//	logger.Warn("convertUser error", logger.Err(err), logger.Any("user", record), interceptor.ServerCtxRequestIDField(ctx))
	//	return nil, ecode.StatusGetByIDUser.Err()
	//}

	//return &userV1.GetUserByIDReply{User: data}, nil
	return &userV1.GetUserByIDReply{User: &userV1.User{
		//Id:   res.ID,
		//Name: res.Name,
		//Age:  int32(res.Age),
	}}, nil
}

// GetByCondition get a record by id
func (s *user) GetByCondition(ctx context.Context, req *userV1.GetUserByConditionRequest) (*userV1.GetUserByConditionReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	conditions := &query.Conditions{}
	for _, v := range req.Conditions.GetColumns() {
		column := query.Column{}
		_ = copier.Copy(&column, v)
		conditions.Columns = append(conditions.Columns, column)
	}
	err = conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error", logger.Err(err), logger.Any("conditions", conditions), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record, err := s.iDao.GetByCondition(ctx, conditions)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByCondition error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusNotFound.Err()
		}
		logger.Error("GetByCondition error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	data, err := convertUser(record)
	if err != nil {
		logger.Warn("convertUser error", logger.Err(err), logger.Any("user", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusGetByConditionUser.Err()
	}

	return &userV1.GetUserByConditionReply{
		User: data,
	}, nil
}

// ListByIDs list of records by batch id
func (s *user) ListByIDs(ctx context.Context, req *userV1.ListUserByIDsRequest) (*userV1.ListUserByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	userMap, err := s.iDao.GetByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("ids", req.Ids), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	users := []*userV1.User{}
	for _, id := range req.Ids {
		if v, ok := userMap[id]; ok {
			record, err := convertUser(v)
			if err != nil {
				logger.Warn("convertUser error", logger.Err(err), logger.Any("user", v), interceptor.ServerCtxRequestIDField(ctx))
				return nil, ecode.StatusInternalServerError.ToRPCErr()
			}
			users = append(users, record)
		}
	}

	return &userV1.ListUserByIDsReply{Users: users}, nil
}

// List of records by query parameters
func (s *user) List(ctx context.Context, req *userV1.ListUserRequest) (*userV1.ListUserReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	params := &query.Params{}
	err = copier.Copy(params, req.Params)
	if err != nil {
		return nil, ecode.StatusListUser.Err()
	}
	params.Size = int(req.Params.Limit)

	ctx = context.WithValue(ctx, interceptor.ContextRequestIDKey, interceptor.ServerCtxRequestID(ctx)) //nolint
	records, total, err := s.iDao.GetByColumns(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "query params error:") {
			logger.Warn("GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusInvalidParams.Err()
		}
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	users := []*userV1.User{}
	for _, record := range records {
		data, err := convertUser(record)
		if err != nil {
			logger.Warn("convertUser error", logger.Err(err), logger.Any("id", record.ID), interceptor.ServerCtxRequestIDField(ctx))
			continue
		}
		users = append(users, data)
	}

	return &userV1.ListUserReply{
		Total: total,
		Users: users,
	}, nil
}

func convertUser(record *model.User) (*userV1.User, error) {
	value := &userV1.User{}
	err := copier.Copy(value, record)
	if err != nil {
		return nil, err
	}
	value.Id = record.ID
	value.CreateTime = record.CreateTime.Format("2006-01-02 15:04:05")
	value.UpdateTime = record.UpdateTime.Format("2006-01-02 15:04:05")

	return value, nil
}
