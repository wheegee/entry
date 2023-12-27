package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type MockSSMClient struct {
	GetParametersByPathFunc func(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
	GetParameterFunc        func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

func (m *MockSSMClient) GetParametersByPath(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error) {
	return m.GetParametersByPathFunc(ctx, params, optFns...)
}

func (m *MockSSMClient) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return m.GetParameterFunc(ctx, params, optFns...)
}
