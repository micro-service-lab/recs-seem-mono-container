package errhandle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/testutils"
)

var successData = map[string]any{
	"message": "hello world",
}

func TestErrorHandle(t *testing.T) {
	t.Parallel()
	type wants struct {
		resType response.APIResponseType
		errAttr map[string]any
		handled bool
	}
	cases := map[string]struct {
		err  error
		want wants
	}{
		"wantErr": {
			err: wantErr{},
			want: wants{
				resType: response.ModelConflict,
				errAttr: nil,
				handled: true,
			},
		},
		"wrap wantErr": {
			err: fmt.Errorf("wrap: %w", wantErr{}),
			want: wants{
				resType: response.ModelConflict,
				errAttr: nil,
				handled: true,
			},
		},
		"wrap wrap wantErr": {
			err: fmt.Errorf("wrap: %w", fmt.Errorf("wrap: %w", wantErr{})),
			want: wants{
				resType: response.ModelConflict,
				errAttr: nil,
				handled: true,
			},
		},
		"wantErrWithAttr": {
			err: wantErrWithAttr{
				attr: map[string]any{
					"key": "value",
				},
			},
			want: wants{
				resType: response.Validation,
				errAttr: map[string]any{
					"key": "value",
				},
				handled: true,
			},
		},
		"wrap wantErrWithAttr": {
			err: fmt.Errorf("wrap: %w", wantErrWithAttr{
				attr: map[string]any{
					"key": "value",
				},
			}),
			want: wants{
				resType: response.Validation,
				errAttr: map[string]any{
					"key": "value",
				},
				handled: true,
			},
		},
		"wrap wrap wantErrWithAttr": {
			err: fmt.Errorf("wrap: %w", fmt.Errorf("wrap: %w", wantErrWithAttr{
				attr: map[string]any{
					"key": "value",
				},
			})),
			want: wants{
				resType: response.Validation,
				errAttr: map[string]any{
					"key": "value",
				},
				handled: true,
			},
		},
		"nil": {
			err: nil,
			want: wants{
				resType: response.Success,
				errAttr: nil,
				handled: false,
			},
		},
		"not application error": {
			err: notApplicationErr{},
			want: wants{
				resType: response.System,
				errAttr: nil,
				handled: true,
			},
		},
		"wrap not application error": {
			err: fmt.Errorf("wrap: %w", notApplicationErr{}),
			want: wants{
				resType: response.System,
				errAttr: nil,
				handled: true,
			},
		},
	}
	ctx := context.Background()

	for name, c := range cases {
		cc := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			var want response.ApplicationResponse
			var data any
			if cc.err == nil {
				data = successData
				err := response.JSONResponseWriter(ctx, w, response.Success, data, nil)
				assert.NoError(t, err)
			}
			handled, err := ErrorHandle(ctx, w, cc.err)
			resp := w.Result()
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, cc.want.handled, handled)
			want = response.ApplicationResponse{
				Success:         cc.want.resType == response.Success,
				Data:            data,
				Code:            cc.want.resType.Code,
				Message:         cc.want.resType.Message,
				ErrorAttributes: cc.want.errAttr,
			}
			wb, err := json.Marshal(want)
			assert.NoError(t, err)
			testutils.AssertResponse(t, resp, cc.want.resType.StatusCode, wb)
		})
	}
}

type notApplicationErr struct{}

func (e notApplicationErr) Error() string {
	return "notApplicationErr"
}

type wantErr struct{}

func (e wantErr) Error() string {
	return "wantErr"
}

func (e wantErr) ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes) {
	return response.ModelConflict, nil
}

type wantErrWithAttr struct {
	attr response.ApplicationErrorAttributes
}

func (e wantErrWithAttr) Error() string {
	return "wantErrWithAttr"
}

func (e wantErrWithAttr) ResolveCodeAndAttribute() (response.APIResponseType, response.ApplicationErrorAttributes) {
	return response.Validation, e.attr
}
