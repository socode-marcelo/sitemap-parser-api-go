package main

import (
	"encoding/xml"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Sitemap represents a sitemap.
type Sitemap struct {
	URLs    []SitemapURL    `xml:"url"`
	Sitemaps []SitemapSitemap `xml:"sitemap"`
}

// SitemapURL represents a URL in a sitemap.
type SitemapURL struct {
	Loc string `xml:"loc"`
}

// SitemapSitemap represents a sitemap in a sitemap index.
type SitemapSitemap struct {
	Loc string `xml:"loc"`
}

// parseSitemapFromRobotsTxt function is used to parse the sitemap from robots.txt.
// Input: robotsTxt string
// Output: sitemap string
func parseSitemapFromRobotsTxt(robotsTxt string) string {

	// Split the robotsTxt by new line character to get the lines
	lines := strings.Split(robotsTxt, "\n")

	// Loop over each line in the lines
	for _, line := range lines {

		// Check if the line starts with "Sitemap:"
		if strings.HasPrefix(line, "Sitemap:") {

			// If it does, then trim the prefix "Sitemap: " from the line
			return strings.TrimPrefix(line, "Sitemap: ")
		}
	}

	// If no sitemap is found, return an empty string
	return ""
}

// Function to check if a given domain string is valid
// @param domain: string - The domain to be checked
// @return: bool - True if the domain is valid, False otherwise
func isValidDomain(domain string) bool {
	// Check if the input is a string
	if len(domain) == 0 {
		return false // Input string is empty
	}

	// Check if the domain starts with "http://" or "https://",
	// if not, prepend it to the domain for proper URL parsing
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}

	parsedURL, err := url.Parse(domain)
	if err != nil {
		return false // Error parsing URL
	}

	// Check if the parsed URL's Host field is not empty
	return parsedURL.Host != "" // URL's Host field is empty
}

// Function to extract a domain from a String
// @param domain: string - The domain to be extracted
// @return: string - The extracted domain
func extractDomain(domain string) string {
	// Check if the input is a string
	if len(domain) == 0 {
		return "" // Input string is empty
	}

	// Prepend "http://" if domain does not start with it
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "http://" + domain
	}

	parsedURL, err := url.Parse(domain)
	if err != nil {
		return "" // Error parsing URL
	}

	// Return the Host field of the parsed URL
	return parsedURL.Host
}

// getSitemapURLFromDomain retrieves the sitemap URL from the given domain.
//
// It takes a domain string as a parameter and returns a string and an error.
func getSitemapURLFromDomain(domain string) (string, error) {
	// Check if the domain is valid. If not, return an error.
	if !isValidDomain(domain) {
		return "", fmt.Errorf("Failed to validate %s", domain)
	}

	// Extract the domain from the input.
	domain = extractDomain(domain)

	// Create an HTTP client with a timeout of 3 seconds.
	client := http.Client{Timeout: time.Second * 3}

	// If no sitemap is found, fetch the robots.txt file.
	robotsURL := fmt.Sprintf("https://%s/robots.txt", domain)
	resp, err := client.Get(robotsURL)
	if err != nil {
		// If the request fails, return the error.
		return "", err
	}
	defer resp.Body.Close()

	// If the response status is OK, parse the sitemap URL from the robots.txt file.
	if resp.StatusCode == http.StatusOK {
		robotsTxt, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		sitemapLoc := parseSitemapFromRobotsTxt(string(robotsTxt))
		// Check if sitemapLoc could be extracted from robots.txt
		if sitemapLoc != "" {
			return sitemapLoc, nil
		}
	}

	// Define a list of possible sitemap locations.
	sitemapLocations := []string{
		"/test.xml",
		"/sitemap.xml",
		"/sitemap1.xml",
		"/sitemap.txt",
		"/sitemap_index.xml",
		"/sitemap/",
		"/sitemap",
		"/sitemap-index.xml",
		"/sitemaps/",
		"/sitemaps",
		"/site-map",
		"/sitemap-indexes/",
		"/post-sitemap.xml",
		"/page-sitemap.xml",
		"/category-sitemap.xml",
		"/tag-sitemap.xml",
		"/pages-sitemap.xml",
		"/blog-pages-sitemap.xml",
		"/member-profile-sitemap.xml",
		"/dynamic-pages-sitemap.xml",
		"/other-pages-sitemap.xml",
		"/sitemap.xml.gz",
		"/sitemapindex.xml",
		"/sitemap_index.xml.gz",
		"/sitemap/index.xml",
		"/sitemap.xml",
		"/sitemap_map.html",
		"/wp-sitemap.xml",
		"/other-pages-sitemap.xml",
		"/category-sitemap.xml",
		"/tag-sitemap.xml",
		"/author-sitemap.xml",
		"/post-sitemap",
		"/sitemaps-2-sitemap.xml",
		"/page-sitemap",
	}

	// Loop through each sitemap location.
	for _, location := range sitemapLocations {
		// Construct the URL.
		url := fmt.Sprintf("https://%s%s", domain, location)
		// Send a GET request to the URL.
		resp, err := client.Get(url)
		if err != nil {
			// If the request fails, return the error.
			return "", err
		}
		defer resp.Body.Close()

		// If the response status is OK, return the URL.
		if resp.StatusCode == http.StatusOK {
			return url, nil
		}
	}

	// If the URL cannot be retrieved, return an error.
	return "", fmt.Errorf("Couldn't find sitemap for %s", domain)
}

// parseSitemap parses a sitemap URL and returns a slice of URLs found in the sitemap.
//
// It takes a string parameter named 'url' which specifies the URL of the sitemap.
// The function returns a slice of strings ([]string) containing the URLs found in the sitemap,
// and an error if there was an error during the parsing process.
func parseSitemap(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sitemap Sitemap
	err = xml.Unmarshal(body, &sitemap)
	if err != nil {
		return nil, err
	}

	// If sitemap contains URLs, return them
	if len(sitemap.URLs) > 0 {
		urls := make([]string, len(sitemap.URLs))
		for i, u := range sitemap.URLs {
			urls[i] = u.Loc
		}
		return urls, nil
	}

	// If sitemap contains sitemaps, parse each of them
	var sitemapIndex Sitemap
	err = xml.Unmarshal(body, &sitemapIndex)
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(sitemapIndex.Sitemaps))
	for i, s := range sitemapIndex.Sitemaps {
		subUrls, err := parseSitemap(s.Loc)
		if err != nil {
			return nil, err
		}
		urls[i] = fmt.Sprintf("Sitemap index: %s", s.Loc)
		urls = append(urls, subUrls...)
	}

	return urls, nil
}

// handleRequest handles the HTTP request for both domain and sitemap endpoints.
//
// It expects a POST request and validates the method.
// It decodes the JSON payload and checks for errors.
// It retrieves the required field (domain or sitemap) from the payload and checks for its existence.
// It processes the request based on the specified requestType ('domain' or 'sitemap').
// It constructs a JSON response with the parsed URLs.
// It marshals the JSON response and checks for errors.
// It sets the Content-Type header to "application/json".
// It writes the JSON response to the HTTP response writer.
func handleRequest(w http.ResponseWriter, r *http.Request, requestType string) {

	// Check if the request method is POST
	if r.Method != http.MethodPost {
		// If it's not POST, return a method not allowed error
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON payload
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		// If the JSON payload is invalid, return a bad request error
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Get the value of the request type field
	fieldValue, exists := data[requestType]
	if !exists {
		// If the request type field is missing, return a bad request error
		http.Error(w, fmt.Sprintf("Missing '%s' field in JSON payload", requestType), http.StatusBadRequest)
		return
	}

	fmt.Println(requestType, fieldValue)

	// Declare the URLs slice and the parse error
	var urls []string
	var parseErr error

	// If the request type is "domain", get the sitemap URL from the domain
	if requestType == "domain" {
		sitemapURL, err := getSitemapURLFromDomain(fieldValue)
		if err != nil {
			// If an error occurs, return an internal server error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		urls, parseErr = parseSitemap(sitemapURL)
	} else if requestType == "sitemap" {
		// check if fieldValue is a valid URL
		_, err := url.ParseRequestURI(fieldValue)
		if err != nil {
		    // If fieldValue is not a valid URL, return a bad request error
		    http.Error(w, "Invalid URL", http.StatusBadRequest)
		    return
		}
		// If the request type is "sitemap", parse the sitemap
		urls, parseErr = parseSitemap(fieldValue)
	}

	// If an error occurs while parsing the sitemap, return an internal server error
	if parseErr != nil {
		http.Error(w, "Failed to parse sitemap", http.StatusInternalServerError)
		return
	}

	// Create the response
	response := map[string]interface{}{
		// "errors": []string{}, // TODO add errors
		"type":   requestType,
		"urls":   urls,
	}

	// Marshal the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// If an error occurs, return an internal server error
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response with a status code of OK
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

// handleDomain handles the HTTP request for the domain endpoint.
func handleDomainEndpoint(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, "domain")
}

// handleSitemap handles the HTTP request for the sitemap endpoint.
func handleSitemapEndpoint(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, "sitemap")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "To request URLs, POST the link to /sitemap as {\"sitemap\":\"https://stackovercode.com/sitemap.xml\"}")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Pong!")
}

func main() {
	http.HandleFunc("/sitemap", handleSitemapEndpoint)
	http.HandleFunc("/domain", handleDomainEndpoint)
	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/", handleRoot)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
