# StatusQuo

Very simple status monitoring for websites

A simple but functional, responsive web application that monitors registered URLs (by pinging them on a regular, defined frequency).  Each website is displayed with it's current status (up or down) along with basic metrics: 

* percentage uptime
* cumulative moving average (mean) response time
* [Apdex score](https://en.wikipedia.org/wiki/Apdex) (along with % satisfied, tolerating, frustrated)
* frequency of check

![Screenshot](https://github.com/james-bowman/statusquo/screenshot.png "StatusQuo Screenshot")

The application is deployed as a single, self contained, binary executable only requiring a separate .json configuration file to define the details of the websites/endpoints to monitor.

StatusQuo is developed using Go for the backend and HTML/JavaScript with Vue.js for the front end.

## Future features

* Reporting including plotting response times and uptime over a historical period.
* Download response times and uptime as CSV files
* Web hooks for alerting when sites go down.