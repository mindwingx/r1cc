package app

import (
	"context"
	"fmt"
	"microservice/internal/adapter/cache"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/metric"
	"microservice/internal/adapter/orm"
	"microservice/internal/adapter/queue"
	"microservice/internal/adapter/registry"
	"microservice/internal/adapter/trace"
	"microservice/pkg/utils"
)

type Clients struct {
	registry     registry.IRegistry
	tracer       trace.IGTracer
	logIngest    logger.ILog
	globalLogger logger.IGLogger
	metric       metric.IMetrics
	locale       locale.ILocale
	cache        cache.ICache
	database     orm.ISqlGeneric
	queue        queue.IQueue
}

func NewClients() *Clients { return &Clients{} }

func (c *Clients) Registry() registry.IRegistry            { return c.registry }
func (c *Clients) SetRegistry(registry registry.IRegistry) { c.registry = registry }

func (c *Clients) Tracer() trace.IGTracer          { return c.tracer }
func (c *Clients) SetTracer(tracer trace.IGTracer) { c.tracer = tracer }

func (c *Clients) Logger() logger.ILog          { return c.logIngest }
func (c *Clients) SetLogger(logger logger.ILog) { c.logIngest = logger }

func (c *Clients) GlobalLogger() logger.IGLogger                { return c.globalLogger }
func (c *Clients) SetGlobalLogger(globalLogger logger.IGLogger) { c.globalLogger = globalLogger }

func (c *Clients) Metric() metric.IMetrics          { return c.metric }
func (c *Clients) SetMetric(metric metric.IMetrics) { c.metric = metric }

func (c *Clients) Locale() locale.ILocale          { return c.locale }
func (c *Clients) SetLocale(locale locale.ILocale) { c.locale = locale }

func (c *Clients) Cache() cache.ICache         { return c.cache }
func (c *Clients) SetCache(cache cache.ICache) { c.cache = cache }

func (c *Clients) Database() orm.ISqlGeneric            { return c.database }
func (c *Clients) SetDatabase(database orm.ISqlGeneric) { c.database = database }

func (c *Clients) Queue() queue.IQueue         { return c.queue }
func (c *Clients) SetQueue(queue queue.IQueue) { c.queue = queue }

// init

func (a *App) InitClients() {
	a.SetClient(NewClients())
	a.initRegistry()
	a.initService()
	a.initTrace()
	a.initLogger()
	a.initMetric()
	a.initLocale()
	a.initCache()
	a.initDatabase()
	a.initQueue()
}

func (a *App) initRegistry() {
	a.Client().SetRegistry(registry.New())
	a.Client().Registry().Init(registry.ConfigTypeEnv, fmt.Sprintf("%s/.env", utils.Root()))
}

func (a *App) initTrace() {
	a.Client().SetTracer(trace.New(*a.Config(), a.Client().Registry()))
	a.Client().Tracer().Init()

	a.SetTrace(a.Client().Tracer().SpanByCtx(context.Background(), "service", "init"))
	a.Span().AddEvent("otel initialized")
}

func (a *App) initLogger() {
	a.Client().SetLogger(logger.NewIngest(*a.Config(), a.Client().Registry()))
	a.Client().Logger().Init()
	a.Span().AddEvent("logstash initialized")

	a.Client().SetGlobalLogger(logger.New(*a.Config(), a.Client().Logger()))
	a.Client().GlobalLogger().Init()
	a.Span().AddEvent("zap initialized")
}

func (a *App) initLocale() {
	a.Client().SetLocale(locale.New(a.Client().Registry()))
	a.Client().Locale().Init()
	a.Span().AddEvent("locale initialized")
}

func (a *App) initMetric() {
	a.Client().SetMetric(metric.New(*a.Config()))
	a.Client().Metric().Init()
	a.Span().AddEvent("metric initialized")
}

func (a *App) initCache() {
	a.Client().SetCache(cache.New(a.Client().Registry()))
	a.Client().Cache().Init()
	a.Span().AddEvent("cache initialized")
}

func (a *App) initDatabase() {
	a.Client().SetDatabase(orm.New(a.Config(), a.Client().Registry()))
	a.Client().Database().Init()
	a.Client().Database().Migrate(fmt.Sprintf("%s/schema/psql", utils.Root()))
	a.Span().AddEvent("database initialized")
}

func (a *App) initQueue() {
	a.Client().SetQueue(queue.New(a.Client().Registry()))
	a.Client().Queue().Init()
	a.Span().AddEvent("queue initialized")
}
