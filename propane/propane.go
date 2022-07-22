package propane

import (
	"errors"
	"path"

	"github.com/pulpfree/gdps-propane-dwnld/awsservices"
	"github.com/pulpfree/gdps-propane-dwnld/config"
	"github.com/pulpfree/gdps-propane-dwnld/graphql"
	"github.com/pulpfree/gdps-propane-dwnld/model"
	"github.com/pulpfree/gdps-propane-dwnld/xlsx"

	log "github.com/sirupsen/logrus"
)

// ReportName constant
const (
	reportFileName = "PropaneReport"
	timeFrmt       = "2006-01"
)

// Report struct
type Report struct {
	authToken string
	cfg       *config.Config
	request   *model.Request
	file      *xlsx.XLSX
	filenm    string
}

// New function
func New(req *model.Request, cfg *config.Config, authToken string) (r *Report, err error) {
	r = &Report{
		authToken: authToken,
		cfg:       cfg,
		request:   req,
	}
	return r, err
}

// Create method
func (r *Report) Create() (err error) {

	r.setFileName()

	// Init graphql and xlsx packages
	client := graphql.New(r.request, r.cfg, r.authToken)
	r.file, err = xlsx.NewFile()
	if err != nil {
		return err
	}

	// Fetch and create Propane Sales
	sales, err := client.PropaneSales()
	if err != nil {
		log.Errorf("Error fetching PropaneSales: %s", err)
		return err
	}
	if len(sales.Report.Sales) <= 0 {
		log.Errorf("Sales fetch from graphql client: %+v", sales.Report.Sales[0])
		return errors.New("Failed to fetch sales")
	}

	err = r.file.PropaneSales(sales)
	if err != nil {
		log.Errorf("failed to create xlsx PropaneSales: %s", err)
		return err
	}

	return err
}

// SaveToDisk method
func (r *Report) SaveToDisk(dir string) (fp string, err error) {

	filePath := path.Join(dir, r.getFileName())
	fp, err = r.file.OutputToDisk(filePath)
	if err != nil {
		log.Errorf("Failed to write to disk: %s", err)
		return "", err
	}
	return fp, err
}

// CreateSignedURL method
func (r *Report) CreateSignedURL() (url string, err error) {

	output, err := r.file.OutputFile()
	if err != nil {
		log.Errorf("failed to create file.OutputFile: %s", err)
		return "", err
	}

	s3Serv, err := awsservices.NewS3(r.cfg)
	if err != nil {
		log.Errorf("failed to create awsservices.NewS3: %s", err)
		return "", err
	}

	return s3Serv.GetSignedURL(r.getFileName(), &output)
}

//
// ======================== Helper Functions =============================== //
//

func (r *Report) setFileName() {
	r.filenm = reportFileName + "_" + r.request.Date.Format(timeFrmt) + ".xlsx"
}

func (r *Report) getFileName() string {
	return r.filenm
}
