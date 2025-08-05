# 🌐 Web Page Analyzer

The **Web Page Analyzer** is a tool that analyzes a given web page and provides detailed technical insights. It helps developers, testers, and SEO specialists quickly evaluate the structure and content of a web page.

---

## 📊 Features

- **Page Title** – Displays the title of the web page.
- **HTML Version** – Detects the version of HTML used.
- **Login Form Detection** – Identifies if a login form is present on the page.
- **Headings Overview** – Shows the count and content of all H1–H6 tags.
- **Links Analysis** – Lists internal, external, and broken links with their HTTP status and response latency.

---

## 🛠️ Technologies Used

- **Backend:** Go (Golang) `v1.24.5`
- **Frontend:** Angular `v15`

---

## 🧩 Design & Architecture

This tool is designed for simplicity and performance:

- The user provides a URL via a form on the web UI.
- Clicking **Run Analyzer** initiates an analysis of the target web page.
- The analysis runs in a **multi-threaded** manner for fast response times and scalable performance.

> 📌 *This section can be enhanced with an architecture diagram.*

---

## 📁 Project Structure (iidiomatic go)
`analyzers`  -> This holds the all types of analyzers.
`config`    -> This holds the application configuration.
`engine`    -> This holds the router for handle request and response.
`fetcher`    -> This acts like the bridge beetween alayzer and outer - (Validates the URL, fetch web pages from internet...).
`handler`    -> This holds the http, model and middleware. it helps to run the analyzer after from route request.
`pool`    -> This help to managing and running analyzers concurrently. in here WaitGroup is used.

## 📊 Sample Architecture (Image Placeholder)
![Web Page Analyzer Architecture](https://www.canva.com/design/DAGvMgPw7eQ/so7hyjJ7WTugcT1LPKhUzA/view?utm_content=DAGvMgPw7eQ&utm_campaign=designshare&utm_medium=link2&utm_source=uniquelinks&utlId=hcc75bbada1)

## 🚀 Run & Deployment
You can run the project in two different ways:

- Run from Source
    `go run main.go` -> Requires Go installed on your PC. inside the root run this command..


## 🔌 API Endpoints
Method	Endpoint	Description
GET	/	Serves static frontend web content (Angular application)
GET	/v1/analyze?url=<URL>	Returns analysis report for the given URL


