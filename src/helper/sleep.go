package helper

import (
	"time"

	"github.com/vault-thirteen/auxie/random"
)

const (
	IntervalBeforeReturningFakeRequestIdSecMin       = 5
	IntervalBeforeReturningFakeRequestIdSecMax       = 15
	IntervalBeforeReportingFailedVerificationCodeMin = 5
	IntervalBeforeReportingFailedVerificationCodeMax = 15
	IntervalFraudResponseMin                         = 30
	IntervalFraudResponseMax                         = 60
)

func SleepBeforeReturningFakeRequestId() (err error) {
	var intervalSec uint
	intervalSec, err = random.Uint(IntervalBeforeReturningFakeRequestIdSecMin, IntervalBeforeReturningFakeRequestIdSecMax)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * time.Duration(intervalSec))
	return nil
}
func SleepBeforeReportingFailedVerificationCode() (err error) {
	var intervalSec uint
	intervalSec, err = random.Uint(IntervalBeforeReportingFailedVerificationCodeMin, IntervalBeforeReportingFailedVerificationCodeMax)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * time.Duration(intervalSec))
	return nil
}
func SleepBeforeFraudResponse() (err error) {
	var intervalSec uint
	intervalSec, err = random.Uint(IntervalFraudResponseMin, IntervalFraudResponseMax)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * time.Duration(intervalSec))
	return nil
}
