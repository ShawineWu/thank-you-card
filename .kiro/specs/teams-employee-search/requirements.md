# Requirements Document

## Introduction

本功能旨在将 Teams 员工信息搜索接口从 Mock 实现替换为真实的 Microsoft Graph API 实现。该接口供 Teams App 和 Web Dashboard 使用，通过 Azure AD 获取组织内员工信息，支持按姓名或邮箱搜索员工。

## Glossary

- **Graph_API_Service**: 使用 Microsoft Graph API 实现的员工服务，负责与 Azure AD 交互获取员工数据
- **Token_Manager**: 负责 OAuth2 Client Credentials 流程的 access token 获取、缓存和刷新
- **Employee_Search_Handler**: 处理 `/api/v1/employees/search` 端点的 HTTP 请求处理器
- **Employee**: 员工数据模型，包含 ID、AADID、Name、Email、Department 字段
- **Azure_AD**: Microsoft Azure Active Directory，组织的身份认证和用户目录服务

## Requirements

### Requirement 1: OAuth2 Token 管理

**User Story:** As a backend service, I want to obtain and cache OAuth2 access tokens, so that I can authenticate requests to Microsoft Graph API efficiently.

#### Acceptance Criteria

1. WHEN the Graph_API_Service initializes, THE Token_Manager SHALL validate that AZURE_TENANT_ID, AZURE_CLIENT_ID, and AZURE_CLIENT_SECRET environment variables are configured
2. WHEN a Graph API request is made and no valid token exists, THE Token_Manager SHALL request a new access token using OAuth2 Client Credentials flow from `https://login.microsoftonline.com/{tenant_id}/oauth2/v2.0/token`
3. WHEN a valid cached token exists and has not expired, THE Token_Manager SHALL return the cached token without making a new request
4. WHEN the cached token is within 60 seconds of expiry, THE Token_Manager SHALL proactively refresh the token
5. IF the token request fails, THEN THE Token_Manager SHALL return a descriptive error including the HTTP status code and error message

### Requirement 2: 员工搜索功能

**User Story:** As a user, I want to search for employees by name or email, so that I can find colleagues to send thank you cards to.

#### Acceptance Criteria

1. WHEN a search query is provided, THE Graph_API_Service SHALL search Azure AD users using the Microsoft Graph API endpoint `https://graph.microsoft.com/v1.0/users`
2. WHEN searching users, THE Graph_API_Service SHALL use `$filter` parameter with `startswith(displayName, '{query}') or startswith(mail, '{query}')` to match users
3. WHEN searching users, THE Graph_API_Service SHALL request only necessary fields using `$select=id,displayName,mail,department,jobTitle`
4. WHEN search results are returned, THE Graph_API_Service SHALL map Graph API response fields to Employee model: id→AADID/ID, displayName→Name, mail→Email, department→Department
5. WHEN the search query is empty, THE Employee_Search_Handler SHALL return a 400 Bad Request error with message "Search query is required"
6. IF the Graph API request fails, THEN THE Graph_API_Service SHALL return an error with the failure reason

### Requirement 3: 员工验证功能

**User Story:** As a backend service, I want to validate that employee IDs correspond to valid Azure AD users, so that I can ensure cards are sent to real employees.

#### Acceptance Criteria

1. WHEN validating employee IDs, THE Graph_API_Service SHALL query each employee by their Azure AD Object ID using `https://graph.microsoft.com/v1.0/users/{id}`
2. WHEN all employee IDs are found in Azure AD, THE Graph_API_Service SHALL return true
3. WHEN any employee ID is not found in Azure AD (404 response), THE Graph_API_Service SHALL return false
4. IF the validation request fails due to network or authentication errors, THEN THE Graph_API_Service SHALL return an error

### Requirement 4: 错误处理和日志

**User Story:** As a system operator, I want comprehensive error handling and logging, so that I can diagnose issues with the employee search service.

#### Acceptance Criteria

1. WHEN a Graph API request fails, THE Graph_API_Service SHALL log the error with request URL, HTTP status code, and error details
2. WHEN a token request fails, THE Token_Manager SHALL log the failure with tenant ID and error message (excluding secrets)
3. WHEN an employee search is performed, THE Graph_API_Service SHALL log the search query and result count
4. IF rate limiting occurs (429 response), THEN THE Graph_API_Service SHALL log the retry-after header value and return an appropriate error

### Requirement 5: 服务配置和切换

**User Story:** As a developer, I want to easily switch between Mock and Graph API implementations, so that I can develop locally without Azure AD credentials.

#### Acceptance Criteria

1. WHEN AZURE_TENANT_ID, AZURE_CLIENT_ID, and AZURE_CLIENT_SECRET are all configured, THE application SHALL use Graph_API_Service
2. WHEN any of the Azure credentials are missing, THE application SHALL fall back to MockEmployeeService
3. WHEN the service implementation is selected, THE application SHALL log which implementation is being used
