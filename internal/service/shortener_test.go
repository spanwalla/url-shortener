package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	encodermocks "github.com/spanwalla/url-shortener/internal/mocks/encoder"
	repomocks "github.com/spanwalla/url-shortener/internal/mocks/repository"
	"github.com/spanwalla/url-shortener/internal/repository/repoerrs"
)

const testAliasLength = 10
const testAttemptsOnCollision = 3

func TestShortenerService_Shorten(t *testing.T) {
	type args struct {
		ctx        context.Context
		alias, uri string
	}

	type MockBehavior func(l *repomocks.MockLink, e *encodermocks.MockEncoder, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantAlias    string
		wantCreated  bool
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				alias: "test123Uri",
				uri:   "https://ya.ru/?npr=1",
			},
			mockBehavior: func(l *repomocks.MockLink, e *encodermocks.MockEncoder, args args) {
				e.EXPECT().Encode(args.uri, testAliasLength).
					Return(args.alias)
				l.EXPECT().Store(args.ctx, args.alias, args.uri).
					Return(args.alias, nil)
			},
			wantAlias:   "test123Uri",
			wantCreated: true,
			wantErr:     false,
		},
		{
			name: "uri already exists",
			args: args{
				ctx:   context.Background(),
				alias: "test123Uri",
				uri:   "https://ya.ru/?npr=1",
			},
			mockBehavior: func(l *repomocks.MockLink, e *encodermocks.MockEncoder, args args) {
				e.EXPECT().Encode(args.uri, testAliasLength).AnyTimes().
					Return(args.alias)
				l.EXPECT().Store(args.ctx, args.alias, args.uri).AnyTimes().
					Return("anotherAli", nil)
			},
			wantAlias:   "anotherAli",
			wantCreated: false,
			wantErr:     false,
		},
		{
			name: "unavoidable collision",
			args: args{
				ctx:   context.Background(),
				alias: "test123Uri",
				uri:   "https://ya.ru/?npr=1",
			},
			mockBehavior: func(l *repomocks.MockLink, e *encodermocks.MockEncoder, args args) {
				e.EXPECT().Encode(args.uri, testAliasLength).AnyTimes().
					Return(args.alias)
				l.EXPECT().Store(args.ctx, args.alias, args.uri).AnyTimes().
					Return("", repoerrs.ErrAlreadyExists)
			},
			wantAlias:   "",
			wantCreated: false,
			wantErr:     true,
		},
		{
			name: "unknown error",
			args: args{
				ctx:   context.Background(),
				alias: "test123Uri",
				uri:   "https://ya.ru/?npr=1",
			},
			mockBehavior: func(l *repomocks.MockLink, e *encodermocks.MockEncoder, args args) {
				e.EXPECT().Encode(args.uri, testAliasLength).
					Return(args.alias)
				l.EXPECT().Store(args.ctx, args.alias, args.uri).
					Return("", errors.New("unknown error"))
			},
			wantAlias:   "",
			wantCreated: false,
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLink := repomocks.NewMockLink(ctrl)
			mockEncoder := encodermocks.NewMockEncoder(ctrl)

			tc.mockBehavior(mockLink, mockEncoder, tc.args)

			s := NewShortenerService(mockLink, mockEncoder, testAliasLength, testAttemptsOnCollision)

			alias, created, err := s.Shorten(tc.args.ctx, tc.args.uri)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tc.wantAlias, alias)
			assert.Equal(t, tc.wantCreated, created)
			assert.NoError(t, err)
		})
	}
}
