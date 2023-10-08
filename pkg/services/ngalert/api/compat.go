package api

import (
	"encoding/json"
	"time"

	"github.com/prometheus/common/model"

	"github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	migmodels "github.com/grafana/grafana/pkg/services/ngalert/migration/models"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/util"
)

// AlertRuleFromProvisionedAlertRule converts definitions.ProvisionedAlertRule to models.AlertRule
func AlertRuleFromProvisionedAlertRule(a definitions.ProvisionedAlertRule) (models.AlertRule, error) {
	return models.AlertRule{
		ID:           a.ID,
		UID:          a.UID,
		OrgID:        a.OrgID,
		NamespaceUID: a.FolderUID,
		RuleGroup:    a.RuleGroup,
		Title:        a.Title,
		Condition:    a.Condition,
		Data:         AlertQueriesFromApiAlertQueries(a.Data),
		Updated:      a.Updated,
		NoDataState:  models.NoDataState(a.NoDataState),          // TODO there must be a validation
		ExecErrState: models.ExecutionErrorState(a.ExecErrState), // TODO there must be a validation
		For:          time.Duration(a.For),
		Annotations:  a.Annotations,
		Labels:       a.Labels,
		IsPaused:     a.IsPaused,
	}, nil
}

// ProvisionedAlertRuleFromAlertRule converts models.AlertRule to definitions.ProvisionedAlertRule and sets provided provenance status
func ProvisionedAlertRuleFromAlertRule(rule models.AlertRule, provenance models.Provenance) definitions.ProvisionedAlertRule {
	return definitions.ProvisionedAlertRule{
		ID:           rule.ID,
		UID:          rule.UID,
		OrgID:        rule.OrgID,
		FolderUID:    rule.NamespaceUID,
		RuleGroup:    rule.RuleGroup,
		Title:        rule.Title,
		For:          model.Duration(rule.For),
		Condition:    rule.Condition,
		Data:         ApiAlertQueriesFromAlertQueries(rule.Data),
		Updated:      rule.Updated,
		NoDataState:  definitions.NoDataState(rule.NoDataState),          // TODO there may be a validation
		ExecErrState: definitions.ExecutionErrorState(rule.ExecErrState), // TODO there may be a validation
		Annotations:  rule.Annotations,
		Labels:       rule.Labels,
		Provenance:   definitions.Provenance(provenance), // TODO validate enum conversion?
		IsPaused:     rule.IsPaused,
	}
}

// ProvisionedAlertRuleFromAlertRules converts a collection of models.AlertRule to definitions.ProvisionedAlertRules with provenance status models.ProvenanceNone
func ProvisionedAlertRuleFromAlertRules(rules []*models.AlertRule, provenances map[string]models.Provenance) definitions.ProvisionedAlertRules {
	result := make([]definitions.ProvisionedAlertRule, 0, len(rules))
	for _, r := range rules {
		result = append(result, ProvisionedAlertRuleFromAlertRule(*r, provenances[r.UID]))
	}
	return result
}

// AlertQueriesFromApiAlertQueries converts a collection of definitions.AlertQuery to collection of models.AlertQuery
func AlertQueriesFromApiAlertQueries(queries []definitions.AlertQuery) []models.AlertQuery {
	result := make([]models.AlertQuery, 0, len(queries))
	for _, q := range queries {
		result = append(result, models.AlertQuery{
			RefID:     q.RefID,
			QueryType: q.QueryType,
			RelativeTimeRange: models.RelativeTimeRange{
				From: models.Duration(q.RelativeTimeRange.From),
				To:   models.Duration(q.RelativeTimeRange.To),
			},
			DatasourceUID: q.DatasourceUID,
			Model:         q.Model,
		})
	}
	return result
}

// ApiAlertQueriesFromAlertQueries converts a collection of models.AlertQuery to collection of definitions.AlertQuery
func ApiAlertQueriesFromAlertQueries(queries []models.AlertQuery) []definitions.AlertQuery {
	result := make([]definitions.AlertQuery, 0, len(queries))
	for _, q := range queries {
		result = append(result, definitions.AlertQuery{
			RefID:     q.RefID,
			QueryType: q.QueryType,
			RelativeTimeRange: definitions.RelativeTimeRange{
				From: definitions.Duration(q.RelativeTimeRange.From),
				To:   definitions.Duration(q.RelativeTimeRange.To),
			},
			DatasourceUID: q.DatasourceUID,
			Model:         q.Model,
		})
	}
	return result
}

func AlertRuleGroupFromApiAlertRuleGroup(a definitions.AlertRuleGroup) (models.AlertRuleGroup, error) {
	ruleGroup := models.AlertRuleGroup{
		Title:     a.Title,
		FolderUID: a.FolderUID,
		Interval:  a.Interval,
	}
	for i := range a.Rules {
		converted, err := AlertRuleFromProvisionedAlertRule(a.Rules[i])
		if err != nil {
			return models.AlertRuleGroup{}, err
		}
		ruleGroup.Rules = append(ruleGroup.Rules, converted)
	}
	return ruleGroup, nil
}

func ApiAlertRuleGroupFromAlertRuleGroup(d models.AlertRuleGroup) definitions.AlertRuleGroup {
	rules := make([]definitions.ProvisionedAlertRule, 0, len(d.Rules))
	for i := range d.Rules {
		rules = append(rules, ProvisionedAlertRuleFromAlertRule(d.Rules[i], d.Provenance))
	}
	return definitions.AlertRuleGroup{
		Title:     d.Title,
		FolderUID: d.FolderUID,
		Interval:  d.Interval,
		Rules:     rules,
	}
}

// AlertingFileExportFromAlertRuleGroupWithFolderTitle creates an definitions.AlertingFileExport DTO from []models.AlertRuleGroupWithFolderTitle.
func AlertingFileExportFromAlertRuleGroupWithFolderTitle(groups []models.AlertRuleGroupWithFolderTitle) (definitions.AlertingFileExport, error) {
	f := definitions.AlertingFileExport{APIVersion: 1}
	for _, group := range groups {
		export, err := AlertRuleGroupExportFromAlertRuleGroupWithFolderTitle(group)
		if err != nil {
			return definitions.AlertingFileExport{}, err
		}
		f.Groups = append(f.Groups, export)
	}
	return f, nil
}

// AlertRuleGroupExportFromAlertRuleGroupWithFolderTitle creates a definitions.AlertRuleGroupExport DTO from models.AlertRuleGroup.
func AlertRuleGroupExportFromAlertRuleGroupWithFolderTitle(d models.AlertRuleGroupWithFolderTitle) (definitions.AlertRuleGroupExport, error) {
	rules := make([]definitions.AlertRuleExport, 0, len(d.Rules))
	for i := range d.Rules {
		alert, err := AlertRuleExportFromAlertRule(d.Rules[i])
		if err != nil {
			return definitions.AlertRuleGroupExport{}, err
		}
		rules = append(rules, alert)
	}
	return definitions.AlertRuleGroupExport{
		OrgID:           d.OrgID,
		Name:            d.Title,
		Folder:          d.FolderTitle,
		FolderUID:       d.FolderUID,
		Interval:        model.Duration(time.Duration(d.Interval) * time.Second),
		IntervalSeconds: d.Interval,
		Rules:           rules,
	}, nil
}

// AlertRuleExportFromAlertRule creates a definitions.AlertRuleExport DTO from models.AlertRule.
func AlertRuleExportFromAlertRule(rule models.AlertRule) (definitions.AlertRuleExport, error) {
	data := make([]definitions.AlertQueryExport, 0, len(rule.Data))
	for i := range rule.Data {
		query, err := AlertQueryExportFromAlertQuery(rule.Data[i])
		if err != nil {
			return definitions.AlertRuleExport{}, err
		}
		data = append(data, query)
	}

	result := definitions.AlertRuleExport{
		UID:          rule.UID,
		Title:        rule.Title,
		For:          model.Duration(rule.For),
		Condition:    rule.Condition,
		Data:         data,
		DashboardUID: rule.DashboardUID,
		PanelID:      rule.PanelID,
		NoDataState:  definitions.NoDataState(rule.NoDataState),
		ExecErrState: definitions.ExecutionErrorState(rule.ExecErrState),
		IsPaused:     rule.IsPaused,
	}
	if rule.For.Seconds() > 0 {
		result.ForSeconds = util.Pointer(int64(rule.For.Seconds()))
	}
	if rule.Annotations != nil {
		result.Annotations = &rule.Annotations
	}
	if rule.Labels != nil {
		result.Labels = &rule.Labels
	}
	return result, nil
}

// AlertQueryExportFromAlertQuery creates a definitions.AlertQueryExport DTO from models.AlertQuery.
func AlertQueryExportFromAlertQuery(query models.AlertQuery) (definitions.AlertQueryExport, error) {
	// We unmarshal the json.RawMessage model into a map in order to facilitate yaml marshalling.
	var mdl map[string]any
	err := json.Unmarshal(query.Model, &mdl)
	if err != nil {
		return definitions.AlertQueryExport{}, err
	}
	var queryType *string
	if query.QueryType != "" {
		queryType = &query.QueryType
	}
	return definitions.AlertQueryExport{
		RefID:     query.RefID,
		QueryType: queryType,
		RelativeTimeRange: definitions.RelativeTimeRangeExport{
			FromSeconds: int64(time.Duration(query.RelativeTimeRange.From).Seconds()),
			ToSeconds:   int64(time.Duration(query.RelativeTimeRange.To).Seconds()),
		},
		DatasourceUID: query.DatasourceUID,
		Model:         mdl,
		ModelString:   string(query.Model),
	}, nil
}

// AlertingFileExportFromEmbeddedContactPoints creates a definitions.AlertingFileExport DTO from []definitions.EmbeddedContactPoint.
func AlertingFileExportFromEmbeddedContactPoints(orgID int64, ecps []definitions.EmbeddedContactPoint) (definitions.AlertingFileExport, error) {
	f := definitions.AlertingFileExport{APIVersion: 1}

	cache := make(map[string]*definitions.ContactPointExport)
	contactPoints := make([]*definitions.ContactPointExport, 0)
	for _, ecp := range ecps {
		c, ok := cache[ecp.Name]
		if !ok {
			c = &definitions.ContactPointExport{
				OrgID:     orgID,
				Name:      ecp.Name,
				Receivers: make([]definitions.ReceiverExport, 0),
			}
			cache[ecp.Name] = c
			contactPoints = append(contactPoints, c)
		}

		recv, err := ReceiverExportFromEmbeddedContactPoint(ecp)
		if err != nil {
			return definitions.AlertingFileExport{}, err
		}
		c.Receivers = append(c.Receivers, recv)
	}

	for _, c := range contactPoints {
		f.ContactPoints = append(f.ContactPoints, *c)
	}
	return f, nil
}

// ReceiverExportFromEmbeddedContactPoint creates a definitions.ReceiverExport DTO from definitions.EmbeddedContactPoint.
func ReceiverExportFromEmbeddedContactPoint(contact definitions.EmbeddedContactPoint) (definitions.ReceiverExport, error) {
	raw, err := contact.Settings.MarshalJSON()
	if err != nil {
		return definitions.ReceiverExport{}, err
	}
	return definitions.ReceiverExport{
		UID:                   contact.UID,
		Type:                  contact.Type,
		Settings:              raw,
		DisableResolveMessage: contact.DisableResolveMessage,
	}, nil
}

// AlertingFileExportFromRoute creates a definitions.AlertingFileExport DTO from definitions.Route.
func AlertingFileExportFromRoute(orgID int64, route definitions.Route) (definitions.AlertingFileExport, error) {
	f := definitions.AlertingFileExport{
		APIVersion: 1,
		Policies: []definitions.NotificationPolicyExport{{
			OrgID:  orgID,
			Policy: RouteExportFromRoute(&route),
		}},
	}
	return f, nil
}

// RouteExportFromRoute creates a definitions.RouteExport DTO from definitions.Route.
func RouteExportFromRoute(route *definitions.Route) *definitions.RouteExport {
	export := definitions.RouteExport{
		Receiver:          route.Receiver,
		GroupByStr:        route.GroupByStr,
		Match:             route.Match,
		MatchRE:           route.MatchRE,
		Matchers:          route.Matchers,
		ObjectMatchers:    route.ObjectMatchers,
		MuteTimeIntervals: route.MuteTimeIntervals,
		Continue:          route.Continue,
		GroupWait:         route.GroupWait,
		GroupInterval:     route.GroupInterval,
		RepeatInterval:    route.RepeatInterval,
	}

	if len(route.Routes) > 0 {
		export.Routes = make([]*definitions.RouteExport, 0, len(route.Routes))
		for _, r := range route.Routes {
			export.Routes = append(export.Routes, RouteExportFromRoute(r))
		}
	}

	return &export
}

func FromMigrationState(summary *migmodels.OrgMigrationState) *definitions.OrgMigrationState {
	result := &definitions.OrgMigrationState{
		OrgID: summary.OrgID,
	}
	result.MigratedChannels = FromContactPairs(summary.MigratedChannels)
	result.MigratedDashboards = FromDashboardUpgrades(summary.MigratedDashboards)
	result.Errors = append(result.Errors, summary.Errors...)

	return result
}

func FromContactPairs(pairs []*migmodels.ContactPair) []*definitions.ContactPair {
	result := make([]*definitions.ContactPair, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, FromContactPair(p))
	}
	return result
}

func FromContactPair(pair *migmodels.ContactPair) *definitions.ContactPair {
	return &definitions.ContactPair{
		LegacyChannel:       FromLegacyChannel(pair.LegacyChannel),
		ContactPointUpgrade: FromContactPointUpgrade(pair.ContactPointUpgrade),
		Provisioned:         pair.Provisioned,
		Error:               pair.Error,
	}
}

func FromLegacyChannel(channel *migmodels.LegacyChannel) *definitions.LegacyChannel {
	if channel == nil {
		return nil
	}
	return &definitions.LegacyChannel{
		ID:                    channel.ID,
		UID:                   channel.UID,
		Name:                  channel.Name,
		Type:                  channel.Type,
		SendReminder:          channel.SendReminder,
		DisableResolveMessage: channel.DisableResolveMessage,
		Frequency:             channel.Frequency,
		IsDefault:             channel.IsDefault,
		Modified:              channel.Modified,
	}
}

func FromContactPointUpgrade(contactPoint *migmodels.ContactPointUpgrade) *definitions.ContactPointUpgrade {
	if contactPoint == nil {
		return nil
	}
	return &definitions.ContactPointUpgrade{
		Name:                  contactPoint.Name,
		UID:                   contactPoint.UID,
		Type:                  contactPoint.Type,
		DisableResolveMessage: contactPoint.DisableResolveMessage,
		RouteLabel:            contactPoint.RouteLabel,
		Modified:              contactPoint.Modified,
	}
}

func FromDashboardUpgrades(dus []*migmodels.DashboardUpgrade) []*definitions.DashboardUpgrade {
	result := make([]*definitions.DashboardUpgrade, 0, len(dus))
	for _, du := range dus {
		result = append(result, FromDashboardUpgrade(du))
	}
	return result
}

func FromDashboardUpgrade(du *migmodels.DashboardUpgrade) *definitions.DashboardUpgrade {
	res := &definitions.DashboardUpgrade{
		DashboardID:   du.DashboardID,
		DashboardUID:  du.DashboardUID,
		DashboardName: du.DashboardName,
		FolderUID:     du.FolderUID,
		FolderName:    du.FolderName,
		NewFolderUID:  du.NewFolderUID,
		NewFolderName: du.NewFolderName,
		Provisioned:   du.Provisioned,
		Errors:        du.Errors,
		Warnings:      du.Warnings,
	}

	for _, a := range du.MigratedAlerts {
		res.MigratedAlerts = append(res.MigratedAlerts, FromAlertPair(a))
	}

	return res
}

func FromAlertPair(pair *migmodels.AlertPair) *definitions.AlertPair {
	return &definitions.AlertPair{
		LegacyAlert: FromLegacyAlert(pair.LegacyAlert),
		AlertRule:   FromAlertRuleUpgrade(pair.AlertRule),
		Error:       pair.Error,
	}
}

func FromLegacyAlert(alert *migmodels.LegacyAlert) *definitions.LegacyAlert {
	if alert == nil {
		return nil
	}
	return &definitions.LegacyAlert{
		ID:             alert.ID,
		DashboardID:    alert.DashboardID,
		PanelID:        alert.PanelID,
		Name:           alert.Name,
		Paused:         alert.Paused,
		Silenced:       alert.Silenced,
		ExecutionError: alert.ExecutionError,
		Frequency:      alert.Frequency,
		For:            alert.For,
		Modified:       alert.Modified,
	}
}

func FromAlertRuleUpgrade(rule *migmodels.AlertRuleUpgrade) *definitions.AlertRuleUpgrade {
	if rule == nil {
		return nil
	}
	return &definitions.AlertRuleUpgrade{
		UID:          rule.UID,
		Title:        rule.Title,
		DashboardUID: rule.DashboardUID,
		PanelID:      rule.PanelID,
		NoDataState:  definitions.NoDataState(rule.NoDataState),
		ExecErrState: definitions.ExecutionErrorState(rule.ExecErrState),
		For:          rule.For,
		Annotations:  rule.Annotations,
		Labels:       rule.Labels,
		IsPaused:     rule.IsPaused,
		Modified:     rule.Modified,
	}
}

func FromOrgMigrationSummary(summary migmodels.OrgMigrationSummary) definitions.OrgMigrationSummary {
	return definitions.OrgMigrationSummary{
		NewDashboards: summary.NewDashboards,
		NewAlerts:     summary.NewAlerts,
		NewChannels:   summary.NewChannels,
		Removed:       summary.Removed,
		HasErrors:     summary.HasErrors,
	}
}
