package xlsx

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-propane-dwnld/config"
	"github.com/pulpfree/gdps-propane-dwnld/graphql"
	"github.com/pulpfree/gdps-propane-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	date             = "2018-05-01"
	defaultsFilePath = "../config/defaults.yaml"
	filePath         = "../tmp/testfile.xlsx"
	timeFormat       = "2006-01-02"
)

// Suite struct
type Suite struct {
	suite.Suite
	cfg     *config.Config
	request *model.Request
	file    *XLSX
	graphql *graphql.Client
}

// SetupTest method
func (suite *Suite) SetupTest() {
	os.Setenv("Stage", "test")
	dte, err := time.Parse(timeFormat, date)
	suite.request = &model.Request{
		Date: dte,
	}
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFilePath}
	err = suite.cfg.Load()
	suite.NoError(err)
	suite.IsType(new(config.Config), suite.cfg)

	suite.file, err = NewFile()
	suite.NoError(err)
	suite.IsType(new(XLSX), suite.file)

	suite.graphql = graphql.New(suite.request, suite.cfg, "")
	suite.IsType(new(graphql.Client), suite.graphql)
}

// TestOutput method
func (suite *Suite) TestOutput() {

	sales, err := suite.graphql.PropaneSales()
	suite.NoError(err)
	suite.IsType(new(model.PropaneSales), sales)

	err = suite.file.PropaneSales(sales)
	suite.NoError(err)

	_, err = suite.file.OutputToDisk(filePath)
	suite.NoError(err)
}

// TestXLSXSuite function
func TestXLSXSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
