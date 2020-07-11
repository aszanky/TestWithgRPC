#### Message Request and Response using gRPC

Run :
    1. API Task 1 First -> This is just for Request Response
    2. API Task 2 For Getting All Messages
    3. API Task 3 For streaming - long live connection
Run on client.go code what task you want to run, comment the rest, do like this:
- Running streaming, comment task 1 & task 2 on main function
- Running task 1 and Task 2, comment task 3
- Notes: Task 3 can't run independently, because the API need request message

To run client and server, just:
   - Please ```dep ensure``` first
   - on client filepath, type ```go run client.go```
   - on server filepath, type ```go run server.go```