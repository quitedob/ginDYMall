package upload // Changed package name to 'upload' to match directory

import (
	"fmt"
	"mime/multipart"
	"path/filepath" // To help with getting extension
	"strings"       // For SanitizeFilename
	"time"          // For generating unique names

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials" // For explicit credential configuration if needed
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	// "douyin/config" // Assuming config.Cfg.OSS might exist for credentials
)

// Client wraps S3 operations.
type Client struct {
	svc    *s3.S3
	bucket string
	region string // Store region for constructing URL if needed, or use a config base URL
}

// Config holds configuration for the OSS client.
// Expect these to be populated from your project's config (e.g., Viper)
type Config struct {
	Type            string // e.g., "s3"
	Endpoint        string // For S3 compatible storage, otherwise AWS default
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	// BaseURL      string // Optional: if you want to force a base URL for links
}

// NewClient creates a new S3 client.
func NewClient(cfg Config) (*Client, error) {
	awsCfg := aws.NewConfig().WithRegion(cfg.Region)

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		awsCfg = awsCfg.WithCredentials(credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""))
	}
	if cfg.Endpoint != "" { // For S3-compatible storage like MinIO
		awsCfg = awsCfg.WithEndpoint(cfg.Endpoint).WithS3ForcePathStyle(true)
	}

	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := s3.New(sess)
	return &Client{svc: svc, bucket: cfg.Bucket, region: cfg.Region}, nil
}

// Upload uploads a file to S3 and returns its public URL.
// The 'uploadPath' is the path within the bucket (e.g., "avatars/", "products/").
func (c *Client) Upload(uploadPath string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Generate a unique filename to prevent overwrites
	ext := filepath.Ext(fileHeader.Filename)
	// Sanitize and truncate the base filename, ensuring total length is reasonable. Max 50 for base.
	sanitizedBase := SanitizeFilename(strings.TrimSuffix(fileHeader.Filename, ext), 50)
	uniqueFileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), sanitizedBase, ext)
	key := filepath.Join(uploadPath, uniqueFileName) // Example: "avatars/1678886400000-my-image.jpg"

	_, err = c.svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        file,
		ACL:         aws.String("public-read"), // As per issue example
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")), // Set content type
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct URL. This might vary based on S3 setup (path-style vs virtual-hosted, custom domain)
	// Default virtual-hosted style:
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.bucket, c.region, key)
	// If using custom endpoint (like MinIO) or need path style, URL construction will change.
	// Example for endpoint-based URL:
	if c.svc.Endpoint != "" && !strings.Contains(c.svc.Endpoint, "amazonaws.com") { // Basic check if it's not AWS default
		if strings.HasSuffix(c.svc.Endpoint, "/") {
			url = fmt.Sprintf("%s%s/%s", c.svc.Endpoint, c.bucket, key)
		} else {
			url = fmt.Sprintf("%s/%s/%s", c.svc.Endpoint, c.bucket, key)
		}
		// Check if S3ForcePathStyle is true, usually it is for MinIO
        if aws.BoolValue(c.svc.Config.S3ForcePathStyle) {
             url = fmt.Sprintf("%s/%s/%s", c.svc.Endpoint, c.bucket, key)
        } else {
            // Virtual hosted style for custom endpoint might be different, e.g. https://bucket.endpoint/key
            // This part might need adjustment based on specific provider's URL format for virtual hosted style.
            // For now, sticking to path style if endpoint is custom.
             url = fmt.Sprintf("https://%s.%s/%s", c.bucket, strings.Replace(c.svc.Endpoint,"https://","",1), key)
             // A more robust way for virtual hosted with custom endpoint:
             // Parse c.svc.Endpoint to get scheme and host, then construct: scheme://bucket.host/key
        }
	}


	return url, nil
}

// SanitizeFilename truncates and sanitizes the filename.
// This is a basic sanitizer.
func SanitizeFilename(filename string, maxLength int) string {
	sanitized := ""
	for _, r := range filename {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			if r == '.' && strings.Contains(sanitized, ".") { // Allow only one dot for extension
				continue
			}
			sanitized += string(r)
		}
	}
	// Remove leading/trailing dots or hyphens that might have been formed
    sanitized = strings.Trim(sanitized, ".-_")

	base := sanitized
	ext := ""
	if dotIndex := strings.LastIndex(sanitized, "."); dotIndex != -1 {
		base = sanitized[:dotIndex]
		ext = sanitized[dotIndex:]
	}

	if len(base) > maxLength {
		base = base[:maxLength]
	}
	return base + ext
}
