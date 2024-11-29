package models

import "net/http"

// SessionKeys holds the access and secret key
type SessionKeys struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// Session contains authentication token
type Session struct {
	Md5sum_wizard_templates string `json:"md5sum_wizard_templates"`
	Md5sum_tenable_links    string `json:"md5sum_tenable_links"`
	Token                   string `json:"token"`
}

// Folder contains folder details
type Folder struct {
	UnreadCount *int   `json:"unread_count"`
	Custom      int    `json:"custom"`
	DefaultTag  int    `json:"default_tag"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
}

// Nessus contains nessus authentication and client details
type Nessus struct {
	Username   string
	Password   string
	Url        string
	Token      string
	AccessKey  string
	SecretKey  string
	ScanName   string
	HttpClient *http.Client
}

// Scan holds scan details
type Scan struct {
	ScanType             string  `json:"scan_type"`
	IsWas                *bool   `json:"isWas"`
	FolderID             int     `json:"folder_id"`
	Type                 *string `json:"type"`
	Read                 bool    `json:"read"`
	LastModificationDate int64   `json:"last_modification_date"`
	CreationDate         int64   `json:"creation_date"`
	Status               string  `json:"status"`
	UUID                 *string `json:"uuid"`
	Shared               bool    `json:"shared"`
	UserPermissions      int     `json:"user_permissions"`
	Owner                string  `json:"owner"`
	Timezone             *string `json:"timezone"`
	RRules               *string `json:"rrules"`
	Starttime            *string `json:"starttime"`
	Enabled              bool    `json:"enabled"`
	Control              bool    `json:"control"`
	LiveResults          int     `json:"live_results"`
	Name                 string  `json:"name"`
	ID                   int     `json:"id"`
}

// Vulnerability contains vulnerability information
type Vulnerability struct {
	PluginID                         string
	AssetName                        string
	PluginName                       string
	Severity                         string
	CVSS3BaseScore                   string
	VulnerabilityPriorityRating      string
	PluginOutput                     string
	Synopsis                         string
	Description                      string
	Solution                         string
	VulnerabilityFirstDiscoveredDate string
	VulnerabilityLastObservedDate    string
	PluginPublicationDate            string
	PatchPublicationDate             string
	Exploitability                   string
}

type ScanResponse struct {
	Folders   []Folder `json:"folders"`
	Scans     []Scan   `json:"scans"`
	Timestamp int64    `json:"timestamp"`
}

// ReportData holds report data for HTML generation
type ReportData struct {
	DateGenerated   string
	Vulnerabilities []Vulnerability
}
