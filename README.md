# Authentication Overview

- User Signup
- User Login & Logout
- Token Management:
    1. Access Tokens (JWT-based authentication) - Short Lived, 
    2. Refresh Tokens (Session persistence) - Long Lived

# Session Management (Redis)

## UUID Mapping:
- Tokens mapped with token’s unique identifiers (UUID),
- Improved security and control over user sessions
## Redis because of 
- Fast in-memory storage
- Efficient handling of frequent reads/writes
- Optimal solution for token validation and session expiry

# Middlewares & Authorization
- Basic Middlewares:
    Logging, CORS, Gzip
- Custom Middlewares:
    1. Authorization Token verification
    2. JWT Parsing and Claims extraction
    3. Maintaining Current Context User
       
# Permission-Based ACL
- ACL Integration: Permissions assigned explicitly per user



<img width="578" height="327" alt="login" src="https://github.com/user-attachments/assets/7beb4028-a2d2-4399-8e8b-09d1f92b37db" />
<img width="533" height="311" alt="login2" src="https://github.com/user-attachments/assets/00eb488a-ffc9-4092-817a-df705de100e4" />
