package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

func main() {
	r := gin.Default()
	r.GET("/os/bucket", listBuckets)
	r.POST("/os/bucket/:bucketName", createBucket)
	r.Run()
}

func listBuckets(ctx *gin.Context) {
	namespace := os.Getenv("NAMESPACE")
	compartmentId := os.Getenv("COMPARTMENT_ID")
	provider, err := auth.OkeWorkloadIdentityConfigurationProvider()
	if err != nil {
		panic(err)
	}
	client, clerr := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if clerr != nil {
		panic(clerr)
	}
	res, rerr := client.ListBuckets(
		ctx,
		objectstorage.ListBucketsRequest{
			NamespaceName: &namespace,
			CompartmentId: &compartmentId,
		},
	)
	if rerr != nil {
		panic(rerr)
	}
	ctx.JSON(200, gin.H{
		"data": res.Items,
	})
}

func createBucket(ctx *gin.Context) {
	namespace := os.Getenv("NAMESPACE")
	compartmentId := os.Getenv("COMPARTMENT_ID")
	bucketName := ctx.Param("bucketName")
	provider, err := auth.OkeWorkloadIdentityConfigurationProvider()
	if err != nil {
		panic(err)
	}
	client, clerr := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if clerr != nil {
		panic(clerr)
	}
	res, rerr := client.CreateBucket(
		ctx,
		objectstorage.CreateBucketRequest{
			NamespaceName: &namespace,
			CreateBucketDetails: objectstorage.CreateBucketDetails{
				Name:          &bucketName,
				CompartmentId: &compartmentId,
			},
		},
	)
	if rerr != nil {
		panic(rerr)
	}
	ctx.JSON(200, gin.H{
		"data": res.Bucket,
	})
}
