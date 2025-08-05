# 🌐 Web Page Analyzer

The **Web Page Analyzer** is a tool that analyzes a given web page and provides detailed technical insights. It helps developers, testers, and SEO specialists quickly evaluate the structure and content of a web page.

---

## Features

- **Page Title** – Displays the title of the web page.
- **HTML Version** – Detects the version of HTML used.
- **Login Form Detection** – Identifies if a login form is present on the page.
- **Headings Overview** – Shows the count and content of all H1–H6 tags.
- **Links Analysis** – Lists internal, external, and broken links with their HTTP status and response latency.

---

## Technologies Used

- **Backend:** Go (Golang) `v1.24.5`
- **Frontend:** Angular `v15`

---

## Design & Architecture

This tool is designed for simplicity and performance:

- The user provides a URL via a form on the web UI.
- Clicking **Run Analyzer** initiates an analysis of the target web page.
- The analysis runs in a **multi-threaded** manner for fast response times and scalable performance.

![Web Page Analyzer Architecture]([![go-analyzer.png](https://i.postimg.cc/jSW9qhhY/go-analyzer.png)](https://postimg.cc/zyrp0KMx))

---

## Project Structure (iidiomatic go)
- `analyzers`  - This holds the all types of analyzers.  
- `config`    -  This holds the application configuration.
- `engine`    -  This holds the router for handle request and response.
- `fetcher`    -  This acts like the bridge beetween alayzer and outer - (Validates the URL, fetch web pages from internet...).
- `handler`    -  This holds the http, model and middleware. it helps to run the analyzer after from route request.
- `pool`    -  This help to managing and running analyzers concurrently. in here WaitGroup is used.

##  Run & Deployment
You can run the project in two different ways:

- Run from Source
    `go run main.go` -> Requires Go installed on your PC. inside the root run this command..

- Using Docker - please refer the Dockerfile here for more details.
    `docker build -t webpage-analyzer .` - for build the docker image.
    `docker run -p 8080:8080 webpage-analyzer .` - for run the containers.

## Sample UI/UX 
- If you ran the project successfully

![Initial view][![web1.png](https://i.postimg.cc/ZYXJhj3p/web1.png)](https://postimg.cc/QVg2QQfd)

[![example.com](https://i.postimg.cc/3JpYrpC6/web2.png)](https://postimg.cc/pp2gqmSY)

[![srilankacricket.lk](https://i.postimg.cc/8zSD75ZC/web3.png)](https://postimg.cc/4mBjMJkC)

## API Endpoints
Method	Endpoint	Description
GET	/	Serves static frontend web content (Angular application)
GET	/v1/analyze?url=<URL>	Returns analysis report for the given URL

## Future Enhancements
1. Docker Compose support for frontend/backend.
2. More detailed link validation and reports.
3. Accessibility checks - Some URLs not allowd to analyze - forbidden.
4. Exportable analysis reports (PDF/CSV).
5. Database integration and time analyzis (Future purposes).


