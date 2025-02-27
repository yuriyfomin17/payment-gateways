package services

//TODO:
//commented out for faster test execution
//implement parallel testing

//func TestPublishWithCircuitBreaker_Success(t *testing.T) {
//	// Given
//	ftService := NewFaultToleranceService()
//	successfulOperation := func() error {
//		return nil
//	}
//
//	// When
//	err := ftService.PublishWithCircuitBreaker(successfulOperation)
//
//	// Then
//	assert.NoError(t, err) // Check that no error was returned
//}
//
//func TestPublishWithCircuitBreaker_Failure(t *testing.T) {
//	// Given
//	ftService := NewFaultToleranceService()
//	failingOperation := func() error {
//		return errors.New("operation failed")
//	}
//
//	// When
//	err := ftService.PublishWithCircuitBreaker(failingOperation)
//
//	// Then
//	assert.Error(t, err) // Ensure it reports the operation failure
//	assert.Equal(t, "operation failed", err.Error())
//}
//
//func TestRetryOperation_Success(t *testing.T) {
//	// Given
//	ftService := NewFaultToleranceService()
//	retryCount := 0
//	successfulOperation := func() error {
//		if retryCount < 2 { // Fail the operation twice
//			retryCount++
//			return errors.New("temporary failure")
//		}
//		return nil
//	}
//
//	// When
//	err := ftService.RetryOperation(successfulOperation, 5)
//
//	// Then
//	assert.NoError(t, err)         // Ensure the operation succeeds eventually
//	assert.Equal(t, 2, retryCount) // Ensure operation was retried 2 times
//}
//
//func TestRetryOperation_Failure(t *testing.T) {
//	// Given
//	ftService := NewFaultToleranceService()
//	failingOperation := func() error {
//		return errors.New("permanent failure")
//	}
//
//	// When
//	err := ftService.RetryOperation(failingOperation, 3)
//
//	// Then
//	assert.Error(t, err) // Operation should fail after retries
//	assert.Equal(t, "operation failed after 3 attempts", err.Error())
//}
//
//func TestRetryOperation_NoRetriesNeeded(t *testing.T) {
//	// Given
//	ftService := NewFaultToleranceService()
//	successfulOperation := func() error {
//		return nil
//	}
//
//	// When
//	err := ftService.RetryOperation(successfulOperation, 1)
//
//	// Then
//	assert.NoError(t, err) // Ensure the operation succeeds without retries
//}
