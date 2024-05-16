package main

import (
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func main() {
	dates := []string{time.Now().Format("2024-01-02")}
	likes := []int{0}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Discord Likes Progression",
			Subtitle: "GitHub Repository: ProjectDGT",
			Left:     "center",
			TitleStyle: &opts.TextStyle{
				Color:    "#333",
				FontSize: 20,
			},
			SubtitleStyle: &opts.TextStyle{
				Color:    "#aaa",
				FontSize: 14,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Date",
			AxisLabel: &opts.AxisLabel{
				Show:   true,
				Rotate: 45,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Number of Likes",
			AxisLabel: &opts.AxisLabel{
				Show: true,
			},
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
			Left: "right",
			TextStyle: &opts.TextStyle{
				Color: "#333",
			},
		}),
	)

	line.SetXAxis(dates).
		AddSeries("Likes", generateLineItems(likes)).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	f, err := os.Create("line_chart.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := line.Render(f); err != nil {
		panic(err)
	}
}

func generateLineItems(data []int) []opts.LineData {
	items := make([]opts.LineData, len(data))
	for i, v := range data {
		items[i] = opts.LineData{Value: v}
	}
	return items
}
