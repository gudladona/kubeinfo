package models

// Response model defines the JSON response to be returned to the client
type Response struct {
	PodCount       int    `json:"pod_count,omitempty"`
	Message string `json:"message"`
}
