---
layout: default
title: Utilities
nav_order: 10
description: "utils — string, slice, map, integer, date, conversion and file helpers."
permalink: /utils
last_modified_date: 2026-09-24
---

# Utilities
{: .no_toc }

Small, dependency-free helpers for everyday Go code. Each lives in its own package under
`utils/`.
{: .fs-6 .fw-300 }

1. TOC
{:toc}

---

{: .note }
Package names differ from their directories for historical reasons — import with an alias
that reads well, e.g. `strutil "github.com/HiWay-Media/hwm-go-utils/utils/strings"`.

## Strings — `utils/strings`
Package `strings_utils`.

| Function | Example | Result |
|:--|:--|:--|
| `EncodeURL(raw)` | `https://cdn/x?name=é x&t=a/b` | `https://cdn/x?name=%C3%A9+x&t=a%2Fb` — values escaped once, userinfo and fragment kept, `""` if unparsable |
| `CleanUrlPath(url)` | `https://host//a///b` | `https://host/a/b` |
| `Contains(list, s)` | `["a","b"], "b"` | `true` |
| `ContainsWord(s, word)` | `"Hello World", "world"` | `true` (case-insensitive substring) |
| `RemoveSlice(list, i)` | `["a","b","c"], 1` | `["a","c"]` — returns a copy, input untouched, out-of-range index is a no-op |
| `ReverseString(s)` | `"héllo"` | `"olléh"` (rune-aware) |
| `CountWords(s)` | `"  one two  three "` | `3` |
| `PrettyPrint(v)` | any value | indented JSON string |

## Slices — `utils/slice`
Package `slice_utils`, generic.

```go
slice_utils.Contains([]int{1, 2, 3}, 2)                     // true
slice_utils.RemoveDuplicates([]string{"eu", "us", "eu"})   // ["eu" "us"], order kept
slice_utils.Sum([]int{1, 2, 3})                             // 6
```

## Maps — `utils/map`
Package `map_utils`, generic.

```go
map_utils.MergeMaps(defaults, overrides) // new map, overrides win on duplicate keys
map_utils.InvertMap(map[string]int{"a": 1}) // map[int]string{1: "a"} — values must be unique
```

## Integers — `utils/int`
Package `int_utils`.

```go
int_utils.Min(4, 2, 9)  // 2   (0 when called with no arguments)
int_utils.Max(4, 2, 9)  // 9
int_utils.IsEven(4)     // true
int_utils.IsOdd(4)      // false
```

## Dates — `utils/date`
Package `date_utils`. An empty layout means `2006-01-02`.

```go
d, _ := date_utils.ParseDate("2026-09-24", "")
date_utils.FormatDate(d, "")               // "2026-09-24"
date_utils.AddDays(d, 7)                   // 2026-10-01
date_utils.AddMonths(d, 1)                 // 2026-10-24
date_utils.AddYears(d, 1)                  // 2027-09-24
date_utils.DaysBetween(d, d.AddDate(0, 0, 10)) // 10 (whole days, truncated)
date_utils.IsWeekend(d)                    // false (Thursday)
date_utils.DaysInMonth(2028, time.February) // 29
```

## Conversions — `utils/conv`
Package `conv_utils`.

```go
n, err := conv_utils.StrToInt("42") // 42, nil
conv_utils.IntToStr(42)             // "42"
conv_utils.IntToFloat(3)            // 3.0
conv_utils.FloatToInt(3.9)          // 3 (truncates)
```

## Files — `utils/file`
Package `file_utils`.

```go
if err := file_utils.WriteToFile(report, "/tmp/report.json"); err != nil { // create or truncate
	return err
}
file_utils.FileExists("/tmp/report.json") // true — false for directories and unreadable paths
```
