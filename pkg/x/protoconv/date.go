// Copyright 2025 Takekazu Omi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protoconv

import (
	"time"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DateNew は、year, month, day から date.Date を生成。
func DateNew(year, month, day int32) *date.Date {
	return &date.Date{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

// DateToTime は、date.Date を time.Time に変換。
func DateToTime(d *date.Date) (time.Time, error) {
	if d == nil {
		return time.Time{}, nil
	}
	return time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.Local), nil
}

// TimeToDate は、time.Time を date.Date に変換。
func TimeToDate(t time.Time) *date.Date {
	return &date.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}
}

// DateToTimestamp は、date.Date を timestamppb.Timestamp に変換。
func DateToTimestamp(d *date.Date) *timestamppb.Timestamp {
	if d == nil {
		return nil
	}
	return timestamppb.New(time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.Local))
}

// TimestampToDate は、timestamppb.Timestamp を date.Date に変換。
func TimestampToDate(ts *timestamppb.Timestamp) *date.Date {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &date.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}
}
