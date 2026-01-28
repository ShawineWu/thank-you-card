# 认证功能设置说明

## 1. 安装依赖

在 `backend` 目录下运行：

```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

## 2. 环境变量

在 `backend/.env` 文件中添加（可选，有默认值）：

```env
JWT_SECRET=your-secret-key-change-in-production
```

## 3. 测试账户

系统会自动创建以下测试账户（密码都是 `password123`）：

- **employee** / password123 - 普通员工角色
- **hr** / password123 - HR 角色（可访问 Analytics）
- **admin** / password123 - 管理员角色（可访问 Analytics）

## 4. 数据库迁移

启动后端服务时，会自动：
- 创建 `users` 表
- 创建测试账户数据

## 5. API 端点

### 登录
```
POST /api/auth/login
Content-Type: application/json

{
  "username": "hr",
  "password": "password123"
}
```

响应：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 2,
    "username": "hr",
    "role": "HR",
    "employee": {
      "id": 2,
      "name": "Mock HR Admin",
      "email": "hradmin@example.com",
      "department": "HR"
    }
  }
}
```

### 使用 Token

在后续请求的 Header 中添加：
```
Authorization: Bearer <token>
```

## 6. 角色权限

- **EMPLOYEE**: 可以访问所有普通功能（Feed, Create Card, My Cards, Stats, Top 10, Values）
- **HR**: 可以访问所有功能 + HR Analytics
- **ADMIN**: 可以访问所有功能 + HR Analytics
