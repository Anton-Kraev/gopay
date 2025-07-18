package gopay_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Anton-Kraev/gopay"
	"github.com/Anton-Kraev/gopay/mocks"
)

type mockFields struct {
	mockLinks    *mocks.MocklinkGenerator
	mockStorage  *mocks.MockpaymentStorage
	mockPayments *mocks.MockpaymentService
}

func setupMocks(ctrl *gomock.Controller) (mockFields, *gopay.PaymentManager) {
	mf := mockFields{
		mockLinks:    mocks.NewMocklinkGenerator(ctrl),
		mockStorage:  mocks.NewMockpaymentStorage(ctrl),
		mockPayments: mocks.NewMockpaymentService(ctrl),
	}

	pm := gopay.NewPaymentManager(mf.mockLinks, mf.mockStorage, mf.mockPayments)

	return mf, pm
}

func TestPaymentManager_CreatePayment(t *testing.T) {
	t.Parallel()

	type (
		args struct {
			template gopay.PaymentTemplate
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
			name: "error generate id and link",
			setupMocks: func(f mockFields) {
				f.mockLinks.EXPECT().GenerateLink().
					Return(gopay.ID(""), gopay.Link(""), errors.New("error generate id and link")).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("error generate id and link"),
			},
		},
		{
			name: "error create payment",
			setupMocks: func(f mockFields) {
				f.mockLinks.EXPECT().GenerateLink().
					Return(gopay.ID("uuid"), gopay.Link("https://redirect.com/uuid"), nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("uuid"), gopay.PaymentTemplate{}).
					Return(nil, errors.New("error create payment")).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("error create payment"),
			},
		},
		{
			name: "error empty payment",
			setupMocks: func(f mockFields) {
				f.mockLinks.EXPECT().GenerateLink().
					Return(gopay.ID("uuid"), gopay.Link("https://redirect.com/uuid"), nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("uuid"), gopay.PaymentTemplate{}).
					Return(nil, nil).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("create payment failed"),
			},
		},
		{
			name: "error set payment",
			setupMocks: func(f mockFields) {
				f.mockLinks.EXPECT().GenerateLink().
					Return(gopay.ID("uuid"), gopay.Link("https://redirect.com/uuid"), nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("uuid"), gopay.PaymentTemplate{}).
					Return(&gopay.Payment{}, nil).Times(1)
				f.mockStorage.EXPECT().Set(gopay.ID("uuid"), gomock.Any()).
					Return(errors.New("error set payment")).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("error set payment"),
			},
		},
		{
			name: "error set link",
			setupMocks: func(f mockFields) {
				f.mockLinks.EXPECT().GenerateLink().
					Return(gopay.ID("uuid"), gopay.Link("https://redirect.com/uuid"), nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("uuid"), gopay.PaymentTemplate{}).
					Return(&gopay.Payment{PaymentLink: "payment"}, nil).Times(1)
				f.mockStorage.EXPECT().Set(gopay.ID("uuid"), gomock.Any()).
					Return(nil).Times(1)
				f.mockStorage.EXPECT().SetLink(gopay.ID("uuid"), gopay.Link("payment")).
					Return(errors.New("error set link")).Times(1)
			},
			expected: expected{
				link: gopay.Link(""),
				err:  errors.New("error set link"),
			},
		},
		{
			name: "success",
			args: args{
				template: gopay.PaymentTemplate{
					Currency:     "RUB",
					Amount:       100,
					Description:  "description",
					ResourceLink: "resource",
				},
				user: gopay.User{
					ID:    "1",
					Name:  "name",
					Email: "email",
				},
			},
			setupMocks: func(f mockFields) {
				f.mockLinks.EXPECT().GenerateLink().
					Return(gopay.ID("uuid"), gopay.Link("https://redirect.com/uuid"), nil).Times(1)
				f.mockPayments.EXPECT().CreatePayment(gopay.ID("uuid"), gopay.PaymentTemplate{
					Currency:     "RUB",
					Amount:       100,
					Description:  "description",
					ResourceLink: "resource",
				}).
					Return(&gopay.Payment{Amount: 100, Status: gopay.StatusPending, PaymentLink: "payment"}, nil).
					Times(1)
				f.mockStorage.EXPECT().Set(gopay.ID("uuid"), gomock.Any()).
					DoAndReturn(func(_ gopay.ID, payment gopay.Payment) error {
						if payment.Amount != 100 || payment.Status != gopay.StatusPending {
							return errors.New("error create new payment")
						}

						if payment.PaymentLink != "payment" || payment.ResourceLink != "resource" || payment.User.ID != "1" {
							return errors.New("error set payment fields")
						}

						return nil
					}).Times(1)
				f.mockStorage.EXPECT().SetLink(gopay.ID("uuid"), gopay.Link("payment")).
					Return(nil).Times(1)
			},
			expected: expected{
				link: gopay.Link("https://redirect.com/uuid"),
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

			link, err := pm.CreatePayment(tt.args.template, tt.args.user)

			if tt.expected.err != nil {
				require.EqualError(t, err, tt.expected.err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expected.link, link)
		})
	}
}

func TestPaymentManager_GetRedirectLink(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mf, pm := setupMocks(ctrl)

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mf.mockStorage.EXPECT().GetLink(gopay.ID("2")).
			Return(gopay.Link(""), errors.New("error get link")).Times(1)

		link, err := pm.GetRedirectLink("2")

		require.EqualError(t, err, "error get link")
		assert.Equal(t, gopay.Link(""), link)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mf.mockStorage.EXPECT().GetLink(gopay.ID("1")).
			Return(gopay.Link("redirect.link"), nil).Times(1)

		link, err := pm.GetRedirectLink("1")

		require.NoError(t, err)
		assert.Equal(t, gopay.Link("redirect.link"), link)
	})
}

func TestPaymentManager_UpdatePaymentStatus(t *testing.T) {
	t.Parallel()

	type args struct {
		id     gopay.ID
		status gopay.Status
	}

	tests := []struct {
		name        string
		args        args
		setupMocks  func(f mockFields)
		errExpected bool
	}{
		{
			name: "error get payment",
			args: args{
				id:     gopay.ID("1"),
				status: gopay.StatusSucceeded,
			},
			setupMocks: func(f mockFields) {
				f.mockStorage.EXPECT().Get(gopay.ID("1")).
					Return(gopay.Payment{}, errors.New("error get payment")).Times(1)
			},
			errExpected: true,
		},
		{
			name: "error set link",
			args: args{
				id:     gopay.ID("1"),
				status: gopay.StatusSucceeded,
			},
			setupMocks: func(f mockFields) {
				f.mockStorage.EXPECT().Get(gopay.ID("1")).
					Return(gopay.Payment{ResourceLink: "resource.link"}, nil).Times(1)
				f.mockStorage.EXPECT().SetLink(gopay.ID("1"), gopay.Link("resource.link")).
					Return(errors.New("error set link")).Times(1)
			},
			errExpected: true,
		},
		{
			name: "error update status",
			args: args{
				id:     gopay.ID("1"),
				status: gopay.StatusWaitingForCapture,
			},
			setupMocks: func(f mockFields) {
				f.mockStorage.EXPECT().UpdateStatus(gopay.ID("1"), gopay.StatusWaitingForCapture).
					Return(errors.New("error update status")).Times(1)
			},
			errExpected: true,
		},
		{
			name: "success",
			args: args{
				id:     gopay.ID("1"),
				status: gopay.StatusSucceeded,
			},
			setupMocks: func(f mockFields) {
				f.mockStorage.EXPECT().Get(gopay.ID("1")).
					Return(gopay.Payment{ResourceLink: "resource.link"}, nil).Times(1)
				f.mockStorage.EXPECT().SetLink(gopay.ID("1"), gopay.Link("resource.link")).
					Return(nil).Times(1)
				f.mockStorage.EXPECT().UpdateStatus(gopay.ID("1"), gopay.StatusSucceeded).
					Return(nil).Times(1)
			},
			errExpected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mf, pm := setupMocks(ctrl)
			tt.setupMocks(mf)

			err := pm.UpdatePaymentStatus(tt.args.id, tt.args.status)

			require.Equal(t, tt.errExpected, err != nil)
		})
	}
}
