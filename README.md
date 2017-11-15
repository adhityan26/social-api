# Social API
For demo open:
[Social API](ant-social-api.herokuapp.com "Social API on Heroku")

# Features
* Manage user
* Connect/remove remove friend connection user
* Subscribe/unsubscribe user update
* Block/unblock user update
* View list user message update

## Component
This API is written in GO Lang 1.9 and uses:
* Go, is an open source programming language that makes it easy to build simple, reliable, and efficient software.
* Iris, is a fast, simple and efficient web framework for Go. Iris provides a beautifully expressive and easy to use foundation for your next website, API, or distributed app.
* Glide, Is Vendor Package Management for Golang. Glide is a tool for managing the vendor directory within a Go package.
* MySQL, is a freely available open source Relational Database Management System (RDBMS) that uses Structured Query Language (SQL)
* Docker, is an open platform for developers and sysadmins to build, ship, and run distributed applications, whether on laptops, data center VMs, or the cloud.

# API routes

### User
This module is for handling user management
* **[Get]** /api/v1/user

  Get user list
    * request

            uri query param:
                name: string [wildcard],
                email: string [wildcard],
                status: int [0|1]
                
    * response success
    
            json:
            {
                "user: [user list],
                "success": true
            }
            response code:
                200
    
    * response error
            
            json:
            {
                "message": failure message
                "success": false
            }
            response code:
                404 -> user not found

* **[Post]** /api/v1/user

  Create new user

    * request
    
            json param:
            {
                "name": string,
                "email": string
            }
            
    * response success
    
            json:
            {
                "user": inserted user value,
                "success": true,
                "message": success message                
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                409 -> duplicate data found
                412 -> invalid format
                428 -> mandatory field not present
                500 -> internal system error
                

* **[Put]** /api/v1/user/{id}

  Update existing user
  
    * request

            path param:
                {id}: int
    
            json param:
            {
                "name": string
                "email": string
                "status": int [0/1]
            }
            
    * response success
    
            json:
            {
                "user": updated user value,
                "success": true,
                "message": success message                
            }
            response code:
                200
            
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                404 -> user not found
                412 -> invalid format
                428 -> mandatory field not present
                500 -> internal system error

* **[Delete]** /api/v1/user/{id}

  Delete existing user
    
    * request
    
            path param:
                {id}: int
            
    * response success
    
            json:
            {
                "success": true,
                "message": success message                
            }
            response code:
                200
            
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                404 -> user not found
                412 -> invalid format
                500 -> internal system error

### Connection
This module is for handling user connection
* **[Post]** /api/v1/connection

  Create new user connection
  
    * request

            json param:
            {
                "friends": [
                    string,
                    string
                ]
            }
            
    * response success
    
            json:
            {
                "success": true            
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                409 -> duplicate data found
                412 -> invalid format
                428 -> mandatory field not present
                500 -> internal system error
                
* **[Delete]** /api/v1/connection

  Delete user connection
  
    * request

            json param:
            {
                "friends": [
                    string,
                    string
                ]
            }
                    
    * response success
    
            json:
            {
                "success": true            
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                404 -> user not found
                412 -> invalid format
                500 -> internal system error
                
* **[Post]** /api/v1/connection/show

  View user friend list
  
    * request

            json param:
            {
                "email": string
            }
            
    * response success
    
            json:
            {
                "count": total user friend
                "friends": list of user friend
                "success": true            
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                412 -> invalid format
                428 -> mandatory field not present
                500 -> internal system error
                

* **[Post]** /api/v1/connection/common

  View common friend between 2 users
  
    * request

            json param:
            {
                "friends": [
                    string,
                    string
                ]
            }
            
    * response success
    
            json:
            {
                "count": total user friend
                "friends": list of user friend
                "success": true            
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                404 -> user not found
                412 -> invalid format
                428 -> mandatory field not present
                500 -> internal system error
                
### Subscribe
This module is for handling user subscription
* **[Post]** /api/v1/subscribe

  Create new subscription from a requestor user to target user
  
    * request

            json param:
            {
                "requestor": string,
                "target": string
            }
            
    * response success
    
            json:
            {
                "success": true            
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                409 -> duplicate data found
                412 -> invalid format
                428 -> mandatory field not present
                500 -> internal system error
        
* **[Delete]** /api/v1/subscribe

  Remove subscription from a requestor user to target user
  
    * request

            json param:
            {
                "requestor": string,
                "target": string
            }
                    
    * response success
    
            json:
            {
                "success": true            
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                404 -> user not found
                412 -> invalid format
                500 -> internal system error

### Block
This module is for handling user block
* **[Post]** /api/v1/block

  Block user update and friend connection from a requestor to target user
  
    * request

            json param:
            {
                "requestor": string,
                "target": string
            }
            
    * response success
    
            json:
            {
                "success": true            
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                409 -> duplicate data found
                412 -> invalid format
                428 -> mandatory field not present
                500 -> internal system error    
    
* **[Delete]** /api/v1/block

  Unblock user update and friend connection from a requestor to target user
 
    * request

            json param:
            {
                "requestor": string,
                "target": string
            }
                    
    * response success
    
            json:
            {
                "success": true            
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                404 -> user not found
                412 -> invalid format
                500 -> internal system error
        
### Message
This module is for handling user updates
* **[Post]** /api/v1/message

  Send update from a user to list followers
 
    * request

            json param:
            {
                "sender": string,
                "text": string
            }
            
    * response success
    
            json:
            {
                "success": true,
                "recipients": list of user that can receive update
            }
            response code:
                200
        
    * response error
    
            json:
            {
                "success": false,
                "message": failure message                
            }
            response code:
                404 -> user not found
                412 -> invalid format
                428 -> mandatory field not present
                500 -> internal system error

# Installation
For running this application locally, you can use docker, simply run on the go project folder:

    $ docker-compose docker-compose up -d --build

After docker build is finished wait about 5-15 minutes for waiting go updating its dependencies depending you internet connection. Run bellow command to view go dependency installation:

    $ docker-compose logs -f social-api
    
Before you can view the application on your browser (http://localhost:8383) wait until you see:

    test-social-api  | [00] Now listening on: http://localhost:8080
    test-social-api  | [00] Application started. Press CTRL+C to shut down.

Then press ctrl+c to close to logs view


### Troubleshot
If you see something like: 

    test-social-api  | [00] main.go:10:2: cannot find package "github.com/joho/godotenv" in any of:
    test-social-api  | [00]         /go/src/social-api/vendor/github.com/joho/godotenv (vendor tree)
    test-social-api  | [00]         /usr/local/go/src/github.com/joho/godotenv (from $GOROOT)
    test-social-api  | [00]         /go/src/github.com/joho/godotenv (from $GOPATH)

It means the dependency loader is failed to fetch all the dependency, you can manualy load the dependency by connecting to bash and execute glide

    1. Access to application container bash
        windows:
        $ winpty docker exec -it test-social-api bash
        
        unix:
        $ docker exec -it test-social-api bash
        
    2. Make sure you are on the go application folder (/go/src/social-api)
    
    3. Run glide install
        $ glide install
        
    4. Wait until the process is finished
    
    5. Restart docker
        $ docker-compose restart social-api
    
    
# Unit Test Result
![N|Solid](http://image.ibb.co/nk7UJR/image.png)

### Detail
Detail function coverage

* Block module

        social-api\apps\controller\block\block_controller.go:24:        Create          97.6%
        social-api\apps\controller\block\block_controller.go:104:       Remove          94.3%
        social-api\apps\controller\block\block_routes.go:14:            Handler         100.0%
        total:                                                          (statements)    96.3%

* Connection module

        social-api\apps\controller\connection\connection_controller.go:28:      Index                   100.0%
        social-api\apps\controller\connection\connection_controller.go:73:      Create                  90.7%
        social-api\apps\controller\connection\connection_controller.go:184:     CreateConnection        77.8%
        social-api\apps\controller\connection\connection_controller.go:203:     Remove                  88.9%
        social-api\apps\controller\connection\connection_controller.go:292:     RemoveConnection        50.0%
        social-api\apps\controller\connection\connection_controller.go:306:     Common                  100.0%
        social-api\apps\controller\connection\connection_routes.go:14:          Handler                 100.0%
        total:                                                                  (statements)            91.7%
    
* Landing module

        social-api\apps\controller\landing\landing_controller.go:14:    Index           100.0%
        social-api\apps\controller\landing\landing_routes.go:14:        Handler         100.0%
        total:                                                          (statements)    100.0%

* Message module
    
        social-api\apps\controller\message\message_controller.go:25:    Create          97.6%
        social-api\apps\controller\message\message_routes.go:14:        Handler         100.0%
        total:                                                          (statements)    97.8%
        
* Subscribe module
    
        social-api\apps\controller\subscribe\subscribe_controller.go:24:        Create          93.5%
        social-api\apps\controller\subscribe\subscribe_controller.go:115:       Remove          88.9%
        social-api\apps\controller\subscribe\subscribe_routes.go:14:            Handler         100.0%
        total:                                                                  (statements)    92.0%
    
* User module
    
        social-api\apps\controller\user\user_controller.go:16:  Index           100.0%
        social-api\apps\controller\user\user_controller.go:57:  Show            95.7%
        social-api\apps\controller\user\user_controller.go:104: Create          96.2%
        social-api\apps\controller\user\user_controller.go:157: Update          95.5%
        social-api\apps\controller\user\user_controller.go:210: Remove          93.8%
        social-api\apps\controller\user\user_routes.go:14:      Handler         100.0%
        total:                                                  (statements)    96.5%
    
For more even detailed coverage you can view on test-result folder on html format
