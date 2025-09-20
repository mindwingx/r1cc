package tenant

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/metric"
	"microservice/internal/adapter/queue"
	"microservice/internal/adapter/trace"
	"microservice/internal/domain"
	"microservice/internal/modules/port"
	"microservice/pkg/meta"
	"microservice/pkg/meta/status"
)

type (
	ITenantHttpHandler interface {
		Create(c echo.Context) error
		Details(c echo.Context) error
		List(c echo.Context) error
	}

	HandlerFx struct {
		fx.In
		Locale   locale.ILocale
		Tracer   trace.ITracer
		Logger   logger.ILogger
		Metric   metric.IMetric
		Queue    queue.IQueue
		TenantUC port.ITenantUsecase
	}

	Handler struct {
		l        locale.ILocale
		trc      trace.ITracer
		lgr      logger.ILogger
		metric   metric.IMetric
		queue    queue.IQueue
		tenantUC port.ITenantUsecase
	}
)

func NewHttpHandlerFx(fx HandlerFx) ITenantHttpHandler {
	return &Handler{
		l:        fx.Locale,
		trc:      fx.Tracer,
		lgr:      fx.Logger,
		metric:   fx.Metric,
		queue:    fx.Queue,
		tenantUC: fx.TenantUC,
	}
}

// Create godoc
// @Summary Create New Tenant
// @Tags Tenant
// @Accept json
// @Produce json
// @Param Request body tenant.CreateRequest true "necessary fields for request"
// @Success 201 {object} meta.Response{data=tenant.CreateResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	404 {object} meta.Response{data=nil} "not found"
// @Failure	409 {object} meta.Response{data=nil} "already exists"
// @Failure	422 {object} meta.Response{data=nil} "unprocessable"
// @Router /api/v1/tenant/create [post]
func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := meta.ReqBodyToDomain[*CreateRequest, domain.Tenant](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	res, ucErr := h.tenantUC.Create(ctx, req)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	return meta.Resp(c, h.l).Status(status.Created).Data(CreateResp(res)).Json()
}

// Details godoc
// @Summary Get Tenant Details
// @Tags Tenant
// @Accept json
// @Produce json
// @Param uuid path string true "Tenant UUID" example(f81eee2d-2cca-4169-8062-7404a78d5c3b)
// @Success 200 {object}  meta.Response{data=tenant.DetailsResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	404 {object} meta.Response{data=nil} "no Tenant found"
// @Failure	422 {object} meta.Response{data=nil} "invalid data types"
// @Router /api/v1/tenant/{uuid} [get]
func (h *Handler) Details(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := meta.ReqRouteParamsToDomain[*DetailsRequest, domain.Tenant](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	details, ucErr := h.tenantUC.GetDetails(ctx, req)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	return meta.Resp(c, h.l).Status(status.Success).Data(DetailsResp(details)).Json()
}

// List godoc
// @Summary Get Tenant List
// @Tags Tenant
// @Accept json
// @Produce json
// @Param page query int false "Page Number"
// @Param limit query int false "Page Limit"
// @Param sort query string false "id, username, tenant_name, created_at, updated_at\n(other valid columns are acceptable)"
// @Param order query string false "`asc` or `desc`"
// @Param search query string false "Search the Tenant Username and Name"
// @Success 200 {object} meta.Response{data=tenant.ListResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	422 {object} meta.Response{data=nil} "database error while retrieving"
// @Router /api/v1/tenant/list [get]
func (h *Handler) List(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := meta.ReqQryParamToDomain[*ListQryRequest, domain.TenantListReqQryParam](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	res, err := h.tenantUC.GetList(ctx, req)
	if err != nil {
		return meta.Resp(c, h.l).Status(status.Failed).Err(err).Json()
	}

	return meta.Resp(c, h.l).Status(status.Success).Data(ListResp(req, res)).Json()
}
