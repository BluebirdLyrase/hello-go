package main

import "fmt"

func main() {

	query := ListItem("2024-11-01", "2024-11-30", "037492", "month")
	fmt.Print(query)
}

func ListItem(start_date string, end_date string, asccode string, view_type string) string {

	conditionSelect := ""
	conditionSum := ""
	conditionGroup := ""
	if view_type == "day" {
		conditionSelect = `SELECT FORMAT(x.date_plan,'dd-MM-yyyy')  AS Date,
		ISNULL(s.total,0) AS SALE,
		ISNULL(e.total,0) AS ETA_COMMISSION,
		ISNULL(i.total,0) AS ITEM,
		ISNULL(c.total,0) AS CONVERSION,
		ISNULL(p.total,0) AS PACKAGE,`
		conditionSum = fmt.Sprintf(`DATEDIFF(DAY, '%s', '%s' )+1 AS TOTAL`, start_date, end_date)
		conditionGroup = `ORDER BY x.date_plan`
	} else {
		conditionSelect = `SELECT FORMAT(x.date_plan,'MM-yyyy')  AS Date,
		SUM(ISNULL(s.total,0)) AS SALE,
		SUM(ISNULL(e.total,0)) AS ETA_COMMISSION,
		SUM(ISNULL(i.total,0)) AS ITEM,
		SUM(ISNULL(c.total,0)) AS CONVERSION,
		SUM(ISNULL(p.total,0)) AS PACKAGE,`
		conditionSum = fmt.Sprintf(`DATEDIFF(MONTH, '%s', '%s' )+1 AS TOTAL`, start_date, end_date)
		conditionGroup = `GROUP BY FORMAT(x.date_plan,'MM-yyyy') ORDER BY FORMAT(x.date_plan,'MM-yyyy')`

	}

	query := fmt.Sprintf(`WITH SALE (create_date, total)
	AS (
	SELECT create_date_format as create_date, SUM(product_amt) as total
	FROM performance_report_dtl
	WHERE asc_code = '%s'
	AND create_date_format BETWEEN '%s' AND '%s'
	GROUP BY create_date_format
	),
	CONVERSION (create_date, total)
	AS (
	select a.create_date, count(*) as total  from	(SELECT create_date_format as create_date, (order_no)
	FROM performance_report_dtl
	WHERE asc_code = '%s'
	AND create_date_format BETWEEN '%s' AND '%s'
	GROUP BY create_date_format , order_no) a
	GROUP BY create_date
	),
	PACKAGE (create_date, total)
	AS (
	select a.create_date, count(*) as order_qty  from (
	SELECT create_date_format as create_date, case when (revenue_amt > 0) then 1 else 0 end as revenue_count
	FROM performance_report_dtl
	WHERE asc_code = '%s'
	AND create_date_format BETWEEN '%s' AND '%s' ) a
	GROUP BY create_date
	),
	ETA_COMMISSION (create_date, total)
	AS (
	SELECT create_date_format as create_date, SUM(est_comm_amt) as total
	FROM performance_report_dtl
	WHERE asc_code = '%s' 
	AND create_date_format BETWEEN '%s' AND '%s'
	GROUP BY create_date_format
	),
	ITEM (create_date, total)
	AS (
	SELECT create_date_format as create_date, count(*) as total
	FROM performance_report_dtl
	WHERE asc_code = '%s'
	AND create_date_format BETWEEN '%s' AND '%s'
	GROUP BY create_date_format 
	)
	%s
	%s
	FROM date_plan x
	LEFT JOIN SALE s on s.create_date =x.date_plan
	LEFT JOIN CONVERSION c on c.create_date =x.date_plan
	LEFT JOIN PACKAGE p on p.create_date =x.date_plan
	LEFT JOIN ETA_COMMISSION e on e.create_date =x.date_plan
	LEFT JOIN ITEM i on i.create_date =x.date_plan
	WHERE x.date_plan BETWEEN '%s' AND'%s'
	%s
	`,
		asccode, start_date, end_date,
		asccode, start_date, end_date,
		asccode, start_date, end_date,
		asccode, start_date, end_date,
		asccode, start_date, end_date,
		conditionSelect,
		conditionSum,
		start_date, end_date,
		conditionGroup,
	)
	// println(query)
	return query
}
