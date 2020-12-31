package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	ga "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

func main() {
	keyfile := os.Getenv("KEY_FILE")
	viewID := os.Getenv("VIEW_ID")

	data, err := ioutil.ReadFile(keyfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	conf, err := google.JWTConfigFromJSON(data, ga.AnalyticsReadonlyScope)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := oauth2.NoContext
	ts := conf.TokenSource(ctx)
	svc, err := ga.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		fmt.Println(err)
		return
	}

	req := &ga.GetReportsRequest{
		ReportRequests: []*ga.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*ga.DateRange{
					{StartDate: "7daysAgo", EndDate: "today"},
				},
				Metrics: []*ga.Metric{
					{Expression: "ga:uniquePageviews"},
				},
				Dimensions: []*ga.Dimension{
					{Name: "ga:pagePath"},
				},
			},
		},
	}

	res, err := svc.Reports.BatchGet(req).Do()
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.HTTPStatusCode != 200 {
		fmt.Printf("HTTPStatusCode: %d\n", res.HTTPStatusCode)
		return
	}

	for _, report := range res.Reports {
		header := report.ColumnHeader
		dimHdrs := header.Dimensions
		metricHdrs := header.MetricHeader.MetricHeaderEntries
		rows := report.Data.Rows

		if rows == nil {
			fmt.Println("No data")
		}

		for _, row := range rows {
			dims := row.Dimensions
			metrics := row.Metrics

			for i := 0; i < len(dimHdrs) && i < len(dims); i++ {
				fmt.Printf("%s: %s\n", dimHdrs[i], dims[i])
			}

			for _, metric := range metrics {
				for j := 0; j < len(metricHdrs) && j < len(metric.Values); j++ {
					fmt.Printf("%s: %s\n", metricHdrs[j].Name, metric.Values[j])
				}
			}
		}
	}
}
