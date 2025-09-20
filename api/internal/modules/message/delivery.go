package message

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/metric"
	"microservice/internal/adapter/queue"
	"microservice/internal/adapter/trace"
	"microservice/internal/domain"
	"microservice/internal/modules/dto"
	"microservice/internal/modules/port"
	"microservice/pkg/meta"
	"microservice/pkg/meta/status"
)

type (
	IMessageHttpHandler interface {
		Send(c echo.Context) error
		List(c echo.Context) error
	}

	HandlerFx struct {
		fx.In
		Locale    locale.ILocale
		Tracer    trace.ITracer
		Logger    logger.ILogger
		Metric    metric.IMetric
		Queue     queue.IQueue
		TenantUC  port.ITenantUsecase
		MessageUC port.IMessageUsecase
	}

	Handler struct {
		l         locale.ILocale
		trc       trace.ITracer
		lgr       logger.ILogger
		metric    metric.IMetric
		queue     queue.IQueue
		tenantUC  port.ITenantUsecase
		messageUC port.IMessageUsecase
	}
)

func NewHttpHandlerFx(fx HandlerFx) IMessageHttpHandler {
	return &Handler{
		l:         fx.Locale,
		trc:       fx.Tracer,
		lgr:       fx.Logger,
		metric:    fx.Metric,
		queue:     fx.Queue,
		tenantUC:  fx.TenantUC,
		messageUC: fx.MessageUC,
	}
}

// Send godoc
// @Summary Send Message
// @Description request body channel values `event.prod` or `event.express`
// @Tags Message
// @Accept json
// @Produce json
// @Security Bearer
// @Param X.TENANT.UUID header string true "Tenant UUID" example(f81eee2d-2cca-4169-8062-7404a78d5c3b)
// @Param Request body message.SendMessageRequest true "necessary fields for request"
// @Success 201 {object} meta.Response{error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	404 {object} meta.Response{data=nil} "not found"
// @Failure	409 {object} meta.Response{data=nil} "already exists"
// @Failure	422 {object} meta.Response{data=nil} "unprocessable"
// @Router /api/v1/message/send [post]
func (h *Handler) Send(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := meta.ReqHeaderToDomain[*dto.TenantUuid, domain.Tenant](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	message, err := meta.ReqBodyToDomain[*SendMessageRequest, domain.Message](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	tenant, ucErr := h.tenantUC.GetDetails(ctx, req)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	message.SetTenantID(tenant.ID())

	_, ucErr = h.messageUC.Send(ctx, tenant, message)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	return meta.Resp(c, h.l).Status(status.Success).Json()
}

// List godoc
// @Summary Get Sent Message List
// @Tags Message
// @Accept json
// @Produce json
// @Param X.TENANT.UUID header string true "Tenant UUID" example(f81eee2d-2cca-4169-8062-7404a78d5c3b)
// @Param page query int false "Page Number"
// @Param limit query int false "Page Limit"
// @Param sort query string false "id, created_at, updated_at\n(other valid columns are acceptable)"
// @Param order query string false "`asc` or `desc`"
// @Param search query string false "Search the Message"
// @Success 200 {object} meta.Response{data=message.ListResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	422 {object} meta.Response{data=nil} "database error while retrieving"
// @Router /api/v1/message/list [get]
func (h *Handler) List(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := meta.ReqHeaderToDomain[*dto.TenantUuid, domain.Tenant](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	list, err := meta.ReqQryParamToDomain[*ListQryRequest, domain.MessageListReqQryParam](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	tenant, ucErr := h.tenantUC.GetDetails(ctx, req)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	list.SetTenantId(tenant.ID())

	res, err := h.messageUC.GetList(ctx, list)
	if err != nil {
		return meta.Resp(c, h.l).Status(status.Failed).Err(err).Json()
	}

	return meta.Resp(c, h.l).Status(status.Success).Data(ListResp(list, res)).Json()
}
