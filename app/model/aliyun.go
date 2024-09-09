package model

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func SendSMS(accessKeyId, accessKeySecret, signature, templateCode, phoneNumber, code string) error {
	client, err := sdk.NewClientWithAccessKey("default", accessKeyId, accessKeySecret)
	if err != nil {
		return fmt.Errorf("failed to create Alibaba Cloud client: %v", err)
	} //dk.NewClientWithAccessKey 函数是阿里云 SDK（Software Development Kit）中的一个函数，用于创建一个基于访问密钥（Access Key）的客户端对象。
	//阿里云提供了多个服务（例如短信服务、对象存储服务、云服务器等），开发者可以使用相应的 SDK 来访问和管理这些服务。在使用 SDK 之前，需要创建一个客户端对象，该客户端对象将用于与阿里云服务进行通信。
	//NewClientWithAccessKey 函数的作用是创建一个使用访问密钥进行身份验证的客户端对象。访问密钥由阿里云控制台提供，包括 AccessKeyId 和 AccessKeySecret。这些密钥用于对 API 请求进行签名和身份验证，以确保请求的安全性和合法性。
	//通过创建客户端对象，您可以使用相应的 SDK 提供的方法来调用阿里云服务的 API，执行各种操作，例如发送短信、上传文件、创建云服务器等。

	// 设置验证码的有效期为五分钟
	//expiration := time.Now().Add(5 * time.Minute).Format("2006-01-02T15:04:05Z")

	//根据您提供的代码 request := requests.NewCommonRequest()，这段代码创建了一个通用的请求对象 (CommonRequest)。
	//在阿里云 SDK 中，通用请求对象 (CommonRequest) 是用于发送各种类型的请求的通用数据结构。它可以用来发送不同 API 的请求，包括但不限于短信服务、对象存储服务、云服务器等
	//通过 NewCommonRequest() 函数创建一个新的通用请求对象，您可以设置请求的参数、方法、域名、版本等信息，然后使用该对象发送请求。
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"

	request.QueryParams["PhoneNumbers"] = phoneNumber
	request.QueryParams["SignName"] = signature
	request.QueryParams["TemplateCode"] = templateCode
	request.QueryParams["TemplateParam"] = `{"code": "` + code + `"}`

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	fmt.Println("SMS sent successfully!")
	fmt.Println("Response:", response.GetHttpContentString())
	return nil
}
