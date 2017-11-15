# Social API

# Features
* Manage user
* Manage user connection
* Block user update
* View list user message update

## Component
This API is written in GO Lang 1.9 and uses:
* Go, Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
* Iris, Iris is a fast, simple and efficient web framework for Go. Iris provides a beautifully expressive and easy to use foundation for your next website, API, or distributed app.
* Glide, Is Vendor Package Management for Golang. Glide is a tool for managing the vendor directory within a Go package.
* MySQL, MySQL is a freely available open source Relational Database Management System (RDBMS) that uses Structured Query Language (SQL)
* Docker, Docker is an open platform for developers and sysadmins to build, ship, and run distributed applications, whether on laptops, data center VMs, or the cloud.

# API routes

### User
This module is for handling user management
* **[Get]** /api/v1/user
Get user list

        param uri:
            name: string [wildcard]
            email: string [wildcard]
            status: int [0|1]

* **[Post]** /api/v1/user
Create new user

        json param:
        {
            "name": string
            "email": string
        }

* **[Put]** /api/v1/user/{id}
Update existing user

        path param:
            {id}: int

        json param:
        {
            "name": string
            "email": string
            "status": int [0/1]
        }

* **[Delete]** /api/v1/user/{id}
Delete existing user

        path param:
            {id}: int

### Connection
This module is for handling user connection
* **[Post]** /api/v1/connection
Create new user connection

        json param:
        {
            "friends": [
                string,
                string
            ]
        }

* **[Post]** /api/v1/connection/show
View user friend list

        json param:
        {
            "email": string
        }

* **[Post]** /api/v1/connection/common
View common friend between 2 users

        json param:
        {
            "friends": [
                string,
                string
            ]
        }

### Subscribe
This module is for handling user subscription
* **[Post]** /api/v1/subscribe
Create new subscription from a requestor user to target user

        json param:
        {
            "requestor": string,
            "target": string
        }
        
* **[Post]** /api/v1/unsubscribe
Remove subscription from a requestor user to target user

        json param:
        {
            "requestor": string,
            "target": string
        }

### Block
This module is for handling user block
* **[Post]** /api/v1/block
Block user update and friend connection from a requestor to target user

    json param:
    {
        "requestor": string,
        "target": string
    }

### Message
This module is for handling user updates
* **[Post]** /api/v1/message
Send update from a user to list followers

        json param:
        {
            "sender": string,
            "text": string
        }

# Installation
For running this application locally, you can use docker, simply run on the go project folder:

    docker-compose docker-compose up -d --build

After docker build is finished wait about 5-15 minutes for waiting go updating its dependencies depending you internet connection. Run bellow command to view go dependency installation:

    docker-compose logs -f social-api
    
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
        winpty docker exec -it test-social-api bash
        
        unix:
        docker exec -it test-social-api bash
        
    2. Make sure you are on the application folder (/go/src/social-api)
    
    3. Run glide
        glide install
        
    4. Wait until the process is finished
    
    5. Restart docker
    
    
# Unit Test Result
![N|Solid](http://image.ibb.co/j7Rysb/image.png)