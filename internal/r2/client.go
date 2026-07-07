package r2

import (
 "bytes"
 "context"
 "net"
 "net/http"
 "os"
 "time"

 "stream-sync/internal/config"

 "github.com/aws/aws-sdk-go-v2/aws"
 awsconfig "github.com/aws/aws-sdk-go-v2/config"
 "github.com/aws/aws-sdk-go-v2/credentials"
 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
 S3     *s3.Client
 Bucket string
}

func NewClient(cfg *config.Config) (*Client, error) {

 transport := &http.Transport{
  MaxIdleConns:        500,
  MaxIdleConnsPerHost: 500,
  MaxConnsPerHost:     500,
  IdleConnTimeout:     90 * time.Second,
  DisableCompression:  true,
  DialContext: (&net.Dialer{
   Timeout:   5 * time.Second,
   KeepAlive: 30 * time.Second,
  }).DialContext,
 }

 httpClient := &http.Client{
  Transport: transport,
  Timeout:   0,
 }

 awsCfg, err := awsconfig.LoadDefaultConfig(
  context.Background(),
  awsconfig.WithHTTPClient(httpClient),
  awsconfig.WithRegion("auto"),
  awsconfig.WithCredentialsProvider(
   credentials.NewStaticCredentialsProvider(
    cfg.R2.AccessKey,
    cfg.R2.SecretKey,
    "",
   ),
  ),
 )
 if err != nil {
  return nil, err
 }

 client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
  o.BaseEndpoint = aws.String(cfg.R2.Endpoint)
  o.UsePathStyle = true
 })

 return &Client{
  S3:     client,
  Bucket: cfg.R2.Bucket,
 }, nil
}

func (c *Client) Upload(localPath string, objectName string) error {

 file, err := os.Open(localPath)
 if err != nil {
  return err
 }
 defer file.Close()

 _, err = c.S3.PutObject(context.Background(), &s3.PutObjectInput{
  Bucket:       &c.Bucket,
  Key:          &objectName,
  Body:         file,
  ContentType:  aws.String("video/mp2t"),
  CacheControl: aws.String("public, max-age=3600"),
 })

 return err
}

func (c *Client) UploadBytes(data []byte, objectName string) error {

 _, err := c.S3.PutObject(context.Background(), &s3.PutObjectInput{
  Bucket:       &c.Bucket,
  Key:          &objectName,
  Body:         bytes.NewReader(data),
  ContentType:  aws.String("application/x-mpegURL"),
  CacheControl: aws.String("no-cache, no-store, must-revalidate"),
 })

 return err
}

func (c *Client) Exists(objectName string) bool {

 _, err := c.S3.HeadObject(context.Background(), &s3.HeadObjectInput{
  Bucket: &c.Bucket,
  Key:    &objectName,
 })

 return err == nil
}

func (c *Client) Delete(objectName string) error {

 _, err := c.S3.DeleteObject(context.Background(), &s3.DeleteObjectInput{
  Bucket: &c.Bucket,
  Key:    &objectName,
 })

 return err
}

func (c *Client) List(prefix string) ([]string, error) {

 var keys []string

 paginator := s3.NewListObjectsV2Paginator(c.S3, &s3.ListObjectsV2Input{
  Bucket: &c.Bucket,
  Prefix: &prefix,
 })

 for paginator.HasMorePages() {

  page, err := paginator.NextPage(context.Background())
  if err != nil {
   return nil, err
  }

  for _, obj := range page.Contents {
   if obj.Key != nil {
    keys = append(keys, *obj.Key)
   }
  }
 }

 return keys, nil
}