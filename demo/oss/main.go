package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	endpoint           = "https://oss-cn-beijing.aliyuncs.com"
	accessKeyIDEnv     = "OSS_ACCESS_KEY_ID"
	accessKeySecretEnv = "OSS_ACCESS_KEY_SECRET"
)

func listBuckets(client *oss.Client) {
	result, err := client.ListBuckets() // 默认情况下一次返回 100 条记录
	if err != nil {
		panic(err)
	}

	fmt.Println("Total Buckets:", len(result.Buckets))
	for _, bkt := range result.Buckets {
		// SetBucketACL
		// GetBucketACL
		acl, _ := client.GetBucketACL(bkt.Name)      // public-read, private
		lc, _ := client.GetBucketLifecycle(bkt.Name) // 目前线上 bucket 并没有做 life cycle 管理
		fmt.Println("Bucket Name:", bkt.Name, bkt.StorageClass, acl.ACL)
		fmt.Printf("Bucket Lifecycle:%v\n", lc.Rules)
	}

	// 完整列举: ListBuckets + marker
	// marker := ""
	// for {
	// 	lsRes, err := client.ListBuckets(oss.Marker(marker))
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 		os.Exit(-1)
	// 	}

	// 	// 默认情况下一次返回100条记录。
	// 	for _, bucket := range lsRes.Buckets {
	// 		fmt.Println("Bucket: ", bucket.Name)
	// 	}

	// 	if lsRes.IsTruncated {  // 返回结果集被截断了，说明还有数据，继续访问
	// 		marker = lsRes.NextMarker
	// 	} else {
	// 		break
	// 	}
	// }
}

func createBucket(client *oss.Client) {
	// oss.Region()
	if err := client.CreateBucket("zstax-yao-test"); err != nil {
		panic(err)
	}
}

func deleteBucket(client *oss.Client) {
	if err := client.DeleteBucket("zstax-yao-test"); err != nil {
		panic(err)
	}
}

func putObject(client *oss.Client) {
	bucket, err := client.Bucket("zstax-yao-test")
	if err != nil {
		panic(err)
	}

	if err = bucket.PutObject("data19.txt", strings.NewReader("test")); err != nil {
		panic(err)
	}

	// 私有 bucket 生成 sign url，通过该 url 在一定时间内可以访问 bucket 内的资源（存在过期时间）
	// https://github.com/aliyun/aliyun-oss-go-sdk/blob/master/sample/sign_url.go?spm=a2c4g.31952.0.0.30a57d22yFyp6H&file=sign_url.go
	fmt.Println(bucket.SignURL("data10.txt", oss.HTTPGet, 3600))
	// if err = bucket.PutObjectFromFile()
}

func getObject(client *oss.Client) {
	bucket, err := client.Bucket("zstax-yao-test")
	if err != nil {
		panic(err)
	}

	r, err := bucket.GetObject("data4.txt")
	if err != nil {
		panic(err)
	}
	defer r.Close() // close 操作

	// oss.Range()

	// etag := meta.Get(oss.HTTPHeaderEtag)
	// oss.IfMatch(etag)

	data, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	fmt.Println("GetObject:", string(data))
}

func putObjectFromFile(client *oss.Client) {
	bucket, err := client.Bucket("zstax-yao-test")
	if err != nil {
		panic(err)
	}

	if err = bucket.PutObjectFromFile("data3.txt", "./data.txt"); err != nil {
		panic(err)
	}
}

func getObjectFromFile(client *oss.Client) {
	bucket, err := client.Bucket("zstax-yao-test")
	if err != nil {
		panic(err)
	}

	if err = bucket.GetObjectToFile("data2.txt", "result.txt"); err != nil {
		panic(err)
	}
}

func listObjects(client *oss.Client) {
	bucket, err := client.Bucket("zstax-yao-test")
	if err != nil {
		panic(err)
	}

	result, err := bucket.ListObjects()
	if err != nil {
		panic(err)
	}

	for _, obj := range result.Objects {
		fmt.Println("Object:", obj.Key)
	}

	// marker := oss.Marker("")
	// for {
	// 	lor, err = bucket.ListObjects(oss.MaxKeys(3), marker)
	// 	if err != nil {
	// 		HandleError(err)
	// 	}
	// 	marker = oss.Marker(lor.NextMarker)
	// 	fmt.Println("my objects page :", getObjectsFormResponse(lor))
	// 	if !lor.IsTruncated {
	// 		break
	// 	}
	// }
}

func deleteObject(client *oss.Client) {
	bucket, err := client.Bucket("zstax-yao-test")
	if err != nil {
		panic(err)
	}

	// bucket.UploadFile()
	// oss.Checkpoint()

	if err = bucket.DeleteObject("data3.txt"); err != nil {
		panic(err)
	}
}

func appendObject(client *oss.Client) {
	bucket, err := client.Bucket("zstax-yao-test")
	if err != nil {
		panic(err)
	}

	// bucket.InitiateMultipartUpload()
	// bucket.UploadPart()
	// bucket.UploadFile()
	// bucket.UploadFile()

	var nextPos int64
	for i := 0; i < 10; i++ {
		nextPos, err = bucket.AppendObject("data4.txt", strings.NewReader(fmt.Sprintf("Test: %d\n", i)), nextPos)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Append Success!")

	// 通过 Object 元信息，获取下一次的写入位置
	// props, err := bucket.GetObjectDetailedMeta(objectKey)
	// nextPos, err = strconv.ParseInt(props.Get(oss.HTTPHeaderOssNextAppendPosition), 10, 64)
}

// func getBucketLifeCycleConf(client *oss.Client) {
// 	conf, err := client.GetBucketLifecycle()
// }

func main() {
	//fmt.Println(filepath.Ext("cms.oss.json"))
	fmt.Println(os.Getenv(accessKeyIDEnv))
	fmt.Println(os.Getenv(accessKeySecretEnv))
	client, err := oss.New(endpoint, os.Getenv(accessKeyIDEnv), os.Getenv(accessKeySecretEnv))
	if err != nil {
		panic(err)
	}

	putObject(client)

	// listBuckets(client)

	// listBuckets(client)
	// createBucket(client)
	// listBuckets(client)

	// getObject(client)
	// putObjectFromFile(client)
	// getObjectFromFile(client)
	// putObjectFromFile(client)
	// listObjects(client)
	// deleteObject(client)
	// listObjects(client)

	// oss.Routines(3)
	// oss.ObjectStorageClass()
	// oss.Checkpoint()
	// listObjects(client)
	// appendObject(client)
	// getObject(client)
	// listObjects(client)

	// oss.StorageArchive

	// bucket.RestoreObject() : 数据解冻操作
	//
	// client.PutBucketAccessMonitor()

	// oss.LifecycleConfiguration

	// client.SetBucketLogging()
	// SetBucketLogging 设置 bucket 访问日志
	// oss.LifecycleTransition.IsAccessTime：是否按照访问时间进行生命周期管理

	// oss.SplitFileByPartNum()
}

// 更多 demo 参考：https://github.com/aliyun/aliyun-oss-go-sdk/blob/master/sample/put_object.go

/*
Total Buckets: 15
Bucket Name: doupingtai
Bucket Name: live-reservation-delay-queue-log
Bucket Name: soft-taxbook
Bucket Name: taxbook-crawling
Bucket Name: taxbook-def
Bucket Name: taxbook-info
Bucket Name: taxbook-test
Bucket Name: verification-info
Bucket Name: zs-lavt-file
Bucket Name: zs-tax-file
Bucket Name: zstax-optdoc-videos
Bucket Name: zstax-tickets
Bucket Name: zstaxb
Bucket Name: zstaxdata
Bucket Name: zstaxfiles
*/

// check_point 检查点 --> 断点续传
// strings.NewReader
// bytes.NewReader
//

// SetObjectMeta
// GetObjectMeta
// GetObjectDetailedMeta
// GetObjectTagging
// SetTagging：tag 分类处理

// 简单上传  --> 暂且不支持 Put 操作的 Option 设置, 一次 http 请求交互
// PutObject
// PutObjectFromFile
// 上传请求携带禁止覆盖同名 object 的参数，x-oss-forbid-overwrite = true
// 上传大量文件时，采用【随机前缀】，提升性能

// 分片上传
// UploadFile
// oss.WithCheckpoint
// Upload ID : 标识该次上传事件
// 分片上传过程中断后，如果使用同一个Upload ID重新上传所有Part，则会覆盖之前上传的同名Part
// 如果使用新的Upload ID重新上传所有Part，旧的Upload ID中的分片会作为【碎片】继续保留
// OSS不支持自动合并分片, 您需要通过调用CompleteMultipartUpload手动合并分片。

// oss.ObjectStorageClass：存储类型, 目前线上 bucket 均为 Standard
// 默认采用标准存储
/*
type StorageClassType string

const (
	// StorageStandard standard
	StorageStandard StorageClassType = "Standard"

	// StorageIA infrequent access
	StorageIA StorageClassType = "IA"

	// StorageArchive archive
	StorageArchive StorageClassType = "Archive"

	// StorageColdArchive cold archive
	StorageColdArchive StorageClassType = "ColdArchive"
)
*/

// AppendObject

// CopyObject --> 同一 bucket 内 obj1 --> obj2
// CopyObject --> 同一 bucket 内 obj1 --> obj2
// CopyFile --> multi_part copy ，同样 oss.CheckPoint
// oss.ObjectStorageClass("Archive")

// DeleteObject
// DeleteObjects
// 选项参数：oss.DeleteObjectsQuiet
//

// GetObject
// 			oss.Range()
//			oss.IfModifiedSince
// 			oss.IfUnmodifiedSince
// GetObjectToFile
// DownloadFile: multi_part get, os.Checkpoint(true, "")

// 暂且先不支持
// SelectObject: csv/json object --> select sql

//
// 基于生命周期规则
// 基于最后一次修改时间的生命周期规则
// 基于最后一次访问时间的生命周期规则
