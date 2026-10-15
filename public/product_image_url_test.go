package public

import (
	"context"
	"testing"
)

func TestResolveProductImageURLUsesS3PublicURL(t *testing.T) {
	t.Setenv("S3_BUCKET_NAME", "demo-bucket")
	t.Setenv("AWS_REGION", "ap-south-1")

	got := resolveProductImageURL(context.Background(), "products/test.jpg")
	want := "https://demo-bucket.s3.ap-south-1.amazonaws.com/products/test.jpg"

	if got != want {
		t.Fatalf("resolveProductImageURL() = %q, want %q", got, want)
	}
}
