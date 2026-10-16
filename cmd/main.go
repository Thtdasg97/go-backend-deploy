package main

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	app := NewApp()

	//
	defer func() {
		// khi Be tắt ->  thì cái này sẽ chạy, để tắt database
		app.entClient.Close()
		// setup close rabbit MQ
		app.rabbitmq.Conn.Close()
	}()

	app.Start()

}
