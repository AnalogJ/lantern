package backend


//The API should not really know anything about the implementation of the proxy or how to retrieve data from the database
//This interceptor interface creates a logical boundary around all implementation details.
type Interface interface {

	ListenMessages()

	////The interceptor will create cdp compatible events that need to be broadcasted to the Websocket(s).
	////This is done via a Pub/Sub model.
	////Supported events are:
	//// - network.requestWillBeSent
	//// - network.responseReceived
	//// - network.dataReceived
	//// - network.loadingFinished
	//
	//Subscribe()
	//Unsubscribe()
	//
	////Commands
	//// - network.canClearBrowserCache
	//CmdNetworkCanClearBrowserCache()
	//// - network.clearBrowserCache
	//CmdNetworkclearBrowserCache()
	//// - network.getResponseBody
	//CmdNetworkGetResponseBody()
	//// - network.getCookies
	//CmdNetworkGetCookies()
	//// - page.getResourceTree
	//CmdPageGetResourceTree()
	//
	//
	////TODO: commands for hijacking / injecting changes into a pending request.
}