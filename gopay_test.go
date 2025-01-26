package gopay_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Anton-Kraev/gopay"
	"github.com/Anton-Kraev/gopay/mocks"
)

type mockFields struct {
	mockTemplates *mocks.Mocktemplates
	mockLinks     *mocks.MocklinkGenerator
	mockStorage   *mocks.Mockstorage
	mockPayments  *mocks.MockpaymentService
}

func setupMocks(ctrl *gomock.Controller) (mockFields, *gopay.PaymentManager) {
	mf := mockFields{
		mockTemplates: mocks.NewMocktemplates(ctrl),
		mockLinks:     mocks.NewMocklinkGenerator(ctrl),
		mockStorage:   mocks.NewMockstorage(ctrl),
		mockPayments:  mocks.NewMockpaymentService(ctrl),
	}

	pm := gopay.NewPaymentManager(mf.mockTemplates, mf.mockLinks, mf.mockStorage, mf.mockPayments)

	return mf, pm
}

func TestPaymentManager_NewPayment(t *testing.T) {
	t.Parallel()

	type (
		args struct {
			template string
			user     gopay.User
		}

		expected struct {
			link gopay.Link
			err  error
		}
	)

	tests := []struct {
		name       string
		args       args
		setupMocks func(f mockFields)
		expected   expected
	}{
		{
			name: "error getting template",
			args: args{
				template: "payment_template",
			},
			setupMocks: func(f mockFields) {
				f.mockTemplates.EXPECT().GetTemplate("payment_template").
					Return(gopay.PaymentTemplate{}, errors.New("template not found")).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("template not found"),
			},
		},
		{
			name: "error create payment",
			args: args{
				template: "payment_template",
				user:     gopay.User{ID: gopay.ID("1")},
			},
			setupMocks: func(f mockFields) {
				f.mockTemplates.EXPECT().GetTemplate("payment_template").
					Return(gopay.PaymentTemplate{}, nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("1"), gopay.PaymentTemplate{}).
					Return(nil, errors.New("error create payment")).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("error create payment"),
			},
		},
		{
			name: "error empty payment",
			args: args{
				template: "payment_template",
				user:     gopay.User{ID: gopay.ID("1")},
			},
			setupMocks: func(f mockFields) {
				f.mockTemplates.EXPECT().GetTemplate("payment_template").
					Return(gopay.PaymentTemplate{}, nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("1"), gopay.PaymentTemplate{}).
					Return(nil, nil).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("create payment failed"),
			},
		},
		{
			name: "error set user",
			args: args{
				template: "payment_template",
				user:     gopay.User{ID: gopay.ID("1")},
			},
			setupMocks: func(f mockFields) {
				f.mockTemplates.EXPECT().GetTemplate("payment_template").
					Return(gopay.PaymentTemplate{}, nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("1"), gopay.PaymentTemplate{}).
					Return(&gopay.Payment{User: gopay.User{ID: "1"}}, nil).Times(1)
				f.mockStorage.EXPECT().SetUser(gopay.ID("1"), gomock.Any()).
					Return(errors.New("error set user")).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("error set user"),
			},
		},
		{
			name: "error generate link",
			args: args{
				template: "payment_template",
				user:     gopay.User{ID: gopay.ID("1")},
			},
			setupMocks: func(f mockFields) {
				f.mockTemplates.EXPECT().GetTemplate("payment_template").
					Return(gopay.PaymentTemplate{}, nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("1"), gopay.PaymentTemplate{}).
					Return(&gopay.Payment{User: gopay.User{ID: "1"}}, nil).Times(1)
				f.mockStorage.EXPECT().SetUser(gopay.ID("1"), gomock.Any()).
					Return(nil).Times(1)
				f.mockLinks.EXPECT().GenerateLink(gopay.ID("1")).
					Return(gopay.Link(""), errors.New("error generate link")).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("error generate link"),
			},
		},
		{
			name: "success new payment",
			args: args{
				template: "payment_template",
				user: gopay.User{
					ID:    "1",
					Name:  "name",
					Email: "email",
				},
			},
			setupMocks: func(f mockFields) {
				f.mockTemplates.EXPECT().GetTemplate("payment_template").
					Return(gopay.PaymentTemplate{}, nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("1"), gopay.PaymentTemplate{}).
					Return(&gopay.Payment{User: gopay.User{ID: "1"}}, nil).Times(1)
				f.mockStorage.EXPECT().SetUser(gopay.ID("1"), gomock.Any()).
					DoAndReturn(func(id gopay.ID, user gopay.User) error {
						if user.ID != id {
							return errors.New("wrong user id")
						}

						return nil
					}).Times(1)
				f.mockLinks.EXPECT().GenerateLink(gopay.ID("1")).
					Return(gopay.Link("https://redirect.com"), nil).Times(1)
			},
			expected: expected{
				link: gopay.Link("https://redirect.com"),
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mf, pm := setupMocks(ctrl)
			tt.setupMocks(mf)

			link, err := pm.NewPayment(tt.args.template, tt.args.user)

			if tt.expected.err != nil {
				require.EqualError(t, err, tt.expected.err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expected.link, link)
		})
	}
}

func TestPaymentManager_FollowLink(t *testing.T) {
	t.Parallel()

	type (
		args struct {
			id gopay.ID
		}

		expected struct {
			status gopay.Status
			link   gopay.Link
			err    error
		}
	)

	tests := []struct {
		name       string
		args       args
		setupMocks func(f mockFields)
		expected   expected
	}{
		{
			name: "error getting links",
			args: args{
				id: gopay.ID("1"),
			},
			setupMocks: func(f mockFields) {
				f.mockStorage.EXPECT().GetLinks(gopay.ID("1")).
					Return(gopay.Links{}, errors.New("links not found")).Times(1)
			},
			expected: expected{
				status: gopay.Status(""),
				link:   gopay.Link(""),
				err:    errors.New("links not found"),
			},
		},
		{
			name: "success resource redirect",
			args: args{
				id: gopay.ID("1"),
			},
			setupMocks: func(f mockFields) {
				f.mockStorage.EXPECT().GetLinks(gopay.ID("1")).
					Return(gopay.Links{
						Status:       gopay.StatusSucceeded,
						ResourceLink: "https://resource.com",
					}, nil).Times(1)
			},
			expected: expected{
				status: gopay.StatusSucceeded,
				link:   gopay.Link("https://resource.com"),
			},
		},
		{
			name: "success payment redirect",
			args: args{
				id: gopay.ID("1"),
			},
			setupMocks: func(f mockFields) {
				f.mockStorage.EXPECT().GetLinks(gopay.ID("1")).
					Return(gopay.Links{
						Status:      gopay.StatusWaitingForCapture,
						PaymentLink: "https://payment.com",
					}, nil).Times(1)
			},
			expected: expected{
				status: gopay.StatusWaitingForCapture,
				link:   gopay.Link("https://payment.com"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mf, pm := setupMocks(ctrl)
			tt.setupMocks(mf)

			status, link, err := pm.FollowLink(tt.args.id)

			if tt.expected.err != nil {
				require.EqualError(t, err, tt.expected.err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expected.status, status)
			assert.Equal(t, tt.expected.link, link)
		})
	}
}

func TestPaymentManager_Checkout(t *testing.T) {
	t.Parallel()

	t.Run("smoke test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mf, pm := setupMocks(ctrl)
		mf.mockStorage.EXPECT().UpdateStatus(gopay.ID("1"), gopay.StatusSucceeded).
			Return(nil).Times(1)

		err := pm.Checkout("1", gopay.StatusSucceeded)

		require.NoError(t, err)
	})
}
