package routes

import (
	"CS1_ToDoApp/controllers" // Import the application logic (controllers)
	"bytes"
	"encoding/json"                   // Handle JSON data
	"github.com/natefinch/lumberjack" // Manage log files
	"github.com/sirupsen/logrus"      // Powerful logging library
	"io"                              // Handle input/output operations
	"log"
	"net/http"
	_ "net/http"
	"os" // Interact with the file system
	"regexp"
	"runtime/debug"
	"time" // Manage time and dates
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic occurred: %v\nStack Trace: %s", r, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Middleware to log incoming requests
func LoggingRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// Recover from panic if any
			}
		}()

		// Read and log the request body
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var bodyJSON map[string]interface{}
		err := json.Unmarshal(bodyBytes, &bodyJSON)
		if err != nil {
			bodyJSON = nil
		}

		logrus.WithFields(logrus.Fields{
			"method": r.Method + "-request",
			"path":   r.URL.Path,
			"body":   bodyJSON,
		}).Info()

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Wrapper for ResponseWriter to capture response body
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Middleware to log outgoing responses
func LoggingResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// Recover from panic if any
			}
		}()

		// Wrap the ResponseWriter
		rw := &responseWriterWrapper{
			ResponseWriter: w,
			body:           bytes.NewBufferString(""),
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log the response
		if rw.statusCode < 400 {
			logrus.WithFields(logrus.Fields{
				"method": r.Method + "-response",
				"path":   r.URL.Path,
				"body":   rw.body.String(),
			}).Info()
		} else {
			logrus.WithFields(logrus.Fields{
				"method": r.Method + "-response",
				"path":   r.URL.Path,
				"body":   rw.body.String(),
			}).Error()
		}
	})
}

func extractIDFromURL(path string) string {
	// Simple regex to match /tasks/<id>
	re := regexp.MustCompile(`/tasks/(\d+)`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// Function to set up the router and define endpoints
func SetupRouter() http.Handler {
	// Get the current date for naming the log file
	currentTime := time.Now()
	formattedDate := currentTime.Format("02-01-2006") // Format: day-month-year

	// Configure the logger to write to both console and file
	multiWriter := io.MultiWriter(
		os.Stdout, // Log to the terminal
		&lumberjack.Logger{
			Filename: "logs/" + formattedDate + ".log", // Log file named by date
			MaxSize:  10,                               // Maximum file size (MB)
			MaxAge:   7,                                // Retain logs for 7 days
		},
	)

	// Set up Logrus with the configured logger
	logrus.SetOutput(multiWriter)                // Log to both console and file
	logrus.SetFormatter(&logrus.JSONFormatter{}) // Use JSON format for structured logging

	// Define the API routes (endpoints)
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.GetAllTasks(w, r)
		} else if r.Method == http.MethodPost {
			controllers.CreateNewTask(w, r)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id := extractIDFromURL(r.URL.Path)

		if id == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			controllers.GetTaskByID(w, r)
		} else if r.Method == http.MethodPut {
			controllers.UpdateTask(w, r)
		} else if r.Method == http.MethodDelete {
			controllers.DeleteTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	handler := LoggingRequestMiddleware(LoggingResponseMiddleware(RecoveryMiddleware(http.DefaultServeMux)))
	// Return the configured router
	return handler
}
