package validation

import (
	"testing"

	ngmodels "github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/stretchr/testify/require"
)

func TestValidateRule(t *testing.T) {
	gen := ngmodels.RuleGen.With(ngmodels.RuleGen.WithRandomRecordingRules())

	t.Run("validates metric name on recording rules", func(t *testing.T) {
		rules := gen.GenerateMany()
		for _, rule := range rules {
			if rule.Type() == ngmodels.RuleTypeRecording {
				rule.Record.Metric = "invalid metric name"
				_, err := ValidateRule(rule)
				require.Error(t, err)
				require.ErrorContains(t, err, "must be a valid Prometheus metric name")
			}
		}
	})

	t.Run("validation also clears ignored fields on recording rules", func(t *testing.T) {
		rules := gen.GenerateMany()
		for _, rule := range rules {
			rule, err := ValidateRule(rule)
			require.NoError(t, err)
			if rule.Type() == ngmodels.RuleTypeRecording {
				require.Empty(t, rule.NoDataState)
				require.Empty(t, rule.ExecErrState)
				require.Empty(t, rule.Condition)
				require.Zero(t, rule.For)
				require.Nil(t, rule.NotificationSettings)
			}
		}
	})
}
