package vpc

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateFlowLog invokes the vpc.CreateFlowLog API synchronously
func (client *Client) CreateFlowLog(request *CreateFlowLogRequest) (response *CreateFlowLogResponse, err error) {
	response = CreateCreateFlowLogResponse()
	err = client.DoAction(request, response)
	return
}

// CreateFlowLogWithChan invokes the vpc.CreateFlowLog API asynchronously
func (client *Client) CreateFlowLogWithChan(request *CreateFlowLogRequest) (<-chan *CreateFlowLogResponse, <-chan error) {
	responseChan := make(chan *CreateFlowLogResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateFlowLog(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// CreateFlowLogWithCallback invokes the vpc.CreateFlowLog API asynchronously
func (client *Client) CreateFlowLogWithCallback(request *CreateFlowLogRequest, callback func(response *CreateFlowLogResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateFlowLogResponse
		var err error
		defer close(result)
		response, err = client.CreateFlowLog(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// CreateFlowLogRequest is the request struct for api CreateFlowLog
type CreateFlowLogRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer    `position:"Query" name:"ResourceOwnerId"`
	Description          string              `position:"Query" name:"Description"`
	ResourceGroupId      string              `position:"Query" name:"ResourceGroupId"`
	IpVersion            string              `position:"Query" name:"IpVersion"`
	Tag                  *[]CreateFlowLogTag `position:"Query" name:"Tag"  type:"Repeated"`
	ResourceId           string              `position:"Query" name:"ResourceId"`
	ProjectName          string              `position:"Query" name:"ProjectName"`
	LogStoreName         string              `position:"Query" name:"LogStoreName"`
	ResourceOwnerAccount string              `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string              `position:"Query" name:"OwnerAccount"`
	TrafficPath          *[]string           `position:"Query" name:"TrafficPath"  type:"Repeated"`
	AggregationInterval  requests.Integer    `position:"Query" name:"AggregationInterval"`
	OwnerId              requests.Integer    `position:"Query" name:"OwnerId"`
	ResourceType         string              `position:"Query" name:"ResourceType"`
	TrafficType          string              `position:"Query" name:"TrafficType"`
	FlowLogName          string              `position:"Query" name:"FlowLogName"`
}

// CreateFlowLogTag is a repeated param struct in CreateFlowLogRequest
type CreateFlowLogTag struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}

// CreateFlowLogResponse is the response struct for api CreateFlowLog
type CreateFlowLogResponse struct {
	*responses.BaseResponse
	RequestId       string `json:"RequestId" xml:"RequestId"`
	Success         string `json:"Success" xml:"Success"`
	FlowLogId       string `json:"FlowLogId" xml:"FlowLogId"`
	ResourceGroupId string `json:"ResourceGroupId" xml:"ResourceGroupId"`
}

// CreateCreateFlowLogRequest creates a request to invoke CreateFlowLog API
func CreateCreateFlowLogRequest() (request *CreateFlowLogRequest) {
	request = &CreateFlowLogRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "CreateFlowLog", "vpc", "openAPI")
	request.Method = requests.POST
	return
}

// CreateCreateFlowLogResponse creates a response to parse from CreateFlowLog response
func CreateCreateFlowLogResponse() (response *CreateFlowLogResponse) {
	response = &CreateFlowLogResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
