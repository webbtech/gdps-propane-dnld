package graphql

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-propane-dwnld/config"
	"github.com/pulpfree/gdps-propane-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	date             = "2018-08-01"
	defaultsFilePath = "../config/defaults.yaml"
	timeFormat       = "2006-01-02"
)

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	client  *Client
	cfg     *config.Config
	request *model.Request
}

// SetupTest method
func (suite *UnitSuite) SetupTest() {
	os.Setenv("Stage", "test")
	dte, err := time.Parse(timeFormat, date)
	req := &model.Request{
		Date: dte,
	}
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFilePath}
	err = suite.cfg.Load()
	suite.NoError(err)
	suite.IsType(new(config.Config), suite.cfg)

	suite.client = New(req, suite.cfg, "")
	suite.NoError(err)
	suite.IsType(new(Client), suite.client)
}

// TestSales method
func (suite *UnitSuite) TestSales() {
	res, err := suite.client.PropaneSales()
	suite.NoError(err)
	suite.IsType(new(model.PropaneSales), res)

	record := res.Report.Sales[0]
	suite.NotEqualf("", record.Date, "expect date to be string")
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
