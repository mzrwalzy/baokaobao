package oss

import (
	"baokaobao/internal/config"
	"fmt"
	"io"
	"strings"

	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSClient wraps the Aliyun OSS SDK client and bucket.
type OSSClient struct {
	client   *alioss.Client
	bucket   *alioss.Bucket
	endpoint string
	bucketName string
}

// NewOSSClient creates an OSSClient from the given config.
// Returns nil if the config is not fully set (credentials missing).
func NewOSSClient(cfg config.OSSConfig) *OSSClient {
	if !isConfigured(cfg) {
		return nil
	}

	client, err := alioss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil
	}

	return &OSSClient{
		client:     client,
		bucket:     bucket,
		endpoint:   cfg.Endpoint,
		bucketName: cfg.Bucket,
	}
}

// IsConfigured returns true if the OSS client has been initialized.
func (c *OSSClient) IsConfigured() bool {
	return c != nil && c.bucket != nil
}

// Upload uploads data from the reader to OSS at the given objectKey and returns the public URL.
func (c *OSSClient) Upload(file io.Reader, objectKey string) (string, error) {
	if c == nil || c.bucket == nil {
		return "", fmt.Errorf("oss client not configured")
	}

	if err := c.bucket.PutObject(objectKey, file); err != nil {
		return "", fmt.Errorf("oss upload failed: %w", err)
	}

	// Build public URL: https://{bucket}.{endpoint}/{objectKey}
	url := fmt.Sprintf("https://%s.%s/%s", c.bucketName, c.endpoint, objectKey)
	return url, nil
}

// isConfigured checks whether all required OSS fields are present.
func isConfigured(cfg config.OSSConfig) bool {
	return strings.TrimSpace(cfg.Endpoint) != "" &&
		strings.TrimSpace(cfg.AccessKeyID) != "" &&
		strings.TrimSpace(cfg.AccessKeySecret) != "" &&
		strings.TrimSpace(cfg.Bucket) != ""
}
