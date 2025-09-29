# \DefaultAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAlive**](DefaultAPI.md#GetAlive) | **Get** /alive | Get alive
[**GetDeploymentByUrl**](DefaultAPI.md#GetDeploymentByUrl) | **Get** /deployment/{url} | Get deployment by URL
[**PostDeployNew**](DefaultAPI.md#PostDeployNew) | **Post** /deploy/new | Post deploy new
[**PutDeployFiles**](DefaultAPI.md#PutDeployFiles) | **Put** /deploy/files | Put deploy files
[**PutUserRegister**](DefaultAPI.md#PutUserRegister) | **Put** /user/register | Put user register



## GetAlive

> HealthCheckOutputBody GetAlive(ctx).Execute()

Get alive

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetAlive(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetAlive``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAlive`: HealthCheckOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetAlive`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetAliveRequest struct via the builder pattern


### Return type

[**HealthCheckOutputBody**](HealthCheckOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDeploymentByUrl

> GetDeploymentOutputBody GetDeploymentByUrl(ctx, url).Execute()

Get deployment by URL

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	url := "url_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.GetDeploymentByUrl(context.Background(), url).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.GetDeploymentByUrl``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDeploymentByUrl`: GetDeploymentOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.GetDeploymentByUrl`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**url** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDeploymentByUrlRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetDeploymentOutputBody**](GetDeploymentOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PostDeployNew

> SuccessOutputBody PostDeployNew(ctx).DeploymentCreateInputBody(deploymentCreateInputBody).Execute()

Post deploy new

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	deploymentCreateInputBody := *openapiclient.NewDeploymentCreateInputBody(*openapiclient.NewUrl("Domain_example")) // DeploymentCreateInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PostDeployNew(context.Background()).DeploymentCreateInputBody(deploymentCreateInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PostDeployNew``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PostDeployNew`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PostDeployNew`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPostDeployNewRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deploymentCreateInputBody** | [**DeploymentCreateInputBody**](DeploymentCreateInputBody.md) |  | 

### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutDeployFiles

> SuccessOutputBody PutDeployFiles(ctx).Url(url).Contents(contents).KeepLeadingDirectories(keepLeadingDirectories).PreserveExistingFiles(preserveExistingFiles).Execute()

Put deploy files

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	url := "url_example" // string | 
	contents := os.NewFile(1234, "some_file") // *os.File |  (optional)
	keepLeadingDirectories := true // bool |  (optional)
	preserveExistingFiles := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PutDeployFiles(context.Background()).Url(url).Contents(contents).KeepLeadingDirectories(keepLeadingDirectories).PreserveExistingFiles(preserveExistingFiles).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PutDeployFiles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutDeployFiles`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PutDeployFiles`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPutDeployFilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **url** | **string** |  | 
 **contents** | ***os.File** |  | 
 **keepLeadingDirectories** | **bool** |  | 
 **preserveExistingFiles** | **bool** |  | 

### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutUserRegister

> SuccessOutputBody PutUserRegister(ctx).AddExternalUserInputBody(addExternalUserInputBody).Execute()

Put user register

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	addExternalUserInputBody := *openapiclient.NewAddExternalUserInputBody("ExternalUserSource_example") // AddExternalUserInputBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.PutUserRegister(context.Background()).AddExternalUserInputBody(addExternalUserInputBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.PutUserRegister``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutUserRegister`: SuccessOutputBody
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.PutUserRegister`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPutUserRegisterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addExternalUserInputBody** | [**AddExternalUserInputBody**](AddExternalUserInputBody.md) |  | 

### Return type

[**SuccessOutputBody**](SuccessOutputBody.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

