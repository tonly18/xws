/* *
 * Test Service
 */

package service

import (
	"context"
	"fmt"

	"github.com/extra-time-zone/xws/core/xerror"
	"github.com/extra-time-zone/xws/example/model"
)

type TestService struct {
	ctx context.Context
}

func NewTestService(ctx context.Context) *TestService {
	return &TestService{
		ctx: ctx,
	}
}

func (s *TestService) GetData(x int) (string, xerror.Error) {
	testModel := model.NewTestModel(s.ctx)
	data, xerr := testModel.GetData(x)
	if xerr != nil {
		return "", xerror.Wrap(xerr, "test-service-error")
	}

	return fmt.Sprintf("%s_%s", data, "test-service"), nil
}
