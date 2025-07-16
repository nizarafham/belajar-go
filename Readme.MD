# Uber Eats Clone API Documentation

## Base URL
```
https://your-domain.com/api
```

## Authentication
The API uses JWT (JSON Web Token) for authentication. For protected routes, include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format
All responses follow this general structure:

### Success Response
```json
{
  "data": {...},
  "message": "Success message (optional)"
}
```

### Error Response
```json
{
  "error": "Error message"
}
```

---

## 🔐 Authentication Endpoints

### Register User
Creates a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone_number": "+628123456789",
  "password": "securepassword123",
  "role": "user" // optional: "user", "tenant_owner", "driver"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone_number": "+628123456789",
  "role": "user",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Login User
Authenticates a user and returns a JWT token.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

---

## 👤 User Endpoints

### Get User Profile
Get current user profile information.

**Endpoint:** `GET /v1/users/me`

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
{
  "id": 1,
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone_number": "+628123456789",
  "role": "user",
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

## 📍 Location Endpoints

### Get All Locations
Retrieve all available locations.

**Endpoint:** `GET /locations`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Food Court Plaza",
    "type": "food_court",
    "operating_hours": "08:00-22:00",
    "image_url": "https://example.com/plaza.jpg"
  },
  {
    "id": 2,
    "name": "Street Food Area",
    "type": "street_food",
    "operating_hours": "18:00-02:00",
    "image_url": "https://example.com/street.jpg"
  }
]
```

### Get Tenants by Location
Get all tenants/restaurants in a specific location.

**Endpoint:** `GET /locations/:id/tenants`

**Parameters:**
- `id` (path): Location ID

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "location_id": 1,
    "name": "Pizza Palace",
    "description": "Authentic Italian pizza",
    "logo_url": "https://example.com/pizza-logo.jpg",
    "is_open": true,
    "xendit_subaccount_id": "sub_123"
  },
  {
    "id": 2,
    "location_id": 1,
    "name": "Burger House",
    "description": "Gourmet burgers",
    "logo_url": "https://example.com/burger-logo.jpg",
    "is_open": false
  }
]
```

---

## 🏪 Tenant Endpoints

### Get Tenant Details
Get detailed information about a specific tenant.

**Endpoint:** `GET /tenants/:id`

**Parameters:**
- `id` (path): Tenant ID

**Response (200 OK):**
```json
{
  "id": 1,
  "location_id": 1,
  "name": "Pizza Palace",
  "description": "Authentic Italian pizza with fresh ingredients",
  "logo_url": "https://example.com/pizza-logo.jpg",
  "is_open": true,
  "xendit_subaccount_id": "sub_123"
}
```

### Get Tenant Menu
Get all menu items for a specific tenant.

**Endpoint:** `GET /tenants/:id/menus`

**Parameters:**
- `id` (path): Tenant ID

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "tenant_id": 1,
    "name": "Margherita Pizza",
    "description": "Classic pizza with tomato and mozzarella",
    "price": 85000,
    "category": "pizza",
    "image_url": "https://example.com/margherita.jpg",
    "is_available": true
  },
  {
    "id": 2,
    "tenant_id": 1,
    "name": "Pepperoni Pizza",
    "description": "Pizza with pepperoni and cheese",
    "price": 95000,
    "category": "pizza",
    "image_url": "https://example.com/pepperoni.jpg",
    "is_available": true
  }
]
```

---

## 🍕 Menu Endpoints

### Get Menu Details
Get detailed information about a specific menu item.

**Endpoint:** `GET /menus/:id`

**Parameters:**
- `id` (path): Menu ID

**Response (200 OK):**
```json
{
  "id": 1,
  "tenant_id": 1,
  "name": "Margherita Pizza",
  "description": "Classic pizza with tomato sauce, mozzarella cheese, and fresh basil",
  "price": 85000,
  "category": "pizza",
  "image_url": "https://example.com/margherita.jpg",
  "is_available": true
}
```

---

## 🛍️ Order Endpoints

### Create Order
Create a new order with multiple items.

**Endpoint:** `POST /v1/orders`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "tenant_id": 1,
  "order_type": "delivery", // or "pickup"
  "delivery_address": "Jl. Sudirman No. 123, Jakarta",
  "items": [
    {
      "menu_id": 1,
      "quantity": 2
    },
    {
      "menu_id": 2,
      "quantity": 1
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "message": "Pesanan berhasil dibuat, silakan lakukan pembayaran.",
  "order_id": 1,
  "total_price": 265000,
  "payment_url": "https://checkout.xendit.co/web/xyz123"
}
```

### Get Order History
Get all orders for the current user.

**Endpoint:** `GET /v1/orders`

**Headers:** `Authorization: Bearer <token>`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "tenant_id": 1,
    "total_price": 265000,
    "order_type": "delivery",
    "status": "pending_payment",
    "delivery_address": "Jl. Sudirman No. 123, Jakarta",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "tenant_info": {
      "id": 1,
      "name": "Pizza Palace",
      "logo_url": "https://example.com/pizza-logo.jpg"
    }
  }
]
```

### Get Order Details
Get detailed information about a specific order.

**Endpoint:** `GET /v1/orders/:id`

**Headers:** `Authorization: Bearer <token>`

**Parameters:**
- `id` (path): Order ID

**Response (200 OK):**
```json
{
  "id": 1,
  "user_id": 1,
  "tenant_id": 1,
  "total_price": 265000,
  "order_type": "delivery",
  "status": "pending_payment",
  "delivery_address": "Jl. Sudirman No. 123, Jakarta",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "tenant_info": {
    "id": 1,
    "name": "Pizza Palace",
    "description": "Authentic Italian pizza",
    "logo_url": "https://example.com/pizza-logo.jpg",
    "is_open": true
  },
  "order_items": [
    {
      "id": 1,
      "order_id": 1,
      "menu_id": 1,
      "quantity": 2,
      "price_per_item": 85000,
      "menu_info": {
        "id": 1,
        "name": "Margherita Pizza",
        "description": "Classic pizza with tomato and mozzarella",
        "price": 85000,
        "category": "pizza",
        "image_url": "https://example.com/margherita.jpg"
      }
    },
    {
      "id": 2,
      "order_id": 1,
      "menu_id": 2,
      "quantity": 1,
      "price_per_item": 95000,
      "menu_info": {
        "id": 2,
        "name": "Pepperoni Pizza",
        "description": "Pizza with pepperoni and cheese",
        "price": 95000,
        "category": "pizza",
        "image_url": "https://example.com/pepperoni.jpg"
      }
    }
  ]
}
```

---

## 💳 Payment Endpoints

### Xendit Payment Webhook
This endpoint is used by Xendit to notify payment status changes.

**Endpoint:** `POST /payments/xendit-notification`

**Headers:** 
- `X-CALLBACK-TOKEN`: Xendit verification token

**Request Body (from Xendit):**
```json
{
  "id": "invoice_123",
  "external_id": "1",
  "user_id": "user_123",
  "status": "PAID",
  "merchant_name": "Your App Name",
  "amount": 265000,
  "paid_amount": 265000,
  "bank_code": "BCA",
  "paid_at": "2024-01-15T11:00:00Z"
}
```

**Response (200 OK):**
```json
{
  "status": "notifikasi diterima"
}
```

---

## 🔧 Utility Endpoints

### Health Check
Check if the API is running.

**Endpoint:** `GET /ping`

**Response (200 OK):**
```
pong
```

---

## 📊 Data Models

### User Roles
- `user`: Regular customer
- `tenant_owner`: Restaurant owner
- `driver`: Delivery driver

### Order Types
- `delivery`: Order will be delivered
- `pickup`: Customer will pick up the order

### Order Status
- `pending_payment`: Order created, waiting for payment
- `paid`: Payment completed
- `preparing`: Restaurant is preparing the order
- `ready`: Order is ready for pickup/delivery
- `on_the_way`: Order is being delivered
- `delivered`: Order has been delivered
- `cancelled`: Order has been cancelled

---

## 🚨 Error Codes

### Common HTTP Status Codes

| Code | Description |
|------|-------------|
| 200  | Success |
| 201  | Created |
| 400  | Bad Request - Invalid request format |
| 401  | Unauthorized - Invalid or missing token |
| 403  | Forbidden - Access denied |
| 404  | Not Found - Resource not found |
| 500  | Internal Server Error |

### Common Error Messages

```json
{
  "error": "Request tidak valid"
}
```

```json
{
  "error": "Semua field wajib diisi"
}
```

```json
{
  "error": "Email atau password salah"
}
```

```json
{
  "error": "Token tidak valid atau sudah kadaluarsa"
}
```

---

## 🛡️ Security Notes

1. **JWT Tokens**: Tokens expire after 72 hours
2. **Password Security**: Passwords are hashed using bcrypt
3. **CORS**: Configured to allow all origins (change for production)
4. **Webhook Security**: Xendit webhooks are verified using callback tokens

---

## 📝 Usage Examples

### Frontend Integration Example (JavaScript)

```javascript
// Login and store token
const login = async (email, password) => {
  const response = await fetch('/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password })
  });
  
  const data = await response.json();
  if (response.ok) {
    localStorage.setItem('token', data.token);
    return data;
  }
  throw new Error(data.error);
};

// Make authenticated requests
const getProfile = async () => {
  const token = localStorage.getItem('token');
  const response = await fetch('/api/v1/users/me', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return await response.json();
};

// Create order
const createOrder = async (orderData) => {
  const token = localStorage.getItem('token');
  const response = await fetch('/api/v1/orders', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify(orderData)
  });
  
  return await response.json();
};
```

---

## 📱 Mobile App Integration

This API is designed to work with mobile applications. Key considerations:

1. **Token Management**: Store JWT tokens securely
2. **Offline Support**: Cache location and menu data
3. **Push Notifications**: Implement for order status updates
4. **Error Handling**: Graceful handling of network errors

