# code snippets


## working with time

https://pkg.go.dev/time#Time
https://gobyexample.com/time


``` go
package main

import (
    "fmt"
    "time"
)

func main() {
    p := fmt.Println

    now := time.Now()
    p(now)

    then := time.Date(
        2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
    p(then)

    p(then.Year())
    p(then.Month())
    p(then.Day())
    p(then.Hour())
    p(then.Minute())
    p(then.Second())
    p(then.Nanosecond())
    p(then.Location())

    p(then.Weekday())

    p(then.Before(now))
    p(then.After(now))
    p(then.Equal(now))

    diff := now.Sub(then)
    p(diff)

    p(diff.Hours())
    p(diff.Minutes())
    p(diff.Seconds())
    p(diff.Nanoseconds())

    p(then.Add(diff))
    p(then.Add(-diff))
}

```

``` go

package main

import (
    "fmt"
    "time"
)

func main() {
    currentTime := time.Now()
    fmt.Println("Current Time in String: ", currentTime.String())
    fmt.Println("MM-DD-YYYY : ", currentTime.Format("01-02-2006"))
    fmt.Println("YYYY-MM-DD : ", currentTime.Format("2006-01-02"))
    fmt.Println("YYYY.MM.DD : ", currentTime.Format("2006.01.02 15:04:05"))
    fmt.Println("YYYY#MM#DD {Special Character} : ", currentTime.Format("2006#01#02"))
    fmt.Println("YYYY-MM-DD hh:mm:ss : ", currentTime.Format("2006-01-02 15:04:05"))
    fmt.Println("Time with MicroSeconds: ", currentTime.Format("2006-01-02 15:04:05.000000"))
    fmt.Println("Time with NanoSeconds: ", currentTime.Format("2006-01-02 15:04:05.000000000"))
    fmt.Println("ShortNum Month : ", currentTime.Format("2006-1-02"))
    fmt.Println("LongMonth : ", currentTime.Format("2006-January-02"))
    fmt.Println("ShortMonth : ", currentTime.Format("2006-Jan-02"))
    fmt.Println("ShortYear : ", currentTime.Format("06-Jan-02"))
    fmt.Println("LongWeekDay : ", currentTime.Format("2006-01-02 15:04:05 Monday"))
    fmt.Println("ShortWeek Day : ", currentTime.Format("2006-01-02 Mon"))
    fmt.Println("ShortDay : ", currentTime.Format("Mon 2006-01-2"))
    fmt.Println("Short Hour Minute Second: ", currentTime.Format("2006-01-02 3:4:5"))
    fmt.Println("Short Hour Minute Second: ", currentTime.Format("2006-01-02 3:4:5 PM"))
    fmt.Println("Short Hour Minute Second: ", currentTime.Format("2006-01-02 3:4:5 pm"))
}
```

``` sql
-- Source - https://stackoverflow.com/a/1262055
-- Posted by Patrick, modified by community. See post 'Timeline' for change history
-- Retrieved 2026-02-13, License - CC BY-SA 2.5

--Create a table having a CURRENT_TIMESTAMP:
CREATE TABLE FOOBAR (
    RECORD_NO INTEGER NOT NULL,
    TO_STORE INTEGER,
    UPC CHAR(30),
    QTY DECIMAL(15,4),
    EID CHAR(16),
    RECORD_TIME NOT NULL DEFAULT CURRENT_TIMESTAMP)

--Create before update and after insert triggers:
CREATE TRIGGER UPDATE_FOOBAR BEFORE UPDATE ON FOOBAR
    BEGIN
       UPDATE FOOBAR SET record_time = datetime('now', 'localtime')
       WHERE rowid = new.rowid;
    END

CREATE TRIGGER INSERT_FOOBAR AFTER INSERT ON FOOBAR
    BEGIN
       UPDATE FOOBAR SET record_time = datetime('now', 'localtime')
       WHERE rowid = new.rowid;
    END

```

## TLS HTTPS certificate snippet

Based on the gowork/src/goecho-setup example

``` go
// Configuration variables for the server
const (
	HTTPPort  = ":80"
	HTTPSPort = ":443"
	CERT      = "localhost.crt"
	KEY       = "localhost.key"
)

// RedirectHTTP starts a parallel server that monitors port 80 and redirects TLS
func redirectHTTP() {
	e := echo.New()
	e.Use(middleware.HTTPSRedirect())
	go func() { e.Logger.Fatal(e.Start(HTTPPort)) }()
}

// Start will start the server
func Start(e *echo.Echo, TLS bool) {
	if TLS {
		redirectHTTP()
		e.Logger.Fatal(e.StartTLS(HTTPSPort, CERT, KEY))
	} else {
		e.Logger.Fatal(e.Start(HTTPPort))
	}
}
```
## Non CDN sources

<!-- jQuery -->
        <script src="/static/common/jquery.min.js"></script>

        <!-- Bootstrap CSS -->
        <link
            href="/static/common/bootstrap.min.css"
            rel="stylesheet"
            integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH"
            crossorigin="anonymous"
        />
        <!-- Bootstrap JS -->
        <script
            src="/static/common/bootstrap.bundle.min.js"
            integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz"
            crossorigin="anonymous"
        ></script>
        <!-- Toastify JS -->
        <script
            src="/static/common/toastify.js"
            integrity="sha512-MnKz2SbnWiXJ/e0lSfSzjaz9JjJXQNb2iykcZkEY2WOzgJIWVqJBFIIPidlCjak0iTH2bt2u1fHQ4pvKvBYy6Q=="
            crossorigin="anonymous"
            referrerpolicy="no-referrer"
        ></script>
        <!-- Toastify CSS-->
        <link
            rel="stylesheet"
            href="/static/common/toastify.css"
            integrity="sha512-VSD3lcSci0foeRFRHWdYX4FaLvec89irh5+QAGc00j5AOdow2r5MFPhoPEYBUQdyarXwbzyJEO7Iko7+PnPuBw=="
            crossorigin="anonymous"
            referrerpolicy="no-referrer"
        />
		<!-- Leaflet CSS -->
        <link
            rel="stylesheet"
            href="/static/common/leaflet.css"
            integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
            crossorigin=""
        />		
        
        <!-- Leaflet js-->
        <script
            src="/static/common/leaflet.js"
            integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo="
            crossorigin=""
        ></script>

## CDN sourcess

<!-- Toastify JS -->
        <script
            src="https://cdnjs.cloudflare.com/ajax/libs/toastify-js/1.6.1/toastify.js"
            integrity="sha512-MnKz2SbnWiXJ/e0lSfSzjaz9JjJXQNb2iykcZkEY2WOzgJIWVqJBFIIPidlCjak0iTH2bt2u1fHQ4pvKvBYy6Q=="
            crossorigin="anonymous"
            referrerpolicy="no-referrer"
        ></script>
        <!-- Toastify CSS-->
        <link
            rel="stylesheet"
            href="https://cdnjs.cloudflare.com/ajax/libs/toastify-js/1.6.1/toastify.css"
            integrity="sha512-VSD3lcSci0foeRFRHWdYX4FaLvec89irh5+QAGc00j5AOdow2r5MFPhoPEYBUQdyarXwbzyJEO7Iko7+PnPuBw=="
            crossorigin="anonymous"
            referrerpolicy="no-referrer"
        />
        <!-- Leaflet CSS -->
        <link
            rel="stylesheet"
            href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"
            integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
            crossorigin=""
        />
        <!-- jQuery -->
        <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.7.1/jquery.min.js"></script>

        <!-- Bootstrap CSS -->
        <link
            href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css"
            rel="stylesheet"
            integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH"
            crossorigin="anonymous"
        />
        <!-- Bootstrap JS -->
        <script
            src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"
            integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz"
            crossorigin="anonymous"
        ></script>
        <!-- Font Awesome -->
        <link
            href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.6.0/css/all.min.css"
            rel="stylesheet"
        />
		<!-- Leaflet CSS -->
        <link
            rel="stylesheet"
            href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"
            integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
            crossorigin=""
        />
        <!-- Leaflet js-->
        <script
            src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"
            integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo="
            crossorigin=""
        ></script>		