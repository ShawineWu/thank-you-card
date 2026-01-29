# 新应用接入认证系统指南

## 📋 快速检查清单

- [ ] Hydra 服务器端配置
- [ ] 前端环境配置文件
- [ ] 安装依赖包
- [ ] 创建 AuthService 实例
- [ ] 配置路由拦截
- [ ] 实现回调页面
- [ ] 配置 TypeScript 类型
- [ ] 创建 HTTP 请求拦截器
- [ ] 测试登录流程
- [ ] 各环境部署验证

---

## 🔧 详细接入步骤

### 步骤 1: Hydra 服务器端配置

> ⚠️ **需要后端/运维协助**

在 Hydra 服务器注册新的 OAuth 2.0 客户端。

#### 1.1 注册新客户端

```bash
# 使用 Hydra CLI 创建客户端
hydra clients create \
  --endpoint https://hydra-test.castlery.com \
  --id my-new-app-test \
  --name "My New Application (Test)" \
  --grant-types authorization_code,refresh_token \
  --response-types code \
  --scope openid,offline,offline_access \
  --callbacks http://localhost:3001/callback,https://my-new-app-test.castlery.com/callback \
  --post-logout-callbacks http://localhost:3001/,https://my-new-app-test.castlery.com/
```

#### 1.2 配置各环境的客户端

需要在不同环境创建对应的客户端：

| 环境 | Client ID | Callback URL 示例 |
|------|-----------|-------------------|
| Local | `my-new-app-local` | `http://localhost:3001/callback` |
| Dev | `my-new-app-dev` | `https://my-new-app-dev.castlery.com/callback` |
| Test | `my-new-app-test` | `https://my-new-app-test.castlery.com/callback` |
| Prod | `my-new-app-prod` | `https://my-new-app.castlery.com/callback` |

> 💡 **注意事项**:
> - Callback URL 必须精确匹配，包括协议、域名、端口、路径
> - 可以配置多个回调 URL（开发环境 + 生产环境）
> - `offline_access` scope 允许获取 refresh_token

---

### 步骤 2: 安装依赖包

#### 2.1 安装 OIDC 客户端库

```bash
npm install oidc-client-ts@3.2.1
# 或
yarn add oidc-client-ts@3.2.1
# 或
pnpm add oidc-client-ts@3.2.1
```

#### 2.2 如果使用 Cowman 的共享库

```bash
# 安装 Cowman utils 包（已包含 AuthService）
pnpm add @cowman/utils
```

---

### 步骤 3: 创建环境配置文件

#### 3.1 创建配置文件

在项目 `src` 目录下创建各环境配置：

**文件结构**:
```
src/
  ├── config.local.json
  ├── config.dev.json
  ├── config.test.json
  └── config.prod.json
```

#### 3.2 配置内容

**`config.local.json`** (开发环境):

```json
{
  "REACT_APP_ENV": "local",
  "AUTH_SERVER": "https://hydra-test.castlery.com",
  "CLIENT_ID": "my-new-app-local",
  "LOGIN_CALLBACK_URL": "http://localhost:3001/callback",
  "LOGOUT_CALLBACK_URL": "http://localhost:3001/",
  "PLAT_UMS_SERVER": "https://plat-ums-api-test.cslr.io"
}
```

**`config.prod.json`** (生产环境):

```json
{
  "REACT_APP_ENV": "prod",
  "AUTH_SERVER": "https://hydra.castlery.com",
  "CLIENT_ID": "my-new-app-prod",
  "LOGIN_CALLBACK_URL": "https://my-new-app.castlery.com/callback",
  "LOGOUT_CALLBACK_URL": "https://my-new-app.castlery.com/",
  "PLAT_UMS_SERVER": "https://plat-ums-api.cslr.io"
}
```

> 📝 **配置说明**:
> - `AUTH_SERVER`: Hydra 服务器地址（不同环境不同）
> - `CLIENT_ID`: 在 Hydra 注册的客户端 ID
> - `LOGIN_CALLBACK_URL`: 登录成功后的回调地址
> - `LOGOUT_CALLBACK_URL`: 登出后的重定向地址
> - `PLAT_UMS_SERVER`: 用户管理系统 API（获取用户信息和权限）

---

### 步骤 4: 配置 TypeScript 类型声明

#### 4.1 创建 `src/typings.d.ts`

```typescript
// src/typings.d.ts

// 全局配置变量声明
declare const AUTH_SERVER: string;
declare const CLIENT_ID: string;
declare const LOGIN_CALLBACK_URL: string;
declare const LOGOUT_CALLBACK_URL: string;
declare const PLAT_UMS_SERVER: string;
declare const REACT_APP_ENV: 'local' | 'dev' | 'test' | 'prod';

// Window 对象扩展
interface Window {
  // 如果需要在 window 对象上挂载配置
  AUTH_SERVER?: string;
  CLIENT_ID?: string;
}
```

---

### 步骤 5: 实现配置加载逻辑

#### 5.1 创建 `src/setupConfig.ts`

```typescript
// src/setupConfig.ts

type Env = 'local' | 'dev' | 'test' | 'prod';

export async function setupConfig(): Promise<void> {
  // 根据域名判断环境
  const hostnameToEnv: Record<string, Env> = {
    'localhost': 'local',
    'my-new-app-dev.castlery.com': 'dev',
    'my-new-app-test.castlery.com': 'test',
    'my-new-app.castlery.com': 'prod',
  };

  const hostname = window.location.hostname;
  const env = hostnameToEnv[hostname] || 'local';

  // 动态导入对应环境的配置
  const configs = {
    local: () => import('./config.local.json'),
    dev: () => import('./config.dev.json'),
    test: () => import('./config.test.json'),
    prod: () => import('./config.prod.json'),
  };

  const { default: config } = await configs[env]();

  // 将配置注入到 window 对象（只读）
  Object.keys(config).forEach((key) => {
    Object.defineProperty(window, key, {
      value: config[key],
      configurable: false,
      enumerable: false,
      writable: false,
    });
  });
}
```

---

### 步骤 6: 创建或复用 AuthService

#### 选项 A: 使用 Cowman 的 AuthService

```typescript
// src/services/auth.ts
import { AuthService } from '@cowman/utils';

export async function initAuth() {
  return new AuthService({
    authServer: AUTH_SERVER,
    clientID: CLIENT_ID,
    loginCallbackUrl: LOGIN_CALLBACK_URL,
    logoutCallbackUrl: LOGOUT_CALLBACK_URL,
  });
}
```

#### 选项 B: 自己实现 AuthService

如果不使用 Cowman 的包，可以参考 [`AuthService.js`](file:///Users/youshanli/Workspaces/Programs/Castlery/cowman/packages/libs/utils/src/AuthService.js) 创建自己的实现：

```typescript
// src/services/AuthService.ts
import { UserManager, WebStorageStateStore, Log } from 'oidc-client-ts';

export class AuthService {
  private userManager: UserManager;
  private isLogining: boolean = false;
  private renewTokenPromise: Promise<any> | null = null;

  constructor({
    authServer,
    clientID,
    loginCallbackUrl,
    logoutCallbackUrl,
  }: {
    authServer: string;
    clientID: string;
    loginCallbackUrl: string;
    logoutCallbackUrl: string;
  }) {
    const settings = {
      authority: authServer,
      client_id: clientID,
      response_type: 'code',
      scope: 'offline_access offline openid',
      redirect_uri: loginCallbackUrl,
      post_logout_redirect_uri: logoutCallbackUrl,
      userStore: new WebStorageStateStore({ store: localStorage }),
      automaticSilentRenew: true,
    };

    this.userManager = new UserManager(settings);
    this.userManager.clearStaleState();

    Log.setLevel(Log.INFO);
    Log.setLogger(console);
  }

  async login() {
    if (this.isLogining) return;
    
    this.isLogining = true;
    const url = new URL(window.location.href);
    const args = { state: { url: url.href } };
    
    if (window.location.pathname !== '/callback') {
      return this.userManager.signinRedirect(args);
    }
  }

  async signinRedirectCallback() {
    const user = await this.userManager.signinRedirectCallback();
    this.isLogining = false;
    return user;
  }

  async renewToken() {
    if (this.renewTokenPromise) {
      return this.renewTokenPromise;
    }

    this.renewTokenPromise = this.userManager
      .signinSilent()
      .finally(() => {
        this.renewTokenPromise = null;
      });

    return this.renewTokenPromise;
  }

  async getUser() {
    return this.userManager.getUser();
  }

  async logout() {
    return this.userManager.signoutRedirect();
  }
}
```

---

### 步骤 7: 创建应用入口配置

#### 7.1 初始化配置和认证

```typescript
// src/app.ts 或 src/index.tsx

import { setupConfig } from './setupConfig';
import { initAuth } from './services/auth';
import { isTokenValid } from './utils/token'; // 实现 token 验证逻辑

// 初始化流程
const initConfigProcess = setupConfig().then(() => {
  return initAuth();
});

export { initConfigProcess };
```

#### 7.2 实现 Token 验证工具

```typescript
// src/utils/token.ts

export function isTokenValid(user: any): boolean {
  if (!user || !user.access_token) {
    return false;
  }

  // 检查是否过期
  const now = Math.floor(Date.now() / 1000);
  const expiresAt = user.expires_at || 0;
  
  // 提前 5 分钟认为过期（留出刷新时间）
  return expiresAt > now + 300;
}
```

---

### 步骤 8: 实现路由守卫

#### 8.1 React Router 示例

```typescript
// src/routes/ProtectedRoute.tsx

import { useEffect, useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { initConfigProcess } from '../app';
import { isTokenValid } from '../utils/token';

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const location = useLocation();

  useEffect(() => {
    checkAuth();
  }, [location.pathname]);

  async function checkAuth() {
    try {
      const auth = await initConfigProcess;
      const user = await auth.getUser();

      if (!isTokenValid(user)) {
        // Token 过期，尝试刷新
        await auth.renewToken();
        setIsAuthenticated(true);
      } else {
        setIsAuthenticated(true);
      }
    } catch (error) {
      // Token 无效，需要登录
      const auth = await initConfigProcess;
      await auth.login();
    }
  }

  if (isAuthenticated === null) {
    return <div>Loading...</div>;
  }

  return isAuthenticated ? <>{children}</> : null;
}
```

#### 8.2 路由配置

```typescript
// src/App.tsx

import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ProtectedRoute } from './routes/ProtectedRoute';
import { CallbackPage } from './pages/Callback';
import { HomePage } from './pages/Home';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* 回调路由 - 不需要保护 */}
        <Route path="/callback" element={<CallbackPage />} />
        
        {/* 受保护的路由 */}
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <HomePage />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
```

---

### 步骤 9: 实现回调页面

#### 9.1 创建 Callback 页面

```typescript
// src/pages/Callback.tsx

import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { initConfigProcess } from '../app';

export function CallbackPage() {
  const navigate = useNavigate();

  useEffect(() => {
    handleCallback();
  }, []);

  async function handleCallback() {
    try {
      const auth = await initConfigProcess;
      const user = await auth.signinRedirectCallback();

      // 重定向到登录前的页面
      const redirectUrl = user?.state?.url || '/';
      const url = new URL(redirectUrl);
      
      // 避免循环重定向
      if (url.pathname !== '/callback') {
        navigate(url.pathname + url.search, { replace: true });
      } else {
        navigate('/', { replace: true });
      }
    } catch (error) {
      console.error('Callback error:', error);
      // 错误处理：重定向到首页或错误页
      navigate('/', { replace: true });
    }
  }

  return (
    <div style={{ textAlign: 'center', marginTop: '100px' }}>
      <h2>Processing login...</h2>
      <p>Please wait while we complete your authentication.</p>
    </div>
  );
}
```

---

### 步骤 10: 创建 HTTP 请求拦截器

#### 10.1 Axios 拦截器

```typescript
// src/utils/axios.ts

import axios, { AxiosInstance } from 'axios';
import { initConfigProcess } from '../app';

export class AuthorizedAxios {
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = axios.create({
      timeout: 30000,
    });

    // 请求拦截器 - 添加 token
    this.axiosInstance.interceptors.request.use(
      async (config) => {
        const auth = await initConfigProcess;
        const user = await auth.getUser();

        if (user?.access_token) {
          config.headers.Authorization = `Bearer ${user.access_token}`;
        }

        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // 响应拦截器 - 处理 401 错误
    this.axiosInstance.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;

        // 401 错误且未重试过
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            const auth = await initConfigProcess;
            await auth.renewToken();

            // 重新获取 token 并重试请求
            const user = await auth.getUser();
            if (user?.access_token) {
              originalRequest.headers.Authorization = `Bearer ${user.access_token}`;
            }

            return this.axiosInstance(originalRequest);
          } catch (refreshError) {
            // 刷新失败，重新登录
            const auth = await initConfigProcess;
            await auth.login();
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      }
    );
  }

  getInstance(): AxiosInstance {
    return this.axiosInstance;
  }
}

// 导出单例
export const authorizedAxios = new AuthorizedAxios().getInstance();
```

#### 10.2 使用示例

```typescript
// src/services/user.ts

import { authorizedAxios } from '../utils/axios';

export async function getUserInfo() {
  const response = await authorizedAxios.get(
    `${PLAT_UMS_SERVER}/api/v1/user/info`
  );
  return response.data;
}
```

---

### 步骤 11: 实现登出功能

```typescript
// src/components/LogoutButton.tsx

import { initConfigProcess } from '../app';

export function LogoutButton() {
  async function handleLogout() {
    try {
      const auth = await initConfigProcess;
      await auth.logout();
    } catch (error) {
      console.error('Logout error:', error);
    }
  }

  return (
    <button onClick={handleLogout}>
      Logout
    </button>
  );
}
```

---

## 🧪 测试步骤

### 本地测试

#### 1. 启动应用

```bash
npm run dev
# 确保运行在配置的端口，例如 3001
```

#### 2. 访问受保护的页面

```
http://localhost:3001/
```

#### 3. 验证重定向流程

- [x] 应用重定向到 Hydra
- [x] Hydra 重定向到 Microsoft 登录页
- [x] 输入 Microsoft 凭据
- [x] 重定向回 `/callback`
- [x] 最终返回到原页面

#### 4. 检查存储

打开浏览器 DevTools > Application > Local Storage:

```javascript
// 应该看到类似的 key
oidc.user:https://hydra-test.castlery.com:my-new-app-local
```

#### 5. 测试 API 调用

```typescript
// 在控制台测试
import { authorizedAxios } from './utils/axios';

authorizedAxios.get('${PLAT_UMS_SERVER}/api/v1/user/info')
  .then(res => console.log(res.data));
```

检查 Network 面板，请求头应包含：

```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### 6. 测试 Token 刷新

```typescript
// 手动修改 localStorage 中的 expires_at 为过去的时间
// 然后重新加载页面，应该自动刷新 token
```

#### 7. 测试登出

点击登出按钮，应该：
- [x] 清除 localStorage
- [x] 重定向到 Hydra 登出端点
- [x] 最终返回到 `LOGOUT_CALLBACK_URL`

---

## 🚀 部署配置

### Webpack / Vite 配置

确保构建工具能够正确处理环境配置文件：

#### Webpack 示例

```javascript
// webpack.config.js

module.exports = {
  // ...
  plugins: [
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV),
    }),
  ],
};
```

#### Vite 示例

```typescript
// vite.config.ts

import { defineConfig } from 'vite';

export defineConfig({
  // 环境变量会自动加载
  // 确保 config.*.json 文件在 public 目录或动态导入
});
```

---

## 📝 配置检查清单

在部署到各环境前，确保完成以下配置：

### Hydra 服务器端

- [ ] 已在 Hydra 注册 OAuth 客户端（各环境）
- [ ] Callback URL 配置正确
- [ ] 已配置正确的 scopes
- [ ] 已启用 `authorization_code` 和 `refresh_token` grant types

### 前端应用

- [ ] 安装了 `oidc-client-ts` 依赖
- [ ] 创建了各环境的配置文件
- [ ] TypeScript 类型声明完整
- [ ] 实现了 AuthService 或使用 Cowman 的实现
- [ ] 配置了路由守卫
- [ ] 实现了 `/callback` 回调页面
- [ ] 配置了 HTTP 请求拦截器
- [ ] 实现了登出功能

### 域名和 HTTPS

- [ ] 生产环境使用 HTTPS
- [ ] 域名已备案（如需要）
- [ ] DNS 解析配置正确
- [ ] SSL 证书有效

### 环境变量

- [ ] Local: `http://localhost:xxxx/callback`
- [ ] Dev: `https://xxx-dev.castlery.com/callback`
- [ ] Test: `https://xxx-test.castlery.com/callback`
- [ ] Prod: `https://xxx.castlery.com/callback`

---

## ⚠️ 常见问题

### Q1: Redirect URI mismatch 错误

**现象**: 登录时报错 "The redirect URI is invalid"

**解决**:
1. 检查 Hydra 配置中的 callback URL 是否与前端配置完全一致
2. 注意协议（http/https）、域名、端口、路径必须完全匹配
3. 如果是本地开发，确保端口号一致

### Q2: 无限重定向循环

**现象**: 页面不断刷新，在登录和回调之间循环

**解决**:
1. 检查回调页面的重定向逻辑，避免重定向到 `/callback`
2. 确保 state 参数正确保存和恢复
3. 检查 localStorage 是否被正确清理

### Q3: Token 无法自动刷新

**现象**: Token 过期后用户被重新登录

**解决**:
1. 确保 `automaticSilentRenew: true`
2. 检查是否请求了 `offline_access` scope
3. 确保 Hydra 客户端配置允许 `refresh_token` grant type

### Q4: 跨域问题

**现象**: API 请求被 CORS 策略阻止

**解决**:
1. 后端 API 需要配置 CORS 允许前端域名
2. 确保请求头包含正确的 `Origin`
3. 生产环境使用相同的顶级域名（如 `*.castlery.com`）

---

## 📞 需要协助的团队

| 事项 | 联系人/团队 | 职责 |
|------|------------|------|
| Hydra 客户端注册 | 运维团队/后端团队 | 在 Hydra 创建 OAuth 客户端配置 |
| Microsoft SSO 配置 | IT 部门 | Azure AD 配置（如需新增应用） |
| UMS 权限配置 | 后端团队 | 配置应用权限和用户访问权限 |
| 域名和证书 | 运维团队 | DNS 解析、SSL 证书配置 |
| 部署流程 | DevOps | CI/CD 配置、环境变量设置 |

---

## 🎯 总结

接入认证系统的核心工作包括：

1. **后端配置** (需运维支持): 在 Hydra 注册 OAuth 客户端
2. **前端开发**: 集成 OIDC 客户端库，实现登录流程
3. **测试验证**: 完整测试登录、API 调用、刷新、登出
4. **部署配置**: 各环境配置文件和域名设置

遵循本指南，你可以快速将新应用接入到 Cowman 的认证体系中！
