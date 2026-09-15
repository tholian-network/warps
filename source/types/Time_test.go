package types

import "testing"
import "time"

func TestTime(t *testing.T) {

	t.Run("AddHour()", func(t *testing.T) {

		time1 := ToTime("11:59:59")
		time2 := ToTime("23:59:59")

		time1.AddHour()
		time2.AddHour()

		if time1.String() != "12:59:59" {
			t.Errorf("Expected %s to be %s", time1.String(), "12:59:59")
		}

		if time2.String() != "24:59:59" {
			t.Errorf("Expected %s to be %s", time2.String(), "24:59:59")
		}

	})

	t.Run("AddMinute()", func(t *testing.T) {

		time1 := ToTime("11:59:59")
		time2 := ToTime("23:59:59")

		time1.AddMinute()
		time2.AddMinute()

		if time1.String() != "12:00:59" {
			t.Errorf("Expected %s to be %s", time1.String(), "12:00:59")
		}

		if time2.String() != "24:00:59" {
			t.Errorf("Expected %s to be %s", time2.String(), "24:00:59")
		}

	})

	t.Run("AddSecond()", func(t *testing.T) {

		time1 := ToTime("11:59:59")
		time2 := ToTime("23:59:59")

		time1.AddSecond()
		time2.AddSecond()

		if time1.String() != "12:00:00" {
			t.Errorf("Expected %s to be %s", time1.String(), "12:00:00")
		}

		if time2.String() != "24:00:00" {
			t.Errorf("Expected %s to be %s", time2.String(), "24:00:00")
		}

	})

	t.Run("AddTime()", func(t *testing.T) {

		time1 := ToTime("11:59:59")
		time2 := ToTime("11:59:59")
		time3 := ToTime("11:59:59")
		time4 := ToTime("23:59:59")
		time5 := ToTime("23:59:59")
		time6 := ToTime("23:59:59")

		time1.AddTime(ToTime("12:00:01"))
		time2.AddTime(ToTime("24:00:01"))
		time3.AddTime(ToTime("48:00:01"))

		time4.AddTime(ToTime("12:00:01"))
		time5.AddTime(ToTime("24:00:01"))
		time6.AddTime(ToTime("48:00:01"))

		if time1.String() != "24:00:00" {
			t.Errorf("Expected %s to be %s", time1.String(), "24:00:00")
		}

		if time2.String() != "36:00:00" {
			t.Errorf("Expected %s to be %s", time2.String(), "36:00:00")
		}

		if time3.String() != "60:00:00" {
			t.Errorf("Expected %s to be %s", time3.String(), "60:00:00")
		}

		if time4.String() != "36:00:00" {
			t.Errorf("Expected %s to be %s", time4.String(), "36:00:00")
		}

		if time5.String() != "48:00:00" {
			t.Errorf("Expected %s to be %s", time5.String(), "48:00:00")
		}

		if time6.String() != "72:00:00" {
			t.Errorf("Expected %s to be %s", time6.String(), "72:00:00")
		}

	})

	t.Run("IsAfter()", func(t *testing.T) {

		time1 := ToTime("11:59:59")
		time2 := ToTime("12:00:00")
		time3 := ToTime("12:00:01")
		time4 := ToTime("23:59:59")
		time5 := ToTime("24:00:00")
		time6 := ToTime("24:01:00")

		if time1.IsAfter(time2) != false {
			t.Errorf("Expected %s to be not after %s", time1.String(), time2.String())
		}

		if time1.IsAfter(time3) != false {
			t.Errorf("Expected %s to be not after %s", time1.String(), time3.String())
		}

		if time1.IsAfter(time4) != false {
			t.Errorf("Expected %s to be not after %s", time1.String(), time4.String())
		}

		if time1.IsAfter(time5) != false {
			t.Errorf("Expected %s to be not after %s", time1.String(), time5.String())
		}

		if time1.IsAfter(time6) != false {
			t.Errorf("Expected %s to be not after %s", time1.String(), time6.String())
		}

		if time2.IsAfter(time1) != true {
			t.Errorf("Expected %s to be after %s", time2.String(), time1.String())
		}

		if time2.IsAfter(time3) != false {
			t.Errorf("Expected %s to be not after %s", time2.String(), time3.String())
		}

		if time2.IsAfter(time4) != false {
			t.Errorf("Expected %s to be not after %s", time2.String(), time4.String())
		}

		if time2.IsAfter(time5) != false {
			t.Errorf("Expected %s to be not after %s", time2.String(), time5.String())
		}

		if time2.IsAfter(time6) != false {
			t.Errorf("Expected %s to be not after %s", time2.String(), time6.String())
		}

		if time3.IsAfter(time1) != true {
			t.Errorf("Expected %s to be after %s", time3.String(), time1.String())
		}

		if time3.IsAfter(time2) != true {
			t.Errorf("Expected %s to be after %s", time3.String(), time2.String())
		}

		if time3.IsAfter(time4) != false {
			t.Errorf("Expected %s to be not after %s", time3.String(), time4.String())
		}

		if time3.IsAfter(time5) != false {
			t.Errorf("Expected %s to be not after %s", time3.String(), time5.String())
		}

		if time3.IsAfter(time6) != false {
			t.Errorf("Expected %s to be not after %s", time3.String(), time6.String())
		}

		if time4.IsAfter(time1) != true {
			t.Errorf("Expected %s to be after %s", time4.String(), time1.String())
		}

		if time4.IsAfter(time2) != true {
			t.Errorf("Expected %s to be after %s", time4.String(), time2.String())
		}

		if time4.IsAfter(time3) != true {
			t.Errorf("Expected %s to be after %s", time4.String(), time3.String())
		}

		if time4.IsAfter(time5) != false {
			t.Errorf("Expected %s to be not after %s", time4.String(), time5.String())
		}

		if time4.IsAfter(time6) != false {
			t.Errorf("Expected %s to be not after %s", time4.String(), time6.String())
		}

		if time5.IsAfter(time1) != true {
			t.Errorf("Expected %s to be after %s", time5.String(), time1.String())
		}

		if time5.IsAfter(time2) != true {
			t.Errorf("Expected %s to be after %s", time5.String(), time2.String())
		}

		if time5.IsAfter(time3) != true {
			t.Errorf("Expected %s to be after %s", time5.String(), time3.String())
		}

		if time5.IsAfter(time4) != true {
			t.Errorf("Expected %s to be after %s", time5.String(), time4.String())
		}

		if time5.IsAfter(time6) != false {
			t.Errorf("Expected %s to be not after %s", time5.String(), time6.String())
		}

		if time6.IsAfter(time1) != true {
			t.Errorf("Expected %s to be after %s", time6.String(), time1.String())
		}

		if time6.IsAfter(time2) != true {
			t.Errorf("Expected %s to be after %s", time6.String(), time2.String())
		}

		if time6.IsAfter(time3) != true {
			t.Errorf("Expected %s to be after %s", time6.String(), time3.String())
		}

		if time6.IsAfter(time4) != true {
			t.Errorf("Expected %s to be after %s", time6.String(), time4.String())
		}

		if time6.IsAfter(time5) != true {
			t.Errorf("Expected %s to be after %s", time6.String(), time5.String())
		}

	})

	t.Run("IsBefore()", func(t *testing.T) {

		time1 := ToTime("11:59:59")
		time2 := ToTime("12:00:00")
		time3 := ToTime("12:00:01")
		time4 := ToTime("23:59:59")
		time5 := ToTime("24:00:00")
		time6 := ToTime("24:01:00")

		if time1.IsBefore(time2) != true {
			t.Errorf("Expected %s to be before %s", time1.String(), time2.String())
		}

		if time1.IsBefore(time3) != true {
			t.Errorf("Expected %s to be before %s", time1.String(), time3.String())
		}

		if time1.IsBefore(time4) != true {
			t.Errorf("Expected %s to be before %s", time1.String(), time4.String())
		}

		if time1.IsBefore(time5) != true {
			t.Errorf("Expected %s to be before %s", time1.String(), time5.String())
		}

		if time1.IsBefore(time6) != true {
			t.Errorf("Expected %s to be before %s", time1.String(), time6.String())
		}

		if time2.IsBefore(time1) != false {
			t.Errorf("Expected %s to be not before %s", time2.String(), time1.String())
		}

		if time2.IsBefore(time3) != true {
			t.Errorf("Expected %s to be before %s", time2.String(), time3.String())
		}

		if time2.IsBefore(time4) != true {
			t.Errorf("Expected %s to be before %s", time2.String(), time4.String())
		}

		if time2.IsBefore(time5) != true {
			t.Errorf("Expected %s to be before %s", time2.String(), time5.String())
		}

		if time2.IsBefore(time6) != true {
			t.Errorf("Expected %s to be before %s", time2.String(), time6.String())
		}

		if time3.IsBefore(time1) != false {
			t.Errorf("Expected %s to be not before %s", time3.String(), time1.String())
		}

		if time3.IsBefore(time2) != false {
			t.Errorf("Expected %s to be not before %s", time3.String(), time2.String())
		}

		if time3.IsBefore(time4) != true {
			t.Errorf("Expected %s to be before %s", time3.String(), time4.String())
		}

		if time3.IsBefore(time5) != true {
			t.Errorf("Expected %s to be before %s", time3.String(), time5.String())
		}

		if time3.IsBefore(time6) != true {
			t.Errorf("Expected %s to be before %s", time3.String(), time6.String())
		}

		if time4.IsBefore(time1) != false {
			t.Errorf("Expected %s to be not before %s", time4.String(), time1.String())
		}

		if time4.IsBefore(time2) != false {
			t.Errorf("Expected %s to be not before %s", time4.String(), time2.String())
		}

		if time4.IsBefore(time3) != false {
			t.Errorf("Expected %s to be not before %s", time4.String(), time3.String())
		}

		if time4.IsBefore(time5) != true {
			t.Errorf("Expected %s to be before %s", time4.String(), time5.String())
		}

		if time4.IsBefore(time6) != true {
			t.Errorf("Expected %s to be before %s", time4.String(), time6.String())
		}

		if time5.IsBefore(time1) != false {
			t.Errorf("Expected %s to be not before %s", time5.String(), time1.String())
		}

		if time5.IsBefore(time2) != false {
			t.Errorf("Expected %s to be not before %s", time5.String(), time2.String())
		}

		if time5.IsBefore(time3) != false {
			t.Errorf("Expected %s to be not before %s", time5.String(), time3.String())
		}

		if time5.IsBefore(time4) != false {
			t.Errorf("Expected %s to be not before %s", time5.String(), time4.String())
		}

		if time5.IsBefore(time6) != true {
			t.Errorf("Expected %s to be before %s", time5.String(), time6.String())
		}

		if time6.IsBefore(time1) != false {
			t.Errorf("Expected %s to be not before %s", time6.String(), time1.String())
		}

		if time6.IsBefore(time2) != false {
			t.Errorf("Expected %s to be not before %s", time6.String(), time2.String())
		}

		if time6.IsBefore(time3) != false {
			t.Errorf("Expected %s to be not before %s", time6.String(), time3.String())
		}

		if time6.IsBefore(time4) != false {
			t.Errorf("Expected %s to be not before %s", time6.String(), time4.String())
		}

		if time6.IsBefore(time5) != false {
			t.Errorf("Expected %s to be not before %s", time6.String(), time5.String())
		}

	})

	t.Run("IsFuture()", func(t *testing.T) {

		past := time.Now().Local().Add(-1 * time.Second)
		future := time.Now().Local().Add(1 * time.Second)

		time1 := ToTime(past.Format(time.TimeOnly))
		time2 := ToTime(future.Format(time.TimeOnly))

		if time1.IsFuture() != false {
			t.Errorf("Expected %s to be in the past", time1.String())
		}

		if time2.IsFuture() != true {
			t.Errorf("Expected %s to be in the future", time2.String())
		}

	})

	t.Run("IsPast()", func(t *testing.T) {

		past := time.Now().Local().Add(-1 * time.Second)
		future := time.Now().Local().Add(1 * time.Second)

		time1 := ToTime(past.Format(time.TimeOnly))
		time2 := ToTime(future.Format(time.TimeOnly))

		if time1.IsPast() != true {
			t.Errorf("Expected %s to be in the past", time1.String())
		}

		if time2.IsPast() != false {
			t.Errorf("Expected %s to be in the future", time2.String())
		}

	})

	t.Run("Offset()", func(t *testing.T) {

		time1 := ToTime("23:59:59")
		time2 := ToTime("23:59:59")
		time3 := ToTime("00:00:01")
		time4 := ToTime("00:00:01")

		time1.Offset("-12:34")
		time2.Offset("+12:34")
		time3.Offset("-12:34")
		time4.Offset("+12:34")

		want1 := "12:33:59"
		want2 := "11:25:59"
		want3 := "12:34:01"
		want4 := "11:26:01"

		if time1.String() != want1 {
			t.Errorf("Expected %s to be %s", time1.String(), want1)
		}

		if time2.String() != want2 {
			t.Errorf("Expected %s to be %s", time2.String(), want2)
		}

		if time3.String() != want3 {
			t.Errorf("Expected %s to be %s", time3.String(), want3)
		}

		if time4.String() != want4 {
			t.Errorf("Expected %s to be %s", time4.String(), want4)
		}

	})

	t.Run("Parse(HH:ii:ss format)", func(t *testing.T) {

		time1 := ToTime("12:34:56")
		time2 := ToTime("23:59:59")

		want1 := "12:34:56"
		want2 := "23:59:59"

		if time1.String() != want1 {
			t.Errorf("Expected %s to be %s", time1.String(), want1)
		}

		if time2.String() != want2 {
			t.Errorf("Expected %s to be %s", time2.String(), want2)
		}

	})

	t.Run("Parse(HH:ii format)", func(t *testing.T) {

		time1 := ToTime("01:02")
		time2 := ToTime("23:59")

		want1 := "01:02:00"
		want2 := "23:59:00"

		if time1.String() != want1 {
			t.Errorf("Expected %s to be %s", time1.String(), want1)
		}

		if time2.String() != want2 {
			t.Errorf("Expected %s to be %s", time2.String(), want2)
		}

	})

}
