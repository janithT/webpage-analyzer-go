🌐 Web Page Analyzer
The Web Page Analyzer is a tool that inspects any given web page and provides detailed insights, including:

Page Title – Displays the web page’s title.

HTML Version – Detects the version of HTML used.

Login Form Detection – Identifies if a login form is present.

Headings Overview – Shows the count and content of <h1> to <h6> tags.

Links Analysis – Lists internal, external, and broken links along with their HTTP status and latency.

🛠️ Technologies Used
Backend: Go (Golang) v1.24.5

Frontend: Angular v15

🧩 Design & Architecture
The tool is designed for simplicity and performance:

Enter a URL via the input field.

Click Run Analyzer.

The system will analyze the given URL and display results for:

Page title

HTML version

Login form presence

Headings (H1–H6)

Link categorization and latency

🧵 Multi-threaded Execution
The analyzer runs in a multi-threaded manner to ensure optimal response time and scalability.

This allows simultaneous execution of different analysis components like HTML parsing, link checking, etc.

📊 Sample Architecture (Image Placeholder)
Insert architecture diagram here

📁 Development Structure
bash
Copy
Edit
project-root/
├── web/                # Angular frontend application
├── main.go             # Go entry point for backend
├── app.yaml            # Configuration file (used in binary execution)
├── ...
🚀 Run & Deployment
You can run the project in two different ways:

1. Run from Source
bash
Copy
Edit
go run main.go
Requires Go installed. Useful for development and testing.

2. Run from Compiled Binary
bash
Copy
Edit
./webpage-analyzer
Ensure that app.yaml is present in the same directory as the binary or in the correct relative path.

🔌 API Endpoints
Method	Endpoint	Description
GET	/	Serves static frontend web content (Angular application)
GET	/v1/analyze?url=<URL>	Returns analysis report for the given URL


