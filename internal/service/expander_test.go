package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/spanwalla/url-shortener/internal/entity"
	repomocks "github.com/spanwalla/url-shortener/internal/mocks/repository"
	"github.com/spanwalla/url-shortener/internal/repository/repoerrs"
)

func TestExpanderService_Expand(t *testing.T) {
	type args struct {
		ctx   context.Context
		alias string
	}

	type MockBehavior func(l *repomocks.MockLink, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         string
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				alias: "46g1B3ZgAy",
			},
			mockBehavior: func(l *repomocks.MockLink, args args) {
				l.EXPECT().Get(args.ctx, args.alias).
					Return(entity.Link{Alias: "46g1B3ZgAy", URI: "https://github.com/spanwalla"}, nil)
			},
			want:    "https://github.com/spanwalla",
			wantErr: false,
		},
		{
			name: "uri not found",
			args: args{
				ctx:   context.Background(),
				alias: "46g1B3ZgAf",
			},
			mockBehavior: func(l *repomocks.MockLink, args args) {
				l.EXPECT().Get(args.ctx, args.alias).
					Return(entity.Link{}, repoerrs.ErrNotFound)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "unexpected error",
			args: args{
				ctx:   context.Background(),
				alias: "dm749Kmd",
			},
			mockBehavior: func(l *repomocks.MockLink, args args) {
				l.EXPECT().Get(args.ctx, args.alias).
					Return(entity.Link{}, errors.New("error"))
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLink := repomocks.NewMockLink(ctrl)
			tc.mockBehavior(mockLink, tc.args)

			s := NewExpanderService(mockLink)

			got, err := s.Expand(tc.args.ctx, tc.args.alias)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.want, got)
			assert.NoError(t, err)
		})
	}
}
