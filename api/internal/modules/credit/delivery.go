package credit

import (
	"context"
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
	ICreditHttpHandler interface {
		IncreaseCredit(c echo.Context) error
		TransactionsList(c echo.Context) error
	}

	HandlerFx struct {
		fx.In
		Locale   locale.ILocale
		Tracer   trace.ITracer
		Logger   logger.ILogger
		Metric   metric.IMetric
		Queue    queue.IQueue
		TenantUC port.ITenantUsecase
		CreditUC port.ICreditUsecase
	}

	Handler struct {
		l        locale.ILocale
		trc      trace.ITracer
		lgr      logger.ILogger
		metric   metric.IMetric
		queue    queue.IQueue
		tenantUC port.ITenantUsecase
		creditUC port.ICreditUsecase
	}
)

func NewHttpHandlerFx(fx HandlerFx) ICreditHttpHandler {
	return &Handler{
		l:        fx.Locale,
		trc:      fx.Tracer,
		lgr:      fx.Logger,
		metric:   fx.Metric,
		queue:    fx.Queue,
		tenantUC: fx.TenantUC,
		creditUC: fx.CreditUC,
	}
}

// IncreaseCredit godoc
// @Summary Increase Tenant Credit
// @Tags Credit
// @Accept json
// @Produce json
// @Security Bearer
// @Param X.TENANT.UUID header string true "Tenant UUID" example(f81eee2d-2cca-4169-8062-7404a78d5c3b)
// @Param Request body credit.IncreaseCreditRequest true "necessary fields for request"
// @Success 201 {object} meta.Response{data=credit.IncreaseCreditResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	404 {object} meta.Response{data=nil} "not found"
// @Failure	409 {object} meta.Response{data=nil} "already exists"
// @Failure	422 {object} meta.Response{data=nil} "unprocessable"
// @Router /api/v1/credit/increase [post]
func (h *Handler) IncreaseCredit(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := meta.ReqHeaderToDomain[*dto.TenantUuid, domain.Tenant](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	transaction, err := meta.ReqBodyToDomain[*IncreaseCreditRequest, domain.Transaction](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	// as the header reads the key, it reallocates value to be used during request cycle
	ctx = context.WithValue(ctx, "X.TENANT.UUID", req.UUID().String())

	tenant, ucErr := h.tenantUC.GetDetails(ctx, req)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	credit := tenant.Credit()
	transaction.SetCreditID(credit.ID())
	credit.SetTxAmount(transaction)

	res, ucErr := h.creditUC.IncreaseAmount(ctx, credit)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	return meta.Resp(c, h.l).Status(status.Success).Data(IncreaseCreditResp(res)).Json()
}

// TransactionsList godoc
// @Summary Get Tenant Credit and Transaction List
// @Tags Credit
// @Accept json
// @Produce json
// @Security Bearer
// @Param X.TENANT.UUID header string true "Tenant UUID" example(f81eee2d-2cca-4169-8062-7404a78d5c3b)
// @Param page query int false "Page Number"
// @Param limit query int false "Page Limit"
// @Param order query string false "`asc` or `desc` based on `created_at`"
// @Success 200 {object} meta.Response{data=credit.ListResponse, error=nil} "success response"
// @Failure	400 {object} meta.Response{data=nil} "process failure"
// @Failure	404 {object} meta.Response{data=nil} "not found"
// @Failure	409 {object} meta.Response{data=nil} "already exists"
// @Failure	422 {object} meta.Response{data=nil} "unprocessable"
// @Router /api/v1/credit/transactions [get]
func (h *Handler) TransactionsList(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := meta.ReqHeaderToDomain[*dto.TenantUuid, domain.Tenant](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	qp, err := meta.ReqQryParamToDomain[*ListQryRequest, domain.TransactionListReqQryParam](c)
	if err != nil {
		return meta.Resp(c, h.l).ServiceErr(err).Json()
	}

	tenant, ucErr := h.tenantUC.GetDetails(ctx, req)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	res, ucErr := h.creditUC.GetDetails(ctx, tenant.Credit(), qp)
	if ucErr != nil {
		return meta.Resp(c, h.l).ServiceErr(ucErr).Json()
	}

	return meta.Resp(c, h.l).Status(status.Success).Data(ListResp(qp, res)).Json()
}
