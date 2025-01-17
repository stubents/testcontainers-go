package gcloud_test

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
)

func ExampleRunDatastoreContainer() {
	// runDatastoreContainer {
	ctx := context.Background()

	datastoreContainer, err := gcloud.RunDatastoreContainer(
		ctx,
		testcontainers.WithImage("gcr.io/google.com/cloudsdktool/cloud-sdk:367.0.0-emulators"),
		gcloud.WithProjectID("datastore-project"),
	)
	if err != nil {
		panic(err)
	}

	// Clean up the container
	defer func() {
		if err := datastoreContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()
	// }

	// datastoreClient {
	projectID := datastoreContainer.Settings.ProjectID

	options := []option.ClientOption{
		option.WithEndpoint(datastoreContainer.URI),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	}

	dsClient, err := datastore.NewClient(ctx, projectID, options...)
	if err != nil {
		panic(err)
	}
	defer dsClient.Close()
	// }

	type Task struct {
		Description string
	}

	k := datastore.NameKey("Task", "sample", nil)
	data := Task{
		Description: "my description",
	}
	_, err = dsClient.Put(ctx, k, &data)
	if err != nil {
		panic(err)
	}

	saved := Task{}
	err = dsClient.Get(ctx, k, &saved)
	if err != nil {
		panic(err)
	}

	fmt.Println(saved.Description)

	// Output:
	// my description
}
