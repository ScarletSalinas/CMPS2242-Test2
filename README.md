## Test #2: Improving the TCP Echo Server

### Usage
1. Start the server
    ```bash
     go run server.go 
    ```
2. Open another terminal and connect
   ```bash
    nc localhost 4000
    ```
3. Test Special commnds/inputs
   Open another terminal and connect
   ```bash
    # Connect to server
    nc localhost 5000
    
    # Test commands
    hello
    /time
    /echo testing
    /quit

    # Test special cases
    bye
    ```
   
### Advanced Options
1. Port configuration: Allows the port to be set with a --port flag when you are starting the server.
    ```bash
     go run server.go --port 5000
    ```
2. Log all messages from each client with --logdir flag: Saves each client’s messages to a separate file named after their IP address.
   ```bash
    go run server.go --logdir ./test_logs
    ```
   ## References
  - LLM: DeepSeek, for tutoring, learning necessary concepts, and for guidance when needed.
  

### Link to Video
[Watch demo here]()   
