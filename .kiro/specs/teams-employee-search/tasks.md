# Implementation Plan: Teams Employee Search

## Overview

本实现计划将 Teams 员工搜索接口从 Mock 实现替换为 Microsoft Graph API 实现。实现将增强现有的 `GraphEmployeeService`，添加线程安全的 token 管理、完善的错误处理和服务工厂函数。

## Tasks

- [x] 1. 增强 GraphEmployeeService Token 管理
  - [x] 1.1 添加 sync.RWMutex 保护 token 并发访问
    - 在 GraphEmployeeService 结构体中添加 `mu sync.RWMutex` 字段
    - 修改 getAccessToken() 方法使用读写锁保护 token 访问
    - _Requirements: 1.2, 1.3_
  
  - [x] 1.2 实现 token 主动刷新逻辑
    - 修改 token 过期检查，当 token 在 60 秒内过期时触发刷新
    - 确保刷新时使用写锁
    - _Requirements: 1.4_
  
  - [x] 1.3 增强 token 请求错误处理
    - 解析 token 请求的错误响应体
    - 返回包含 HTTP 状态码和错误消息的描述性错误
    - 添加日志记录（不记录 client_secret）
    - _Requirements: 1.5, 4.2_
  
  - [ ]* 1.4 编写 Token 管理属性测试
    - **Property 1: Token Caching Returns Cached Token**
    - **Property 2: Token Proactive Refresh Near Expiry**
    - **Property 3: Token Error Contains Status and Message**
    - **Validates: Requirements 1.3, 1.4, 1.5**

- [x] 2. 完善 SearchEmployees 实现
  - [x] 2.1 优化搜索 URL 构建
    - 确保 $filter 参数正确使用 startswith 条件
    - 确保 $select 参数包含必要字段
    - 正确处理查询字符串的 URL 编码
    - _Requirements: 2.1, 2.2, 2.3_
  
  - [x] 2.2 增强响应映射和错误处理
    - 确保 Graph API 响应正确映射到 Employee 模型
    - 添加搜索日志记录（查询和结果数量）
    - 处理 Graph API 错误并返回描述性错误
    - _Requirements: 2.4, 2.6, 4.1, 4.3_
  
  - [ ]* 2.3 编写搜索功能属性测试
    - **Property 4: Search URL Construction**
    - **Property 5: Graph Response Mapping**
    - **Property 6: Search Error Propagation**
    - **Validates: Requirements 2.1, 2.2, 2.3, 2.4, 2.6**

- [x] 3. Checkpoint - 确保 Token 和搜索功能测试通过
  - Ensure all tests pass, ask the user if questions arise.

- [x] 4. 完善 ValidateEmployees 实现
  - [x] 4.1 增强员工验证逻辑
    - 确保对每个 ID 调用正确的 Graph API 端点
    - 正确处理 404 响应（返回 false）
    - 正确处理网络/认证错误（返回 error）
    - 添加验证日志记录
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 4.1_
  
  - [ ]* 4.2 编写验证功能属性测试
    - **Property 7: Validation Correctness**
    - **Property 8: Validation Error Handling**
    - **Validates: Requirements 3.2, 3.3, 3.4**

- [x] 5. 实现服务工厂函数
  - [x] 5.1 创建 NewEmployeeService 工厂函数
    - 检查 AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET 配置
    - 所有配置存在时返回 GraphEmployeeService
    - 任一配置缺失时返回 MockEmployeeService
    - 添加日志记录使用的服务类型
    - _Requirements: 5.1, 5.2, 5.3_
  
  - [ ]* 5.2 编写服务选择属性测试
    - **Property 9: Service Selection Based on Configuration**
    - **Validates: Requirements 5.1, 5.2**

- [x] 6. 集成和配置更新
  - [x] 6.1 更新应用启动代码使用工厂函数
    - 修改 main.go 或服务初始化代码
    - 使用 NewEmployeeService 替代直接创建服务实例
    - _Requirements: 5.1, 5.2_
  
  - [x] 6.2 更新环境变量文档
    - 在 .env.example 或 README 中记录 Azure 配置变量
    - _Requirements: 1.1_

- [x] 7. Final Checkpoint - 确保所有测试通过
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- 现有的 GraphEmployeeService 已有基础框架，主要工作是增强和完善
- 测试使用 `rapid` 库进行属性测试，`httpmock` 进行 HTTP 模拟
- 所有 token 操作需要线程安全，使用 sync.RWMutex 保护
