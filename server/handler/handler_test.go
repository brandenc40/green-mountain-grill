package handler

// Basic imports
import (
	"testing"

	gmg "github.com/brandenc40/green-mountain-grill"
	"github.com/brandenc40/green-mountain-grill/mocks"
	repoMock "github.com/brandenc40/green-mountain-grill/mocks/respository"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type HandlerTestSuite struct {
	suite.Suite

	handler    *Handler
	repoMock   *repoMock.Repository
	clientMock *mocks.Client
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.clientMock = &mocks.Client{}
	suite.repoMock = &repoMock.Repository{}
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
