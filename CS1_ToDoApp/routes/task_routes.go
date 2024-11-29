package routes

import (
	"CS1_ToDoApp/controllers" // Import the application logic (controllers)
	"bytes"
	"encoding/json"                   // Handle JSON data
	"github.com/natefinch/lumberjack" // Manage log files
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io" // Handle input/output operations
	"net/http"
	_ "net/http"
	"os" // Interact with the file system
	"regexp"
	"runtime/debug"
	"time" // Manage time and dates
)

var logger *zap.Logger

func initLogger() {
	// Get the current date for naming the log file
	currentTime := time.Now()
	formattedDate := currentTime.Format("02-01-2006") // Format: day-month-year

	// Configure the lumberjack logger for file rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename: "logs/" + formattedDate + ".log", // Log file named by date
		MaxSize:  10,                               // Maximum file size (MB)
		MaxAge:   7,                                // Retain logs for 7 days
		Compress: true,                             // Compress older logs
	}

	// Create a WriteSyncer for both console and file
	writerSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))

	// Set the log level and format for the logger
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 for timestamps

	// Create a zapcore for logging
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON format for log entries
		writerSyncer,
		zapcore.InfoLevel, // Set the default log level
	)

	// Create the logger
	logger = zap.New(core)

	// Set the global logger
	zap.ReplaceGlobals(logger)
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic occurred", zap.Any("error", r), zap.ByteString("stack_trace", debug.Stack()))
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

		// Log the request using zap
		logger.With(zap.String("method", r.Method+"-request"), zap.String("path", r.URL.Path)).
			Info("Incoming request", zap.Any("body", bodyJSON))

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

		// Log the response using zap
		if rw.statusCode < 400 {
			logger.With(zap.String("method", r.Method+"-response"), zap.String("path", r.URL.Path)).
				Info("Successful response", zap.String("body", rw.body.String()))
		} else {
			logger.With(zap.String("method", r.Method+"-response"), zap.String("path", r.URL.Path)).
				Error("Error response", zap.String("body", rw.body.String()))
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
	// Initialize the logger
	initLogger()

	// Create a new ServeMux instance to avoid global route conflicts
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.GetAllTasks(w, r)
		} else if r.Method == http.MethodPost {
			controllers.CreateNewTask(w, r)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
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

	handler := LoggingRequestMiddleware(LoggingResponseMiddleware(RecoveryMiddleware(mux)))
	// Return the configured router
	return handler
}
