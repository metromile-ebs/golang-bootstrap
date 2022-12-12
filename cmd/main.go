package main

import "metromile-ebs/streamline-graph-manager/internal/apiservice"

func main() {

	graphManagerApi := apiservice.NewgraphManagerService()
	graphManagerApi.Start()
}
