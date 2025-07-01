// To deploy to AWS Lambda, ensure you add these dependencies to your go.mod:
// github.com/aws/aws-lambda-go/lambda
// github.com/awslabs/aws-lambda-go-api-proxy/fiber

package main

import (
	"fitness-hack/internal/server"

	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
)

func main() {
	app := server.NewFiberApp() // You should have a function that returns your *fiber.App
	adapter := fiberadapter.New(app)
	lambda.Start(adapter.ProxyWithContext)
}
