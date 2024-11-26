
---

## **Receipt Processor**

A web service that processes receipts and awards points based on predefined rules. Points are calculated using a set of criteria, and the service provides endpoints to process receipts and fetch points.

---

### **Endpoints**

1. **Process Receipts**  
   **Path**: `/receipts/process`  
   **Method**: `POST`  
   **Request**:
   ```json
   {
     "retailer": "Target",
     "purchaseDate": "2022-01-01",
     "purchaseTime": "13:01",
     "items": [
       { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
       { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
     ],
     "total": "18.74"
   }
   ```
   **Response**:
   ```json
   { "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
   ```

2. **Get Points**  
   **Path**: `/receipts/{id}/points`  
   **Method**: `GET`  
   **Response**:
   ```json
   { "points": 32 }
   ```

3. **Get All Receipts**  
   **Path**: `/receipts/all`  
   **Method**: `GET`  
   **Query Parameters**:
   - `limit`: Number of receipts to return (default: 10).
   - `offset`: Starting index for pagination (default: 0).  
   **Response**:
   ```json
   {
     "receipts": {
       "1": 32,
       "2": 45
     },
     "total": 11,
     "limit": 2,
     "offset": 0,
     "currentPage": 1,
     "totalPages": 6
   }
   ```

---

### **Running the Application**

#### **Using Docker**
1. **Build the Docker Image**:
   ```bash
   docker build -t receipt-processor .
   ```

2. **Run the Docker Container**:
   ```bash
   docker run -p 8080:8080 receipt-processor
   ```

3. **Access the Application**:
   - Use tools like `curl`, Postman, or your browser to access:
     - `http://localhost:8080/receipts/process`
     - `http://localhost:8080/receipts/{id}/points`
     - `http://localhost:8080/receipts/all`

---

#### **Using Air for Development**
`air` is used for live reloading during development.

1. **Install Air**:
   ```bash
   go install github.com/air-verse/air@latest
   ```

2. **Run the Application with Air**:
   ```bash
   air
   ```

3. **Modify Code and Watch Changes**:
   - `air` automatically reloads the application whenever changes are made.

---

### **Project Structure**

```
receipt-processor/
├── cmd/
│   └── main.go                  # Entry point of the application
├── internal/
│   ├── handlers/
│   │   ├── process.go           # Handle receipt processing
│   │   ├── get_points.go        # Handle fetching points for a receipt
│   │   └── get_all_receipts.go  # Handle fetching all receipts
│   ├── models/
│   │   └── receipt.go           # Receipt data model
│   └── store/
│       └── store.go             # In-memory storage
├── go.mod                       # Go module definition
├── go.sum                       # Dependency lock file
├── Dockerfile                   # Dockerfile for containerizing the app
└── README.md                    # Project documentation
```

---

### **Example Usage**

#### **Process a Receipt**
```bash
curl -X POST http://localhost:8080/receipts/process \
-H "Content-Type: application/json" \
-d '{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
        { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
        { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
    ],
    "total": "18.74"
}'
```

**Response**:
```json
{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
```

#### **Get Points for a Receipt**
```bash
curl http://localhost:8080/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points
```

**Response**:
```json
{ "points": 28 }
```

#### **Get All Receipts**
```bash
curl "http://localhost:8080/receipts/all?limit=2&offset=0"
```

**Response**:
```json
{
  "receipts": {
    "1": 32,
    "2": 45
  },
  "total": 11,
  "limit": 2,
  "offset": 0,
  "currentPage": 1,
  "totalPages": 6
}
```

---

### **Notes**
- The application stores data in memory and does not persist data between restarts.
- Make sure Docker and Go are installed to run the application.

---
