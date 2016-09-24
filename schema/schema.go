package schema

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Luzifer/go_helpers/str"
	couch "github.com/lancecarlson/couchgo"
)

const (
	TagIll        = "ill"
	TagVacation   = "vacation"
	TagWeekend    = "weekend"
	TagHoliday    = "holiday"
	TagEvent      = "event"
	TagHomeoffice = "homeoffice"
	TagBreak      = "break"
	TagAutotrack  = "autotrack"
	TagOnCall     = "on-call"
)

func evalTags(tags []string, tag string) []string {
	rawTag := strings.TrimLeft(tag, "+-")
	out := tags

	switch tag[0] {
	default:
		fallthrough
	case '+':
		if !str.StringInSlice(rawTag, tags) {
			out = append(tags, rawTag)
		}
	case '-':
		if str.StringInSlice(rawTag, tags) {
			out = []string{}
			for _, t := range tags {
				if t != rawTag {
					out = append(out, t)
				}
			}
		}
	}

	sort.Strings(out)

	return out
}

type Day struct {
	DayID    string `json:"_id"`
	Revision string `json:"_rev,omitempty"`

	Times []*Time  `json:"times"`
	Tags  []string `json:"tags,omitempty"`

	// deprecated tags, have auto-migration
	IsIll      bool `json:"is_ill,omitempty"`
	IsVacation bool `json:"is_vacation,omitempty"`
	IsWeekend  bool `json:"is_weekend,omitempty"`
	IsHoliday  bool `json:"is_holiday,omitempty"`
	IsEvent    bool `json:"is_event,omitempty"`
	Homeoffice bool `json:"homeoffice,omitempty"`
}

func (d *Day) Tag(tag string) {
	d.Tags = evalTags(d.Tags, tag)
}

func (d *Day) migrate() {
	if d.IsIll {
		d.Tag(TagIll)
	}
	if d.IsVacation {
		d.Tag(TagVacation)
	}
	if d.IsWeekend {
		d.Tag(TagWeekend)
	}
	if d.IsHoliday {
		d.Tag(TagHoliday)
	}
	if d.IsEvent {
		d.Tag(TagEvent)
	}
	if d.Homeoffice {
		d.Tag(TagHomeoffice)
	}
	d.IsIll = false
	d.IsEvent = false
	d.IsVacation = false
	d.IsWeekend = false
	d.IsHoliday = false
	d.Homeoffice = false

	for _, t := range d.Times {
		t.migrate()
	}
}

func (d *Day) validate() error {

	for _, t := range d.Times {
		if err := t.validate(); err != nil {
			return err
		}
	}

	return nil
}

func LoadDay(db *couch.Client, date time.Time, mayCreate bool) (*Day, error) {
	id := date.Format("2006-01-02")
	doc := &Day{}
	if err := db.Get(id, doc); err != nil {
		if strings.Contains(err.Error(), "not_found") && mayCreate {
			doc = &Day{DayID: id, Times: []*Time{}}
		} else {
			return nil, err
		}
	}
	doc.migrate()
	return doc, nil
}

func (d *Day) Save(db *couch.Client) error {
	if err := d.validate(); err != nil {
		return err
	}

	res, err := db.Save(d)
	if err != nil {
		return err
	}
	d.Revision = res.Rev
	return nil
}

type Time struct {
	ID    string `json:"id"`
	Start string `json:"start"`
	End   string `json:"end"`

	Tags []string `json:"tags,omitempty"`

	// deprecated tags, have auto-migration
	IsBreak     bool `json:"is_break,omitempty"`
	IsAutotrack bool `json:"is_autotrack,omitempty"`
	IsOnCall    bool `json:"is_on_call,omitempty"`
}

func (w *Time) Tag(tag string) {
	w.Tags = evalTags(w.Tags, tag)
}

func (w *Time) migrate() {
	if w.IsBreak {
		w.Tag(TagBreak)
	}
	if w.IsAutotrack {
		w.Tag(TagAutotrack)
	}
	if w.IsOnCall {
		w.Tag(TagOnCall)
	}
	w.IsBreak = false
	w.IsAutotrack = false
	w.IsOnCall = false
}

func (t *Time) validate() error {
	now := time.Now().Format("15:04:05")
	if t.Start == "now" {
		t.Start = now
	}
	if t.End == "now" {
		t.End = now
	}

	if _, err := time.Parse("15:04:05", t.Start); err != nil {
		return fmt.Errorf("Time %.7s has invalid start date: %s", t.ID, err)
	}

	if _, err := time.Parse("15:04:05", t.End); err != nil {
		return fmt.Errorf("Time %.7s has invalid end date: %s", t.ID, err)
	}

	return nil
}

type Overtime struct {
	Value float64 `json:"value"`
}

func GetOvertime(db *couch.Client, day time.Time) (Overtime, error) {
	var opts *url.Values

	if !day.IsZero() {
		opts = &url.Values{}
		opts.Set("reduce", "false")
		opts.Set("startkey", fmt.Sprintf("\"%s\"", day.Format("2006-01-02")))
		opts.Set("endkey", fmt.Sprintf("\"%s\"", day.Format("2006-01-02")))
	}

	mdr, err := db.View("analysis", "overtime", opts, nil)
	if err != nil {
		return Overtime{}, err
	}

	if len(mdr.Rows) == 0 {
		return Overtime{}, fmt.Errorf("Did not find any results in view")
	}

	result := Overtime{}

	return result, couch.Remarshal(mdr.Rows[0], &result)
}
