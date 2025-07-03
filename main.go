// my-go-app/main.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template" // For HTML templating (if you choose to use it)
	"log"
	"net/http"
	"my-go-app/model" // Adjust module path if different
	"my-go-app/repository" // Adjust module path if different

	"golang.org/x/crypto/bcrypt"
)

// In-memory session store (for simplicity, not for production!)
var sessions = map[string]string{} // token -> username

func main() {
	port := "8000"

	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close() // Ensure database connection is closed when main exits

	// Serve static HTML files from the "static" directory
	// http.Handle("/", http.FileServer(http.Dir("./static"))) // Serves entire static dir
	// We'll serve index.html directly for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// User registration endpoint
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		user := &model.User{
			Username: req.Username,
			Password: string(hashedPassword),
		}

		err = repo.CreateUser(user)
		if err != nil {
			if err.Error() == "failed to create user: ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)" {
				http.Error(w, "Username already exists", http.StatusConflict)
			} else {
				log.Printf("Error creating user: %v", err)
				http.Error(w, "Error creating user", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User registered successfully!")
	})

	// User sign-in endpoint
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := repo.GetUserByUsername(req.Username)
		if err != nil {
			log.Printf("Error getting user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Compare the provided password with the hashed password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// For simplicity: generate a simple session token (NOT SECURE FOR PRODUCTION)
		sessionToken := "some_random_token_for_" + user.Username // In real app, use UUIDs and proper session management
		sessions[sessionToken] = user.Username

		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: sessionToken,
			Path:  "/",
			HttpOnly: true, // Prevent client-side JS access
			Secure:    false, // Set to true for HTTPS in production
			SameSite:  http.SameSiteLaxMode,
		})

		fmt.Fprintf(w, "Signed in successfully!")
	})

	// A protected endpoint (example)
	// A protected endpoint (example)
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
    		cookie, err := r.Cookie("session_token")
    		if err != nil {
        http.Redirect(w, r, "/", http.StatusFound) // Redirect to home/signin if no session cookie
        return
    	}

    // We no longer directly use 'username' or 'ok' in the Go code after switching to static file serving.
    // However, we still need to check if the session is valid.
   	 _, ok := sessions[cookie.Value] // Use '_' to discard the 'username' value as it's not used
    	if !ok {
        	http.Redirect(w, r, "/", http.StatusFound) // Redirect if session invalid
        	return
   		}

    		// Serve the dashboard HTML content from a static file
    		http.ServeFile(w, r, "static/dashboard.html")
    		// The username for display is now passed via the redirect URL to the frontend JavaScript.
	})


	log.Printf("Server starting on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Simple helper to load HTML templates (optional, can be expanded for more complex UIs)
func renderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles(tmpl)
	parsedTemplate.Execute(w, nil)
}
