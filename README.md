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
- [Practical Go Lessons by Maximilien Andile](https://www.practical-go-lessons.com/)
- LLM: DeepSeek, for tutoring, learning necessary concepts, and for guidance when needed.
- [W3Schools](https://www.w3schools.com/go/go_switch.php)


  ### Most educationally enriching functionality
  I was most looking forward to the special responses, originally I was just implementing this functionality with if else statements but I revised it multiple times until I ended up at switch statements
  
  ### Functionality requiring the most research
  Definitely the client message logging, had to do lots of research to understand the requirements and formulate the steps required to implement the functionality.
      

### Link to Video
[Watch demo here](https://www.youtube.com/watch?v=-pOBUj6WGMM)   
