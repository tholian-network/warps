package types

import "strconv"
import "testing"
import "time"

func TestDatetime(t *testing.T) {

	t.Run("IsAfter()", func(t *testing.T) {

		datetime1 := ToDatetime("2022-12-31T23:59:59Z")
		datetime2 := ToDatetime("2023-01-01T00:00:01Z")

		if datetime1.IsAfter(datetime2) != false {
			t.Errorf("Expected %s to be after %s", datetime2.String(), datetime1.String())
		}

		if datetime2.IsAfter(datetime1) != true {
			t.Errorf("Expected %s to be after %s", datetime2.String(), datetime1.String())
		}

	})

	t.Run("IsBefore()", func(t *testing.T) {

		datetime1 := ToDatetime("2022-12-31T23:59:59Z")
		datetime2 := ToDatetime("2023-01-01T00:00:01Z")

		if datetime1.IsBefore(datetime2) != true {
			t.Errorf("Expected %s to be before %s", datetime1.String(), datetime2.String())
		}

		if datetime2.IsBefore(datetime1) != false {
			t.Errorf("Expected %s to be before %s", datetime1.String(), datetime2.String())
		}

	})

	t.Run("IsFuture()", func(t *testing.T) {

		past := time.Now().Local().Add(-1 * time.Second)
		future := time.Now().Local().Add(1 * time.Second)

		datetime1 := ToDatetime(past.Format(time.RFC3339))
		datetime2 := ToDatetime(future.Format(time.RFC3339))

		if datetime1.IsFuture() != false {
			t.Errorf("Expected %s to be in the past", datetime1.String())
		}

		if datetime2.IsFuture() != true {
			t.Errorf("Expected %s to be in the future", datetime2.String())
		}

	})

	t.Run("IsPast()", func(t *testing.T) {

		past := time.Now().Local().Add(-1 * time.Second)
		future := time.Now().Local().Add(1 * time.Second)

		datetime1 := ToDatetime(past.Format(time.RFC3339))
		datetime2 := ToDatetime(future.Format(time.RFC3339))

		if datetime1.IsPast() != true {
			t.Errorf("Expected %s to be in the past", datetime1.String())
		}

		if datetime2.IsPast() != false {
			t.Errorf("Expected %s to be in the future", datetime2.String())
		}

	})

	t.Run("Offset()", func(t *testing.T) {

		datetime1 := ToDatetime("2022-12-31T23:59:59Z")
		datetime2 := ToDatetime("2022-12-31T23:59:59Z")
		datetime3 := ToDatetime("2023-01-01T00:00:01Z")
		datetime4 := ToDatetime("2023-01-01T00:00:01Z")

		datetime1.Offset("-12:34")
		datetime2.Offset("+12:34")
		datetime3.Offset("-12:34")
		datetime4.Offset("+12:34")

		want1 := "2023-01-01 12:33:59"
		want2 := "2022-12-31 11:25:59"
		want3 := "2023-01-01 12:34:01"
		want4 := "2022-12-31 11:26:01"

		if datetime1.String() != want1 {
			t.Errorf("Expected %s to be %s", datetime1.String(), want1)
		}

		if datetime2.String() != want2 {
			t.Errorf("Expected %s to be %s", datetime2.String(), want2)
		}

		if datetime3.String() != want3 {
			t.Errorf("Expected %s to be %s", datetime3.String(), want3)
		}

		if datetime4.String() != want4 {
			t.Errorf("Expected %s to be %s", datetime4.String(), want4)
		}

	})

	t.Run("Parse(en_US local format)", func(t *testing.T) {

		current_year, _, _ := time.Now().Date()

		datetime1 := ToDatetime("Sun Jan 01 12:34:56 2023")
		datetime2 := ToDatetime("Thu Feb 29 23:59:59 2024")
		datetime3 := ToDatetime("Jan 01 12:34:56")
		datetime4 := ToDatetime("Feb 29 23:59:59")

		want1 := ToDatetime("2023-01-01 12:34:56")
		want1.ToZulu()

		want2 := ToDatetime("2024-02-29 23:59:59")
		want2.ToZulu()

		want3 := ToDatetime(strconv.FormatUint(uint64(current_year), 10) + "-01-01 12:34:56")
		want3.ToZulu()

		want4 := ToDatetime(strconv.FormatUint(uint64(current_year), 10) + "-02-29 23:59:59")
		want4.ToZulu()

		if datetime1.String() != want1.String() {
			t.Errorf("Expected %s to be %s", datetime1.String(), want1.String())
		}

		if datetime2.String() != want2.String() {
			t.Errorf("Expected %s to be %s", datetime2.String(), want2.String())
		}

		if datetime3.String() != want3.String() {
			t.Errorf("Expected %s to be %s", datetime3.String(), want3.String())
		}

		if datetime4.String() != want4.String() {
			t.Errorf("Expected %s to be %s", datetime4.String(), want4.String())
		}

	})

	t.Run("Parse(en_US meridiem format)", func(t *testing.T) {

		datetime1 := ToDatetime("Sat Dec 31 01:02:03 AM 2022")
		datetime2 := ToDatetime("Sat Dec 31 01:02:03 PM 2022")

		want1 := ToDatetime("2022-12-31 01:02:03")
		want1.ToZulu()

		want2 := ToDatetime("2022-12-31 13:02:03")
		want2.ToZulu()

		if datetime1.String() != want1.String() {
			t.Errorf("Expected %s to be %s", datetime1.String(), want1.String())
		}

		if datetime2.String() != want2.String() {
			t.Errorf("Expected %s to be %s", datetime2.String(), want2.String())
		}

	})

	t.Run("Parse(en_US timezone format)", func(t *testing.T) {

		datetime1 := ToDatetime("Sat Dec 31 01:02:03 AM CET 2022")
		datetime2 := ToDatetime("Sat Dec 31 01:02:03 PM CET 2022")
		datetime3 := ToDatetime("Sat Dec 31 01:02:03 PM CET 2022")
		datetime4 := ToDatetime("Sat Dec 31 12:13:14 CET 2022")

		want1 := ToDatetime("2022-12-31 01:02:03")
		want1.ToZulu()

		want2 := ToDatetime("2022-12-31 13:02:03")
		want2.ToZulu()

		want3 := ToDatetime("2022-12-31 13:02:03")
		want3.ToZulu()

		want4 := ToDatetime("2022-12-31 12:13:14")
		want4.ToZulu()

		if datetime1.String() != want1.String() {
			t.Errorf("Expected %s to be %s", datetime1.String(), want1.String())
		}

		if datetime2.String() != want2.String() {
			t.Errorf("Expected %s to be %s", datetime2.String(), want2.String())
		}

		if datetime3.String() != want3.String() {
			t.Errorf("Expected %s to be %s", datetime3.String(), want3.String())
		}

		if datetime4.String() != want4.String() {
			t.Errorf("Expected %s to be %s", datetime4.String(), want4.String())
		}

	})

	t.Run("Parse(ISO8601 zulu format)", func(t *testing.T) {

		datetime1 := ToDatetime("2022-12-31 01:02:03")
		datetime2 := ToDatetime("2022-12-31T01:02:03Z")
		datetime3 := ToDatetime("2022-12-31T13:14:15.123Z")

		want1 := "2022-12-31 01:02:03"
		want2 := "2022-12-31 01:02:03"
		want3 := "2022-12-31 13:14:15"

		if datetime1.String() != want1 {
			t.Errorf("Expected %s to be %s", datetime1.String(), want1)
		}

		if datetime2.String() != want2 {
			t.Errorf("Expected %s to be %s", datetime2.String(), want2)
		}

		if datetime3.String() != want3 {
			t.Errorf("Expected %s to be %s", datetime3.String(), want3)
		}

	})

	t.Run("Parse(ISO8601 juliet format)", func(t *testing.T) {

		datetime1 := ToDatetime("2022-12-31T01:02:03+12:30")
		datetime2 := ToDatetime("2022-12-31T13:14:15.123+12:30")
		datetime3 := ToDatetime("2022-12-31T01:02:03-03:45")
		datetime4 := ToDatetime("2022-12-31T13:14:15.123-03:45")

		want1 := "2022-12-31 01:02:03"
		want2 := "2022-12-31 13:14:15"
		want3 := "2022-12-31 01:02:03"
		want4 := "2022-12-31 13:14:15"

		if datetime1.String() != want1 {
			t.Errorf("Expected %s to be %s", datetime1.String(), want1)
		}

		if datetime2.String() != want2 {
			t.Errorf("Expected %s to be %s", datetime2.String(), want2)
		}

		if datetime3.String() != want3 {
			t.Errorf("Expected %s to be %s", datetime3.String(), want3)
		}

		if datetime4.String() != want4 {
			t.Errorf("Expected %s to be %s", datetime4.String(), want4)
		}

	})

	t.Run("ToDays()", func(t *testing.T) {

		datetime1 := ToDatetime("2021-01-01T23:59:59Z")
		datetime2 := ToDatetime("2022-02-14T12:34:56Z")
		datetime3 := ToDatetime("2023-02-01T00:00:01Z")
		datetime4 := ToDatetime("2024-02-14T12:34:56Z")

		if datetime1.ToDays() != 31 {
			t.Errorf("Expected %d to be %d", datetime1.ToDays(), 31)
		}

		if datetime2.ToDays() != 28 {
			t.Errorf("Expected %d to be %d", datetime2.ToDays(), 28)
		}

		if datetime3.ToDays() != 28 {
			t.Errorf("Expected %d to be %d", datetime3.ToDays(), 28)
		}

		if datetime4.ToDays() != 29 {
			t.Errorf("Expected %d to be %d", datetime4.ToDays(), 29)
		}

	})

	t.Run("ToDatetimeDifference()", func(t *testing.T) {

		datetime11 := ToDatetime("2021-01-10 12:01:01")
		datetime12 := ToDatetime("2024-03-11 11:00:00")

		datetime21 := ToDatetime("2024-01-01 12:01:01")
		datetime22 := ToDatetime("2024-03-10 12:01:01")

		datetime31 := ToDatetime("2024-01-10 12:01:01")
		datetime32 := ToDatetime("2024-01-13 11:00:00")

		datetime41 := ToDatetime("2023-03-13 12:01:01")
		datetime42 := ToDatetime("2024-01-11 11:00:00")

		diff1 := datetime11.ToDatetimeDifference(datetime12)
		want1 := "0003-02-00 22:59:59"

		diff2 := datetime21.ToDatetimeDifference(datetime22)
		want2 := "0000-02-09 00:00:00"

		diff3 := datetime31.ToDatetimeDifference(datetime32)
		want3 := "0000-00-02 22:59:59"

		diff4 := datetime41.ToDatetimeDifference(datetime42)
		want4 := "0000-09-28 22:59:59"

		if diff1.String() != want1 {
			t.Errorf("Expected %s to be %s", diff1.String(), want1)
		}

		if diff2.String() != want2 {
			t.Errorf("Expected %s to be %s", diff2.String(), want2)
		}

		if diff3.String() != want3 {
			t.Errorf("Expected %s to be %s", diff3.String(), want3)
		}

		if diff4.String() != want4 {
			t.Errorf("Expected %s to be %s", diff4.String(), want4)
		}

	})

	t.Run("ToTimeDifference()", func(t *testing.T) {

		datetime11 := ToDatetime("2021-01-10 12:01:01")
		datetime12 := ToDatetime("2024-03-11 11:00:00")

		datetime21 := ToDatetime("2024-01-01 12:01:01")
		datetime22 := ToDatetime("2024-03-10 12:01:01")

		datetime31 := ToDatetime("2024-01-10 12:01:01")
		datetime32 := ToDatetime("2024-01-13 11:00:00")

		datetime41 := ToDatetime("2024-01-10 12:01:01")
		datetime42 := ToDatetime("2024-01-11 11:00:00")

		diff1 := datetime11.ToTimeDifference(datetime12)
		want1 := "27742:59:59"

		diff2 := datetime21.ToTimeDifference(datetime22)
		want2 := "1656:00:00"

		diff3 := datetime31.ToTimeDifference(datetime32)
		want3 := "70:59:59"

		diff4 := datetime41.ToTimeDifference(datetime42)
		want4 := "22:59:59"

		if diff1.String() != want1 {
			t.Errorf("Expected %s to be %s", diff1.String(), want1)
		}

		if diff2.String() != want2 {
			t.Errorf("Expected %s to be %s", diff2.String(), want2)
		}

		if diff3.String() != want3 {
			t.Errorf("Expected %s to be %s", diff3.String(), want3)
		}

		if diff4.String() != want4 {
			t.Errorf("Expected %s to be %s", diff4.String(), want4)
		}

	})

	t.Run("ToWeekday()", func(t *testing.T) {

		datetime1 := ToDatetime("2021-01-01T23:59:59Z")
		datetime2 := ToDatetime("2022-02-14T12:34:56Z")
		datetime3 := ToDatetime("2023-02-01T00:00:01Z")
		datetime4 := ToDatetime("2024-02-14T12:34:56Z")

		if datetime1.ToWeekday() != "Friday" {
			t.Errorf("Expected %s to be %s", datetime1.String(), "Friday")
		}

		if datetime2.ToWeekday() != "Monday" {
			t.Errorf("Expected %s to be %s", datetime2.String(), "Monday")
		}

		if datetime3.ToWeekday() != "Wednesday" {
			t.Errorf("Expected %s to be %s", datetime3.String(), "Wednesday")
		}

		if datetime4.ToWeekday() != "Wednesday" {
			t.Errorf("Expected %s to be %s", datetime4.String(), "Wednesday")
		}

	})

}
