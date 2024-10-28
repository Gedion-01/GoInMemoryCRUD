# Go CRUD Challenge

## Tech Stack
- **Go**: Programming language used for the server and in-memory database.
- **Fiber**: Web framework for building the API.
- **UUID**: Library for generating unique identifiers.

## API Endpoints
### Base URL: `http://localhost:5000/api/v1`

#### Create a Person
- **Endpoint**: `POST /person`
- **Description**: Creates a new person.
  
#### Get a Person
- **Endpoint**: `GET /person/:id`
- **Description**: Retrieves a person by ID.
  
#### Get All Persons
- **Endpoint**: `GET /person`
- **Description**: Retrieves all persons.
  
#### Update a Person
- **Endpoint**: `PUT /person/:id`
- **Description**: Updates a person by ID.
  
#### Delete a Person
- **Endpoint**: `DELETE /person/:id`
- **Description**: Deletes a person by ID.
  
## How to Use
### Prerequisites
- **Go**: Ensure Go is installed on your machine.
- **Fiber**: Ensure Fiber is installed.

### Running the Server
1. **Clone the Repository**:
```
   git clone https://github.com/Gedion-01/GoInMemoryCRUD
```
```
   cd go-crud-challenge
```
   
2. **Install Dependencies**:
```
   go mod tidy
```
   
4. **Run the Server**:
```
   make run
```

## Accessing the In-Memory Database

The in-memory database is accessed through the API endpoints. Each endpoint interacts with the in-memory database to perform CRUD (Create, Read, Update, Delete) operations.

## Connecting to the Database Using Telnet

You can connect to the in-memory database using `telnet` to access and query the database directly.

### Steps to Connect:

1. **Open a Terminal**:
   Open a new terminal window on your system.

2. **Connect Using Telnet**:
   Use the following command to connect to the database:

   ```
   telnet localhost 8080
   ```

## Available Commands

Once connected, you can use the following commands to interact with the database:

### Set a Person

Add a new person to the database.

```
set <name> <age> <hobbies>
```

Example:
```
set ghost 23 running,reading
```

### Get a Person

Retrieve information about a specific person using their unique ID.

```
get <id>
```

Example:
```
get unique-id
```

### Update a Person

Update the information of an existing person.

```
update <id> <name> <age> <hobbies>
```

Example:
```
update unique-id ghost 24 running,reading
```

### Delete a Person

Remove a person from the database using their unique ID.

```
delete <id>
```

Example:
```
delete unique-id
```

### Get All Persons

Retrieve information about all persons in the database.

```
all
```

### Exit the Connection

Close the Telnet connection to the database.

```
exit
```
