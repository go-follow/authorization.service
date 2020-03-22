package localtime

//Date - приведение time.Time к часовому поясу Asia/Krasnoyarsk
func Date(d time.Time) time.Time {
	timeZone := "Asia/Krasnoyarsk"
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		logger.Fatalf("unable to set time zone: %s:%v", timeZone, err)
	}
	return d.UTC().In(loc)
}