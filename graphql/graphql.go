package graphql

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/machinebox/graphql"
	"github.com/pulpfree/gdps-propane-dwnld/config"
	"github.com/pulpfree/gdps-propane-dwnld/model"

	log "github.com/sirupsen/logrus"
)

// Client struct
type Client struct {
	hdrs    http.Header
	client  *graphql.Client
	request *model.Request
}

const timeLongFrmt = "2006-01-02"

// New graphql client
func New(req *model.Request, cfg *config.Config, authToken string) (c *Client) {

	hdrs := http.Header{}
	if len(authToken) > 0 {
		hdrs.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	c = &Client{
		client:  graphql.NewClient(cfg.GraphqlURI),
		hdrs:    hdrs,
		request: req,
	}

	return c
}

// PropaneSales method
func (c *Client) PropaneSales() (rpt *model.PropaneSales, err error) {

	req := graphql.NewRequest(`
    query ($date: String!) {
      propaneReportDwnld(date: $date) {
        date
        deliveries
        sales {
          date
          sales
        }
      }
    }
  `)

	req.Var("date", formattedDate(c.request.Date))
	req.Header = c.hdrs

	ctx := context.Background()
	err = c.client.Run(ctx, req, &rpt)
	if err != nil {
		log.Errorf("error running graphql client: %s", err.Error())
		return nil, err
	}

	return rpt, err
}

// formattedDate function
func formattedDate(date time.Time) string {
	return date.Format(timeLongFrmt)
}
