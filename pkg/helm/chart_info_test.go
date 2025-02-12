package helm_test

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ChartParsingTestSuite struct {
	suite.Suite
}

func (*ChartParsingTestSuite) SetupTest() {
	viper.Set("values-file", "values.yaml")
}

func TestChartParsingTestSuite(t *testing.T) {
	suite.Run(t, new(ChartParsingTestSuite))
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOff() {
	chartPath := filepath.Join("test-fixtures", "full-template")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: false,
	})
	suite.NoError(err)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOn() {
	chartPath := filepath.Join("test-fixtures", "full-template")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
	})
	expectedError := `values without documentation: 
controller
controller.name
controller.image
controller.image.repository
controller.image.tag
controller.extraVolumes
controller.extraVolumes.[0].name
controller.extraVolumes.[0].configMap
controller.extraVolumes.[0].configMap.name
controller.publishService
controller.service
controller.service.annotations
controller.service.annotations.external-dns.alpha.kubernetes.io/hostname
controller.service.type`
	suite.EqualError(err, expectedError)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOnIgnores() {
	chartPath := filepath.Join("test-fixtures", "full-template")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
		AllowedMissingValuePaths: []string{
			"controller",
			"controller.image",
			"controller.name",
			"controller.image.repository",
			"controller.image.tag",
			"controller.extraVolumes",
			"controller.extraVolumes.[0].name",
			"controller.extraVolumes.[0].configMap",
			"controller.extraVolumes.[0].configMap.name",
			"controller.publishService",
			"controller.service",
			"controller.service.annotations",
			"controller.service.annotations.external-dns.alpha.kubernetes.io/hostname",
			"controller.service.type",
		},
	})
	suite.NoError(err)
}

func (suite *ChartParsingTestSuite) TestNotFullyDocumentedChartStrictModeOnIgnoresRegexp() {
	chartPath := filepath.Join("test-fixtures", "full-template")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
		AllowedMissingValueRegexps: []*regexp.Regexp{
			regexp.MustCompile("controller.*"),
		},
	})
	suite.NoError(err)
}

func (suite *ChartParsingTestSuite) TestFullyDocumentedChartStrictModeOn() {
	chartPath := filepath.Join("test-fixtures", "fully-documented")
	_, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{
		StrictMode: true,
	})
	suite.NoError(err)
}

func (suite *ChartParsingTestSuite) TestFullFullTemplateChartLock() {
	chartPath := filepath.Join("test-fixtures", "full-template")

	want := helm.ChartRequirements{
		Dependencies: []helm.ChartRequirementsItem{
			{Name: "nginx-ingress", Version: "~2.0.0", Repository: "oci://ghcr.io/nginx/charts", Alias: "", LockedVersion: "2.0.1"},
			{Name: "nginx-ingress", Version: "~1.0.0", Repository: "oci://ghcr.io/nginx/charts", Alias: "nginxA", LockedVersion: "1.0.2"},
		},
	}
	info, err := helm.ParseChartInformation(chartPath, helm.ChartValuesDocumentationParsingConfig{})
	suite.NoError(err)

	suite.Equal(want, info.ChartRequirements)
}
