package sms

import (
	"errors"
	"go.uber.org/mock/gomock"
	"math/rand"
	"microservice/pkg/mock"
)

func New() ISmsProvider {
	return initMock()
}

//

func initMock() *MockISmsProvider {
	ctrl := gomock.NewController(mock.NewController())
	defer ctrl.Finish()

	mocked := NewMockISmsProvider(ctrl)
	mocked.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(func(phone string, message string) (map[string]interface{}, error) {
		success := map[string]interface{}{"status": 200, "result": "sent"}
		failed := map[string]interface{}{"status": 400, "result": "failed"}

		if rand.Intn(10) == 1 {
			return failed, errors.New("bad request. try again")
		}

		return success, nil
	}).AnyTimes()

	return mocked
}
