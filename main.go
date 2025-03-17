package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SpotifyAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SearchResult struct {
	Tracks struct {
		Items []struct {
			Name     string `json:"name"`
			ID       string `json:"id"`
			Duration int    `json:"duration_ms"`
			Album    struct {
				Name   string `json:"name"`
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"album"`
			Artists []struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"artists"`
			PreviewURL string `json:"preview_url"`
		} `json:"items"`
	} `json:"tracks"`
}

var spotifyAuth SpotifyAuth
var authExpiry time.Time

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file, using environment variables")
	}

	r := gin.Default()

	r.GET("/api/search", searchTracks)
	r.GET("/api/track/:id", getTrackDetails)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func getSpotifyToken() error {
	if !authExpiry.IsZero() && time.Now().Before(authExpiry) {
		return nil
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("missing Spotify credentials")
	}

	authURL := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get token, status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&spotifyAuth)
	if err != nil {
		return err
	}

	authExpiry = time.Now().Add(time.Duration(spotifyAuth.ExpiresIn-60) * time.Second)
	return nil
}

func searchTracks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	err := getSpotifyToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate with Spotify"})
		return
	}

	searchURL := "https://api.spotify.com/v1/search"
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("type", "track")
	q.Add("limit", "10")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", spotifyAuth.TokenType+" "+spotifyAuth.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tracks"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": "Spotify API error"})
		return
	}

	var result SearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getTrackDetails(c *gin.Context) {
	trackID := c.Param("id")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "track ID is required"})
		return
	}

	err := getSpotifyToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate with Spotify"})
		return
	}

	trackURL := "https://api.spotify.com/v1/tracks/" + trackID
	req, err := http.NewRequest("GET", trackURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Authorization", spotifyAuth.TokenType+" "+spotifyAuth.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get track details"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": "Spotify API error"})
		return
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, result)
}