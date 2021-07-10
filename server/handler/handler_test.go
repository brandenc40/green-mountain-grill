package handler

// Basic imports
import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	gmg "github.com/brandenc40/green-mountain-grill"
	"github.com/brandenc40/green-mountain-grill/mocks"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type HandlerTestSuite struct {
	suite.Suite

	handler    *Handler
	repoMock   *mocks.Repository
	clientMock *mocks.Client
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.clientMock = &mocks.Client{}
	suite.repoMock = &mocks.Repository{}
	params := Params{
		Logger:      logrus.New(),
		GrillClient: suite.clientMock,
		Repository:  suite.repoMock,
		Scheduler:   gocron.NewScheduler(),
	}
	suite.handler = New(params)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestGetGrillState_Success() {
	state := gmg.State{}
	suite.clientMock.On("GetState").Return(&state, nil)
	ctx := fasthttp.RequestCtx{}
	suite.handler.GetGrillState(&ctx)
	expected := `{"current_temperature":0,"target_temperature":0,"probe1_temperature":0,"probe1_target_temperature":0,"probe2_temperature":0,"probe2_target_temperature":0,"warn_code":"WarnCodeNone","power_state":"PowerStateOff","fire_state":"FireStateDefault"}`
	suite.Equal(expected, string(ctx.Response.Body()))
	suite.Equal(200, ctx.Response.StatusCode())
	suite.clientMock.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetGrillState_Error() {
	suite.clientMock.On("GetState").Return(nil, errors.New("error"))
	ctx := fasthttp.RequestCtx{}
	suite.handler.GetGrillState(&ctx)
	suite.Equal("error", string(ctx.Response.Body()))
	suite.Equal(fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
	suite.clientMock.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestGetGrillState_ErrorUnavailable() {
	suite.clientMock.On("GetState").Return(nil, gmg.GrillUnreachableErr{Err: errors.New("error")})
	ctx := fasthttp.RequestCtx{}
	suite.handler.GetGrillState(&ctx)
	suite.Equal("grill is unreachable: error", string(ctx.Response.Body()))
	suite.Equal(fasthttp.StatusServiceUnavailable, ctx.Response.StatusCode())
	suite.clientMock.AssertExpectations(suite.T())
}

func (suite *HandlerTestSuite) TestNewSession_Success() {
	suite.clientMock.On("IsAvailable").Return(true)
	ctx := fasthttp.RequestCtx{}
	suite.handler.NewSession(&ctx)
	_, err := uuid.Parse(string(ctx.Response.Body()))
	suite.NoError(err)
	suite.Equal(200, ctx.Response.StatusCode())
	time.Sleep(time.Millisecond) // allow goroutine to start
	suite.NotNil(suite.handler.stopChannel)
	suite.True(suite.handler.isMonitoring)
	suite.clientMock.AssertExpectations(suite.T())
	suite.repoMock.AssertExpectations(suite.T())
	suite.handler.stopMonitoringGrill()
	suite.False(suite.handler.isMonitoring)
	suite.Nil(suite.handler.stopChannel)
}
