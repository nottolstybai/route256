// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package service

//go:generate minimock -i route256.ozon.ru/project/cart/internal/service.Repository -o repository_mock_test.go -n RepositoryMock -p service

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/cart/internal/entity"
)

// RepositoryMock implements Repository
type RepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcInsert          func(ctx context.Context, userID int64, skuID int64, count uint16) (err error)
	inspectFuncInsert   func(ctx context.Context, userID int64, skuID int64, count uint16)
	afterInsertCounter  uint64
	beforeInsertCounter uint64
	InsertMock          mRepositoryMockInsert

	funcList          func(ctx context.Context, userID int64) (pa1 []entity.ProductInfo, err error)
	inspectFuncList   func(ctx context.Context, userID int64)
	afterListCounter  uint64
	beforeListCounter uint64
	ListMock          mRepositoryMockList

	funcRemove          func(ctx context.Context, userID int64, skuID int64) (err error)
	inspectFuncRemove   func(ctx context.Context, userID int64, skuID int64)
	afterRemoveCounter  uint64
	beforeRemoveCounter uint64
	RemoveMock          mRepositoryMockRemove

	funcRemoveByUserID          func(ctx context.Context, userID int64) (err error)
	inspectFuncRemoveByUserID   func(ctx context.Context, userID int64)
	afterRemoveByUserIDCounter  uint64
	beforeRemoveByUserIDCounter uint64
	RemoveByUserIDMock          mRepositoryMockRemoveByUserID
}

// NewRepositoryMock returns a mock for Repository
func NewRepositoryMock(t minimock.Tester) *RepositoryMock {
	m := &RepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.InsertMock = mRepositoryMockInsert{mock: m}
	m.InsertMock.callArgs = []*RepositoryMockInsertParams{}

	m.ListMock = mRepositoryMockList{mock: m}
	m.ListMock.callArgs = []*RepositoryMockListParams{}

	m.RemoveMock = mRepositoryMockRemove{mock: m}
	m.RemoveMock.callArgs = []*RepositoryMockRemoveParams{}

	m.RemoveByUserIDMock = mRepositoryMockRemoveByUserID{mock: m}
	m.RemoveByUserIDMock.callArgs = []*RepositoryMockRemoveByUserIDParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mRepositoryMockInsert struct {
	mock               *RepositoryMock
	defaultExpectation *RepositoryMockInsertExpectation
	expectations       []*RepositoryMockInsertExpectation

	callArgs []*RepositoryMockInsertParams
	mutex    sync.RWMutex
}

// RepositoryMockInsertExpectation specifies expectation struct of the Repository.Insert
type RepositoryMockInsertExpectation struct {
	mock    *RepositoryMock
	params  *RepositoryMockInsertParams
	results *RepositoryMockInsertResults
	Counter uint64
}

// RepositoryMockInsertParams contains parameters of the Repository.Insert
type RepositoryMockInsertParams struct {
	ctx    context.Context
	userID int64
	skuID  int64
	count  uint16
}

// RepositoryMockInsertResults contains results of the Repository.Insert
type RepositoryMockInsertResults struct {
	err error
}

// Expect sets up expected params for Repository.Insert
func (mmInsert *mRepositoryMockInsert) Expect(ctx context.Context, userID int64, skuID int64, count uint16) *mRepositoryMockInsert {
	if mmInsert.mock.funcInsert != nil {
		mmInsert.mock.t.Fatalf("RepositoryMock.Insert mock is already set by Set")
	}

	if mmInsert.defaultExpectation == nil {
		mmInsert.defaultExpectation = &RepositoryMockInsertExpectation{}
	}

	mmInsert.defaultExpectation.params = &RepositoryMockInsertParams{ctx, userID, skuID, count}
	for _, e := range mmInsert.expectations {
		if minimock.Equal(e.params, mmInsert.defaultExpectation.params) {
			mmInsert.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmInsert.defaultExpectation.params)
		}
	}

	return mmInsert
}

// Inspect accepts an inspector function that has same arguments as the Repository.Insert
func (mmInsert *mRepositoryMockInsert) Inspect(f func(ctx context.Context, userID int64, skuID int64, count uint16)) *mRepositoryMockInsert {
	if mmInsert.mock.inspectFuncInsert != nil {
		mmInsert.mock.t.Fatalf("Inspect function is already set for RepositoryMock.Insert")
	}

	mmInsert.mock.inspectFuncInsert = f

	return mmInsert
}

// Return sets up results that will be returned by Repository.Insert
func (mmInsert *mRepositoryMockInsert) Return(err error) *RepositoryMock {
	if mmInsert.mock.funcInsert != nil {
		mmInsert.mock.t.Fatalf("RepositoryMock.Insert mock is already set by Set")
	}

	if mmInsert.defaultExpectation == nil {
		mmInsert.defaultExpectation = &RepositoryMockInsertExpectation{mock: mmInsert.mock}
	}
	mmInsert.defaultExpectation.results = &RepositoryMockInsertResults{err}
	return mmInsert.mock
}

// Set uses given function f to mock the Repository.Insert method
func (mmInsert *mRepositoryMockInsert) Set(f func(ctx context.Context, userID int64, skuID int64, count uint16) (err error)) *RepositoryMock {
	if mmInsert.defaultExpectation != nil {
		mmInsert.mock.t.Fatalf("Default expectation is already set for the Repository.Insert method")
	}

	if len(mmInsert.expectations) > 0 {
		mmInsert.mock.t.Fatalf("Some expectations are already set for the Repository.Insert method")
	}

	mmInsert.mock.funcInsert = f
	return mmInsert.mock
}

// When sets expectation for the Repository.Insert which will trigger the result defined by the following
// Then helper
func (mmInsert *mRepositoryMockInsert) When(ctx context.Context, userID int64, skuID int64, count uint16) *RepositoryMockInsertExpectation {
	if mmInsert.mock.funcInsert != nil {
		mmInsert.mock.t.Fatalf("RepositoryMock.Insert mock is already set by Set")
	}

	expectation := &RepositoryMockInsertExpectation{
		mock:   mmInsert.mock,
		params: &RepositoryMockInsertParams{ctx, userID, skuID, count},
	}
	mmInsert.expectations = append(mmInsert.expectations, expectation)
	return expectation
}

// Then sets up Repository.Insert return parameters for the expectation previously defined by the When method
func (e *RepositoryMockInsertExpectation) Then(err error) *RepositoryMock {
	e.results = &RepositoryMockInsertResults{err}
	return e.mock
}

// Insert implements Repository
func (mmInsert *RepositoryMock) Insert(ctx context.Context, userID int64, skuID int64, count uint16) (err error) {
	mm_atomic.AddUint64(&mmInsert.beforeInsertCounter, 1)
	defer mm_atomic.AddUint64(&mmInsert.afterInsertCounter, 1)

	if mmInsert.inspectFuncInsert != nil {
		mmInsert.inspectFuncInsert(ctx, userID, skuID, count)
	}

	mm_params := RepositoryMockInsertParams{ctx, userID, skuID, count}

	// Record call args
	mmInsert.InsertMock.mutex.Lock()
	mmInsert.InsertMock.callArgs = append(mmInsert.InsertMock.callArgs, &mm_params)
	mmInsert.InsertMock.mutex.Unlock()

	for _, e := range mmInsert.InsertMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmInsert.InsertMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmInsert.InsertMock.defaultExpectation.Counter, 1)
		mm_want := mmInsert.InsertMock.defaultExpectation.params
		mm_got := RepositoryMockInsertParams{ctx, userID, skuID, count}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmInsert.t.Errorf("RepositoryMock.Insert got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmInsert.InsertMock.defaultExpectation.results
		if mm_results == nil {
			mmInsert.t.Fatal("No results are set for the RepositoryMock.Insert")
		}
		return (*mm_results).err
	}
	if mmInsert.funcInsert != nil {
		return mmInsert.funcInsert(ctx, userID, skuID, count)
	}
	mmInsert.t.Fatalf("Unexpected call to RepositoryMock.Insert. %v %v %v %v", ctx, userID, skuID, count)
	return
}

// InsertAfterCounter returns a count of finished RepositoryMock.Insert invocations
func (mmInsert *RepositoryMock) InsertAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmInsert.afterInsertCounter)
}

// InsertBeforeCounter returns a count of RepositoryMock.Insert invocations
func (mmInsert *RepositoryMock) InsertBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmInsert.beforeInsertCounter)
}

// Calls returns a list of arguments used in each call to RepositoryMock.Insert.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmInsert *mRepositoryMockInsert) Calls() []*RepositoryMockInsertParams {
	mmInsert.mutex.RLock()

	argCopy := make([]*RepositoryMockInsertParams, len(mmInsert.callArgs))
	copy(argCopy, mmInsert.callArgs)

	mmInsert.mutex.RUnlock()

	return argCopy
}

// MinimockInsertDone returns true if the count of the Insert invocations corresponds
// the number of defined expectations
func (m *RepositoryMock) MinimockInsertDone() bool {
	for _, e := range m.InsertMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.InsertMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterInsertCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcInsert != nil && mm_atomic.LoadUint64(&m.afterInsertCounter) < 1 {
		return false
	}
	return true
}

// MinimockInsertInspect logs each unmet expectation
func (m *RepositoryMock) MinimockInsertInspect() {
	for _, e := range m.InsertMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RepositoryMock.Insert with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.InsertMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterInsertCounter) < 1 {
		if m.InsertMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RepositoryMock.Insert")
		} else {
			m.t.Errorf("Expected call to RepositoryMock.Insert with params: %#v", *m.InsertMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcInsert != nil && mm_atomic.LoadUint64(&m.afterInsertCounter) < 1 {
		m.t.Error("Expected call to RepositoryMock.Insert")
	}
}

type mRepositoryMockList struct {
	mock               *RepositoryMock
	defaultExpectation *RepositoryMockListExpectation
	expectations       []*RepositoryMockListExpectation

	callArgs []*RepositoryMockListParams
	mutex    sync.RWMutex
}

// RepositoryMockListExpectation specifies expectation struct of the Repository.List
type RepositoryMockListExpectation struct {
	mock    *RepositoryMock
	params  *RepositoryMockListParams
	results *RepositoryMockListResults
	Counter uint64
}

// RepositoryMockListParams contains parameters of the Repository.List
type RepositoryMockListParams struct {
	ctx    context.Context
	userID int64
}

// RepositoryMockListResults contains results of the Repository.List
type RepositoryMockListResults struct {
	pa1 []entity.ProductInfo
	err error
}

// Expect sets up expected params for Repository.List
func (mmList *mRepositoryMockList) Expect(ctx context.Context, userID int64) *mRepositoryMockList {
	if mmList.mock.funcList != nil {
		mmList.mock.t.Fatalf("RepositoryMock.List mock is already set by Set")
	}

	if mmList.defaultExpectation == nil {
		mmList.defaultExpectation = &RepositoryMockListExpectation{}
	}

	mmList.defaultExpectation.params = &RepositoryMockListParams{ctx, userID}
	for _, e := range mmList.expectations {
		if minimock.Equal(e.params, mmList.defaultExpectation.params) {
			mmList.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmList.defaultExpectation.params)
		}
	}

	return mmList
}

// Inspect accepts an inspector function that has same arguments as the Repository.List
func (mmList *mRepositoryMockList) Inspect(f func(ctx context.Context, userID int64)) *mRepositoryMockList {
	if mmList.mock.inspectFuncList != nil {
		mmList.mock.t.Fatalf("Inspect function is already set for RepositoryMock.List")
	}

	mmList.mock.inspectFuncList = f

	return mmList
}

// Return sets up results that will be returned by Repository.List
func (mmList *mRepositoryMockList) Return(pa1 []entity.ProductInfo, err error) *RepositoryMock {
	if mmList.mock.funcList != nil {
		mmList.mock.t.Fatalf("RepositoryMock.List mock is already set by Set")
	}

	if mmList.defaultExpectation == nil {
		mmList.defaultExpectation = &RepositoryMockListExpectation{mock: mmList.mock}
	}
	mmList.defaultExpectation.results = &RepositoryMockListResults{pa1, err}
	return mmList.mock
}

// Set uses given function f to mock the Repository.List method
func (mmList *mRepositoryMockList) Set(f func(ctx context.Context, userID int64) (pa1 []entity.ProductInfo, err error)) *RepositoryMock {
	if mmList.defaultExpectation != nil {
		mmList.mock.t.Fatalf("Default expectation is already set for the Repository.List method")
	}

	if len(mmList.expectations) > 0 {
		mmList.mock.t.Fatalf("Some expectations are already set for the Repository.List method")
	}

	mmList.mock.funcList = f
	return mmList.mock
}

// When sets expectation for the Repository.List which will trigger the result defined by the following
// Then helper
func (mmList *mRepositoryMockList) When(ctx context.Context, userID int64) *RepositoryMockListExpectation {
	if mmList.mock.funcList != nil {
		mmList.mock.t.Fatalf("RepositoryMock.List mock is already set by Set")
	}

	expectation := &RepositoryMockListExpectation{
		mock:   mmList.mock,
		params: &RepositoryMockListParams{ctx, userID},
	}
	mmList.expectations = append(mmList.expectations, expectation)
	return expectation
}

// Then sets up Repository.List return parameters for the expectation previously defined by the When method
func (e *RepositoryMockListExpectation) Then(pa1 []entity.ProductInfo, err error) *RepositoryMock {
	e.results = &RepositoryMockListResults{pa1, err}
	return e.mock
}

// List implements Repository
func (mmList *RepositoryMock) List(ctx context.Context, userID int64) (pa1 []entity.ProductInfo, err error) {
	mm_atomic.AddUint64(&mmList.beforeListCounter, 1)
	defer mm_atomic.AddUint64(&mmList.afterListCounter, 1)

	if mmList.inspectFuncList != nil {
		mmList.inspectFuncList(ctx, userID)
	}

	mm_params := RepositoryMockListParams{ctx, userID}

	// Record call args
	mmList.ListMock.mutex.Lock()
	mmList.ListMock.callArgs = append(mmList.ListMock.callArgs, &mm_params)
	mmList.ListMock.mutex.Unlock()

	for _, e := range mmList.ListMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.pa1, e.results.err
		}
	}

	if mmList.ListMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmList.ListMock.defaultExpectation.Counter, 1)
		mm_want := mmList.ListMock.defaultExpectation.params
		mm_got := RepositoryMockListParams{ctx, userID}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmList.t.Errorf("RepositoryMock.List got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmList.ListMock.defaultExpectation.results
		if mm_results == nil {
			mmList.t.Fatal("No results are set for the RepositoryMock.List")
		}
		return (*mm_results).pa1, (*mm_results).err
	}
	if mmList.funcList != nil {
		return mmList.funcList(ctx, userID)
	}
	mmList.t.Fatalf("Unexpected call to RepositoryMock.List. %v %v", ctx, userID)
	return
}

// ListAfterCounter returns a count of finished RepositoryMock.List invocations
func (mmList *RepositoryMock) ListAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmList.afterListCounter)
}

// ListBeforeCounter returns a count of RepositoryMock.List invocations
func (mmList *RepositoryMock) ListBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmList.beforeListCounter)
}

// Calls returns a list of arguments used in each call to RepositoryMock.List.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmList *mRepositoryMockList) Calls() []*RepositoryMockListParams {
	mmList.mutex.RLock()

	argCopy := make([]*RepositoryMockListParams, len(mmList.callArgs))
	copy(argCopy, mmList.callArgs)

	mmList.mutex.RUnlock()

	return argCopy
}

// MinimockListDone returns true if the count of the List invocations corresponds
// the number of defined expectations
func (m *RepositoryMock) MinimockListDone() bool {
	for _, e := range m.ListMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ListMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterListCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcList != nil && mm_atomic.LoadUint64(&m.afterListCounter) < 1 {
		return false
	}
	return true
}

// MinimockListInspect logs each unmet expectation
func (m *RepositoryMock) MinimockListInspect() {
	for _, e := range m.ListMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RepositoryMock.List with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ListMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterListCounter) < 1 {
		if m.ListMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RepositoryMock.List")
		} else {
			m.t.Errorf("Expected call to RepositoryMock.List with params: %#v", *m.ListMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcList != nil && mm_atomic.LoadUint64(&m.afterListCounter) < 1 {
		m.t.Error("Expected call to RepositoryMock.List")
	}
}

type mRepositoryMockRemove struct {
	mock               *RepositoryMock
	defaultExpectation *RepositoryMockRemoveExpectation
	expectations       []*RepositoryMockRemoveExpectation

	callArgs []*RepositoryMockRemoveParams
	mutex    sync.RWMutex
}

// RepositoryMockRemoveExpectation specifies expectation struct of the Repository.Remove
type RepositoryMockRemoveExpectation struct {
	mock    *RepositoryMock
	params  *RepositoryMockRemoveParams
	results *RepositoryMockRemoveResults
	Counter uint64
}

// RepositoryMockRemoveParams contains parameters of the Repository.Remove
type RepositoryMockRemoveParams struct {
	ctx    context.Context
	userID int64
	skuID  int64
}

// RepositoryMockRemoveResults contains results of the Repository.Remove
type RepositoryMockRemoveResults struct {
	err error
}

// Expect sets up expected params for Repository.Remove
func (mmRemove *mRepositoryMockRemove) Expect(ctx context.Context, userID int64, skuID int64) *mRepositoryMockRemove {
	if mmRemove.mock.funcRemove != nil {
		mmRemove.mock.t.Fatalf("RepositoryMock.Remove mock is already set by Set")
	}

	if mmRemove.defaultExpectation == nil {
		mmRemove.defaultExpectation = &RepositoryMockRemoveExpectation{}
	}

	mmRemove.defaultExpectation.params = &RepositoryMockRemoveParams{ctx, userID, skuID}
	for _, e := range mmRemove.expectations {
		if minimock.Equal(e.params, mmRemove.defaultExpectation.params) {
			mmRemove.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRemove.defaultExpectation.params)
		}
	}

	return mmRemove
}

// Inspect accepts an inspector function that has same arguments as the Repository.Remove
func (mmRemove *mRepositoryMockRemove) Inspect(f func(ctx context.Context, userID int64, skuID int64)) *mRepositoryMockRemove {
	if mmRemove.mock.inspectFuncRemove != nil {
		mmRemove.mock.t.Fatalf("Inspect function is already set for RepositoryMock.Remove")
	}

	mmRemove.mock.inspectFuncRemove = f

	return mmRemove
}

// Return sets up results that will be returned by Repository.Remove
func (mmRemove *mRepositoryMockRemove) Return(err error) *RepositoryMock {
	if mmRemove.mock.funcRemove != nil {
		mmRemove.mock.t.Fatalf("RepositoryMock.Remove mock is already set by Set")
	}

	if mmRemove.defaultExpectation == nil {
		mmRemove.defaultExpectation = &RepositoryMockRemoveExpectation{mock: mmRemove.mock}
	}
	mmRemove.defaultExpectation.results = &RepositoryMockRemoveResults{err}
	return mmRemove.mock
}

// Set uses given function f to mock the Repository.Remove method
func (mmRemove *mRepositoryMockRemove) Set(f func(ctx context.Context, userID int64, skuID int64) (err error)) *RepositoryMock {
	if mmRemove.defaultExpectation != nil {
		mmRemove.mock.t.Fatalf("Default expectation is already set for the Repository.Remove method")
	}

	if len(mmRemove.expectations) > 0 {
		mmRemove.mock.t.Fatalf("Some expectations are already set for the Repository.Remove method")
	}

	mmRemove.mock.funcRemove = f
	return mmRemove.mock
}

// When sets expectation for the Repository.Remove which will trigger the result defined by the following
// Then helper
func (mmRemove *mRepositoryMockRemove) When(ctx context.Context, userID int64, skuID int64) *RepositoryMockRemoveExpectation {
	if mmRemove.mock.funcRemove != nil {
		mmRemove.mock.t.Fatalf("RepositoryMock.Remove mock is already set by Set")
	}

	expectation := &RepositoryMockRemoveExpectation{
		mock:   mmRemove.mock,
		params: &RepositoryMockRemoveParams{ctx, userID, skuID},
	}
	mmRemove.expectations = append(mmRemove.expectations, expectation)
	return expectation
}

// Then sets up Repository.Remove return parameters for the expectation previously defined by the When method
func (e *RepositoryMockRemoveExpectation) Then(err error) *RepositoryMock {
	e.results = &RepositoryMockRemoveResults{err}
	return e.mock
}

// Remove implements Repository
func (mmRemove *RepositoryMock) Remove(ctx context.Context, userID int64, skuID int64) (err error) {
	mm_atomic.AddUint64(&mmRemove.beforeRemoveCounter, 1)
	defer mm_atomic.AddUint64(&mmRemove.afterRemoveCounter, 1)

	if mmRemove.inspectFuncRemove != nil {
		mmRemove.inspectFuncRemove(ctx, userID, skuID)
	}

	mm_params := RepositoryMockRemoveParams{ctx, userID, skuID}

	// Record call args
	mmRemove.RemoveMock.mutex.Lock()
	mmRemove.RemoveMock.callArgs = append(mmRemove.RemoveMock.callArgs, &mm_params)
	mmRemove.RemoveMock.mutex.Unlock()

	for _, e := range mmRemove.RemoveMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRemove.RemoveMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRemove.RemoveMock.defaultExpectation.Counter, 1)
		mm_want := mmRemove.RemoveMock.defaultExpectation.params
		mm_got := RepositoryMockRemoveParams{ctx, userID, skuID}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRemove.t.Errorf("RepositoryMock.Remove got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRemove.RemoveMock.defaultExpectation.results
		if mm_results == nil {
			mmRemove.t.Fatal("No results are set for the RepositoryMock.Remove")
		}
		return (*mm_results).err
	}
	if mmRemove.funcRemove != nil {
		return mmRemove.funcRemove(ctx, userID, skuID)
	}
	mmRemove.t.Fatalf("Unexpected call to RepositoryMock.Remove. %v %v %v", ctx, userID, skuID)
	return
}

// RemoveAfterCounter returns a count of finished RepositoryMock.Remove invocations
func (mmRemove *RepositoryMock) RemoveAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRemove.afterRemoveCounter)
}

// RemoveBeforeCounter returns a count of RepositoryMock.Remove invocations
func (mmRemove *RepositoryMock) RemoveBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRemove.beforeRemoveCounter)
}

// Calls returns a list of arguments used in each call to RepositoryMock.Remove.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRemove *mRepositoryMockRemove) Calls() []*RepositoryMockRemoveParams {
	mmRemove.mutex.RLock()

	argCopy := make([]*RepositoryMockRemoveParams, len(mmRemove.callArgs))
	copy(argCopy, mmRemove.callArgs)

	mmRemove.mutex.RUnlock()

	return argCopy
}

// MinimockRemoveDone returns true if the count of the Remove invocations corresponds
// the number of defined expectations
func (m *RepositoryMock) MinimockRemoveDone() bool {
	for _, e := range m.RemoveMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoveMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRemoveCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemove != nil && mm_atomic.LoadUint64(&m.afterRemoveCounter) < 1 {
		return false
	}
	return true
}

// MinimockRemoveInspect logs each unmet expectation
func (m *RepositoryMock) MinimockRemoveInspect() {
	for _, e := range m.RemoveMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RepositoryMock.Remove with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoveMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRemoveCounter) < 1 {
		if m.RemoveMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RepositoryMock.Remove")
		} else {
			m.t.Errorf("Expected call to RepositoryMock.Remove with params: %#v", *m.RemoveMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemove != nil && mm_atomic.LoadUint64(&m.afterRemoveCounter) < 1 {
		m.t.Error("Expected call to RepositoryMock.Remove")
	}
}

type mRepositoryMockRemoveByUserID struct {
	mock               *RepositoryMock
	defaultExpectation *RepositoryMockRemoveByUserIDExpectation
	expectations       []*RepositoryMockRemoveByUserIDExpectation

	callArgs []*RepositoryMockRemoveByUserIDParams
	mutex    sync.RWMutex
}

// RepositoryMockRemoveByUserIDExpectation specifies expectation struct of the Repository.RemoveByUserID
type RepositoryMockRemoveByUserIDExpectation struct {
	mock    *RepositoryMock
	params  *RepositoryMockRemoveByUserIDParams
	results *RepositoryMockRemoveByUserIDResults
	Counter uint64
}

// RepositoryMockRemoveByUserIDParams contains parameters of the Repository.RemoveByUserID
type RepositoryMockRemoveByUserIDParams struct {
	ctx    context.Context
	userID int64
}

// RepositoryMockRemoveByUserIDResults contains results of the Repository.RemoveByUserID
type RepositoryMockRemoveByUserIDResults struct {
	err error
}

// Expect sets up expected params for Repository.RemoveByUserID
func (mmRemoveByUserID *mRepositoryMockRemoveByUserID) Expect(ctx context.Context, userID int64) *mRepositoryMockRemoveByUserID {
	if mmRemoveByUserID.mock.funcRemoveByUserID != nil {
		mmRemoveByUserID.mock.t.Fatalf("RepositoryMock.RemoveByUserID mock is already set by Set")
	}

	if mmRemoveByUserID.defaultExpectation == nil {
		mmRemoveByUserID.defaultExpectation = &RepositoryMockRemoveByUserIDExpectation{}
	}

	mmRemoveByUserID.defaultExpectation.params = &RepositoryMockRemoveByUserIDParams{ctx, userID}
	for _, e := range mmRemoveByUserID.expectations {
		if minimock.Equal(e.params, mmRemoveByUserID.defaultExpectation.params) {
			mmRemoveByUserID.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRemoveByUserID.defaultExpectation.params)
		}
	}

	return mmRemoveByUserID
}

// Inspect accepts an inspector function that has same arguments as the Repository.RemoveByUserID
func (mmRemoveByUserID *mRepositoryMockRemoveByUserID) Inspect(f func(ctx context.Context, userID int64)) *mRepositoryMockRemoveByUserID {
	if mmRemoveByUserID.mock.inspectFuncRemoveByUserID != nil {
		mmRemoveByUserID.mock.t.Fatalf("Inspect function is already set for RepositoryMock.RemoveByUserID")
	}

	mmRemoveByUserID.mock.inspectFuncRemoveByUserID = f

	return mmRemoveByUserID
}

// Return sets up results that will be returned by Repository.RemoveByUserID
func (mmRemoveByUserID *mRepositoryMockRemoveByUserID) Return(err error) *RepositoryMock {
	if mmRemoveByUserID.mock.funcRemoveByUserID != nil {
		mmRemoveByUserID.mock.t.Fatalf("RepositoryMock.RemoveByUserID mock is already set by Set")
	}

	if mmRemoveByUserID.defaultExpectation == nil {
		mmRemoveByUserID.defaultExpectation = &RepositoryMockRemoveByUserIDExpectation{mock: mmRemoveByUserID.mock}
	}
	mmRemoveByUserID.defaultExpectation.results = &RepositoryMockRemoveByUserIDResults{err}
	return mmRemoveByUserID.mock
}

// Set uses given function f to mock the Repository.RemoveByUserID method
func (mmRemoveByUserID *mRepositoryMockRemoveByUserID) Set(f func(ctx context.Context, userID int64) (err error)) *RepositoryMock {
	if mmRemoveByUserID.defaultExpectation != nil {
		mmRemoveByUserID.mock.t.Fatalf("Default expectation is already set for the Repository.RemoveByUserID method")
	}

	if len(mmRemoveByUserID.expectations) > 0 {
		mmRemoveByUserID.mock.t.Fatalf("Some expectations are already set for the Repository.RemoveByUserID method")
	}

	mmRemoveByUserID.mock.funcRemoveByUserID = f
	return mmRemoveByUserID.mock
}

// When sets expectation for the Repository.RemoveByUserID which will trigger the result defined by the following
// Then helper
func (mmRemoveByUserID *mRepositoryMockRemoveByUserID) When(ctx context.Context, userID int64) *RepositoryMockRemoveByUserIDExpectation {
	if mmRemoveByUserID.mock.funcRemoveByUserID != nil {
		mmRemoveByUserID.mock.t.Fatalf("RepositoryMock.RemoveByUserID mock is already set by Set")
	}

	expectation := &RepositoryMockRemoveByUserIDExpectation{
		mock:   mmRemoveByUserID.mock,
		params: &RepositoryMockRemoveByUserIDParams{ctx, userID},
	}
	mmRemoveByUserID.expectations = append(mmRemoveByUserID.expectations, expectation)
	return expectation
}

// Then sets up Repository.RemoveByUserID return parameters for the expectation previously defined by the When method
func (e *RepositoryMockRemoveByUserIDExpectation) Then(err error) *RepositoryMock {
	e.results = &RepositoryMockRemoveByUserIDResults{err}
	return e.mock
}

// RemoveByUserID implements Repository
func (mmRemoveByUserID *RepositoryMock) RemoveByUserID(ctx context.Context, userID int64) (err error) {
	mm_atomic.AddUint64(&mmRemoveByUserID.beforeRemoveByUserIDCounter, 1)
	defer mm_atomic.AddUint64(&mmRemoveByUserID.afterRemoveByUserIDCounter, 1)

	if mmRemoveByUserID.inspectFuncRemoveByUserID != nil {
		mmRemoveByUserID.inspectFuncRemoveByUserID(ctx, userID)
	}

	mm_params := RepositoryMockRemoveByUserIDParams{ctx, userID}

	// Record call args
	mmRemoveByUserID.RemoveByUserIDMock.mutex.Lock()
	mmRemoveByUserID.RemoveByUserIDMock.callArgs = append(mmRemoveByUserID.RemoveByUserIDMock.callArgs, &mm_params)
	mmRemoveByUserID.RemoveByUserIDMock.mutex.Unlock()

	for _, e := range mmRemoveByUserID.RemoveByUserIDMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRemoveByUserID.RemoveByUserIDMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRemoveByUserID.RemoveByUserIDMock.defaultExpectation.Counter, 1)
		mm_want := mmRemoveByUserID.RemoveByUserIDMock.defaultExpectation.params
		mm_got := RepositoryMockRemoveByUserIDParams{ctx, userID}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRemoveByUserID.t.Errorf("RepositoryMock.RemoveByUserID got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRemoveByUserID.RemoveByUserIDMock.defaultExpectation.results
		if mm_results == nil {
			mmRemoveByUserID.t.Fatal("No results are set for the RepositoryMock.RemoveByUserID")
		}
		return (*mm_results).err
	}
	if mmRemoveByUserID.funcRemoveByUserID != nil {
		return mmRemoveByUserID.funcRemoveByUserID(ctx, userID)
	}
	mmRemoveByUserID.t.Fatalf("Unexpected call to RepositoryMock.RemoveByUserID. %v %v", ctx, userID)
	return
}

// RemoveByUserIDAfterCounter returns a count of finished RepositoryMock.RemoveByUserID invocations
func (mmRemoveByUserID *RepositoryMock) RemoveByUserIDAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRemoveByUserID.afterRemoveByUserIDCounter)
}

// RemoveByUserIDBeforeCounter returns a count of RepositoryMock.RemoveByUserID invocations
func (mmRemoveByUserID *RepositoryMock) RemoveByUserIDBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRemoveByUserID.beforeRemoveByUserIDCounter)
}

// Calls returns a list of arguments used in each call to RepositoryMock.RemoveByUserID.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRemoveByUserID *mRepositoryMockRemoveByUserID) Calls() []*RepositoryMockRemoveByUserIDParams {
	mmRemoveByUserID.mutex.RLock()

	argCopy := make([]*RepositoryMockRemoveByUserIDParams, len(mmRemoveByUserID.callArgs))
	copy(argCopy, mmRemoveByUserID.callArgs)

	mmRemoveByUserID.mutex.RUnlock()

	return argCopy
}

// MinimockRemoveByUserIDDone returns true if the count of the RemoveByUserID invocations corresponds
// the number of defined expectations
func (m *RepositoryMock) MinimockRemoveByUserIDDone() bool {
	for _, e := range m.RemoveByUserIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoveByUserIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRemoveByUserIDCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemoveByUserID != nil && mm_atomic.LoadUint64(&m.afterRemoveByUserIDCounter) < 1 {
		return false
	}
	return true
}

// MinimockRemoveByUserIDInspect logs each unmet expectation
func (m *RepositoryMock) MinimockRemoveByUserIDInspect() {
	for _, e := range m.RemoveByUserIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RepositoryMock.RemoveByUserID with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoveByUserIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRemoveByUserIDCounter) < 1 {
		if m.RemoveByUserIDMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RepositoryMock.RemoveByUserID")
		} else {
			m.t.Errorf("Expected call to RepositoryMock.RemoveByUserID with params: %#v", *m.RemoveByUserIDMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemoveByUserID != nil && mm_atomic.LoadUint64(&m.afterRemoveByUserIDCounter) < 1 {
		m.t.Error("Expected call to RepositoryMock.RemoveByUserID")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *RepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockInsertInspect()

			m.MinimockListInspect()

			m.MinimockRemoveInspect()

			m.MinimockRemoveByUserIDInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *RepositoryMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *RepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockInsertDone() &&
		m.MinimockListDone() &&
		m.MinimockRemoveDone() &&
		m.MinimockRemoveByUserIDDone()
}
