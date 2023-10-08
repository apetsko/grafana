/*Package api contains base API implementation of unified alerting
 *
 *Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 *
 *Do not manually edit these files, please find ngalert/api/swagger-codegen/ for commands on how to generate them.
 */
package api

import (
	"net/http"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/middleware/requestmeta"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/ngalert/metrics"
	"github.com/grafana/grafana/pkg/web"
)

type UpgradeApi interface {
	RouteDeleteOrgUpgrade(*contextmodel.ReqContext) response.Response
	RouteGetOrgUpgrade(*contextmodel.ReqContext) response.Response
	RoutePostUpgradeAlert(*contextmodel.ReqContext) response.Response
	RoutePostUpgradeAllChannels(*contextmodel.ReqContext) response.Response
	RoutePostUpgradeAllDashboards(*contextmodel.ReqContext) response.Response
	RoutePostUpgradeChannel(*contextmodel.ReqContext) response.Response
	RoutePostUpgradeDashboard(*contextmodel.ReqContext) response.Response
	RoutePostUpgradeOrg(*contextmodel.ReqContext) response.Response
}

func (f *UpgradeApiHandler) RouteDeleteOrgUpgrade(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteDeleteOrgUpgrade(ctx)
}
func (f *UpgradeApiHandler) RouteGetOrgUpgrade(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetOrgUpgrade(ctx)
}
func (f *UpgradeApiHandler) RoutePostUpgradeAlert(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	dashboardIDParam := web.Params(ctx.Req)[":DashboardID"]
	panelIDParam := web.Params(ctx.Req)[":PanelID"]
	return f.handleRoutePostUpgradeAlert(ctx, dashboardIDParam, panelIDParam)
}
func (f *UpgradeApiHandler) RoutePostUpgradeAllChannels(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRoutePostUpgradeAllChannels(ctx)
}
func (f *UpgradeApiHandler) RoutePostUpgradeAllDashboards(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRoutePostUpgradeAllDashboards(ctx)
}
func (f *UpgradeApiHandler) RoutePostUpgradeChannel(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	channelIDParam := web.Params(ctx.Req)[":ChannelID"]
	return f.handleRoutePostUpgradeChannel(ctx, channelIDParam)
}
func (f *UpgradeApiHandler) RoutePostUpgradeDashboard(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	dashboardIDParam := web.Params(ctx.Req)[":DashboardID"]
	return f.handleRoutePostUpgradeDashboard(ctx, dashboardIDParam)
}
func (f *UpgradeApiHandler) RoutePostUpgradeOrg(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRoutePostUpgradeOrg(ctx)
}

func (api *API) RegisterUpgradeApiEndpoints(srv UpgradeApi, m *metrics.API) {
	api.RouteRegister.Group("", func(group routing.RouteRegister) {
		group.Delete(
			toMacaronPath("/api/v1/upgrade/org"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			api.authorize(http.MethodDelete, "/api/v1/upgrade/org"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/upgrade/org",
				api.Hooks.Wrap(srv.RouteDeleteOrgUpgrade),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/upgrade/org"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			api.authorize(http.MethodGet, "/api/v1/upgrade/org"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/upgrade/org",
				api.Hooks.Wrap(srv.RouteGetOrgUpgrade),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/upgrade/dashboards/{DashboardID}/panels/{PanelID}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			api.authorize(http.MethodPost, "/api/v1/upgrade/dashboards/{DashboardID}/panels/{PanelID}"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/upgrade/dashboards/{DashboardID}/panels/{PanelID}",
				api.Hooks.Wrap(srv.RoutePostUpgradeAlert),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/upgrade/channels"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			api.authorize(http.MethodPost, "/api/v1/upgrade/channels"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/upgrade/channels",
				api.Hooks.Wrap(srv.RoutePostUpgradeAllChannels),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/upgrade/dashboards"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			api.authorize(http.MethodPost, "/api/v1/upgrade/dashboards"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/upgrade/dashboards",
				api.Hooks.Wrap(srv.RoutePostUpgradeAllDashboards),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/upgrade/channels/{ChannelID}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			api.authorize(http.MethodPost, "/api/v1/upgrade/channels/{ChannelID}"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/upgrade/channels/{ChannelID}",
				api.Hooks.Wrap(srv.RoutePostUpgradeChannel),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/upgrade/dashboards/{DashboardID}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			api.authorize(http.MethodPost, "/api/v1/upgrade/dashboards/{DashboardID}"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/upgrade/dashboards/{DashboardID}",
				api.Hooks.Wrap(srv.RoutePostUpgradeDashboard),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/upgrade/org"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			api.authorize(http.MethodPost, "/api/v1/upgrade/org"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/upgrade/org",
				api.Hooks.Wrap(srv.RoutePostUpgradeOrg),
				m,
			),
		)
	}, middleware.ReqSignedIn)
}
