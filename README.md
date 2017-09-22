# Luzifer / worktime

Worktime is intended as a personal time tracker for cases where the employer trusts you to work the amount of hours set in your contract but you want to make sure you really did.

## Usage

```bash
# worktime help
Manage worktimes in CouchDB

Usage:
  worktime [command]

Available Commands:
  overtime    Shows total overtime over all time
  show        Display a summary of the given / current day
  tag         Adds or removes a tag from the day or time entry inside the day
  time        manipulate times of a day

Flags:
      --config string    config file (default is $HOME/.worktime.yaml)
      --couchdb string   URL to access couchdb (http://user:pass@host:port/database)
```

## Setup

You will need:

- The worktime binary (see [latest release](https://github.com/Luzifer/worktime/releases/latest))
- A [CouchDB](https://couchdb.apache.org/) instance already set up with a database
- One design document (`analysis`) with (at least) one view (`overtime`) in it (see below for an example)

After you've set up all of this you will need to export the `COUCHDB` environment variable in the form `https://user:pass@host/db` when using the `worktime` commandline tool. (Alternatively you could pass the `--couchdb` commandline flag to worktime every time you call it)

### Example of required view

This view (`analysis/overtime`) does the time calculation for you so you definitely need to adjust that one for your personal needs (weekly hours, tags, ...).

**Map Function**

```javascript
function(day) {
  var hasTag = function(day, tag) { return day.tags && day.tags.indexOf(tag) > -1 }

  var worked_time = 0.0;
  var daily_hours = 8.0;

  for (var idx = 0; idx < day.times.length; idx++) {
    var time = day.times[idx];
    var s = new Date(day._id + 'T' + time.start);
    var e = new Date(day._id + 'T' + time.end);
    var diff = (e - s) / 1000.0 / 3600.0;
    
    if (hasTag(time, "break")) {
      worked_time = worked_time - diff;
    } else {
      worked_time = worked_time + diff;
    }
  }
  
  if (hasTag(day, "holiday") || hasTag(day, "vacation") || hasTag(day, "ill") || hasTag(day, "weekend")) {
    required_time = 0.0;
  } else {
    required_time = daily_hours;
  }
  
  var outcome = worked_time - required_time;
  
  emit(day._id, outcome);
}
```

**Reduce Function**

```javascript
function(keys, values, rereduce) {
  var sum = 0.0;
  for(var idx = 0; idx < values.length; idx++) {
    sum = sum + values[idx];
  }
  
  return sum;
}
```

Of course this does not need to be the only view: You can for example add a view to count your days of vacation or even any other evaluation you can imagine. The tool just relies on `analysis/overtime` to be present.
