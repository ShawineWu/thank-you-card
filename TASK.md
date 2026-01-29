# Microsoft Teams App 项目任务清单

## 项目概述
构建 Microsoft Teams App，用于企业内部感谢卡系统，包含三个主要部分：
1. **Teams App** (位于 `src/frontend/teams-app`) - 用户在 Teams 中快速创建 Thank You Card
2. **后台服务** (位于 `src/backend`) - 处理表单提交，发送消息到 Teams Channel
   - Card Service - 卡片管理服务
   - Analytics Service - 数据分析服务
3. **Web 管理后台** (位于 `src/frontend/web-dashboard`) - 查看发送和接收情况

## 目录结构
```
src/
├── frontend/
│   ├── teams-app/          # Microsoft Teams App
│   └── web-dashboard/      # Web 管理后台
└── backend/                # 后台服务（Go）
    ├── cmd/                # 服务入口
    ├── internal/           # 业务逻辑
    └── pkg/                # 可复用包
```

## 任务清单

### ✅ Phase 1: 项目基础设施搭建 (已完成)
- [x] 创建项目目录结构
- [x] 创建 README.md 和 .gitignore
- [x] 初始化 Go 项目和前端项目
  - [x] Go module 初始化
  - [x] 使用 pnpm 创建 Teams App 项目（React + TypeScript + Vite）
  - [x] 使用 pnpm 创建 Web Dashboard 项目（React + TypeScript + Vite）
- [x] 设置开发环境配置
  - [x] 后端 .env.example 配置
  - [x] Web Dashboard 环境配置文件（config.local/test/prod.json）
  - [x] 数据库初始化脚本
- [x] 确定技术栈和依赖

### ✅ Phase 2: 认证配置（复用现有 Hydra）(已完成)
- [x] **Web Dashboard 认证配置** ✅
  - [x] 复用 Cowman 的 Hydra Client ID（cowboy-test/dev/prod）
  - [x] 复用 Cowman 的 callback URL（无需运维配置变更）
  - [x] 创建各环境配置文件（config.local/dev/test/prod.json）
  - [x] 安装 `oidc-client-ts@3.2.1` 依赖
  - [x] 实现 AuthService（基于 oidc-client-ts）
  - [x] 实现 ProtectedRoute 路由保护
  - [x] 实现 OAuth Callback 处理
  - [x] 实现 Axios HTTP 客户端（自动 token 注入和刷新）
  - [x] 实现 Logout 功能
  - [x] 构建验证通过
- [ ] Teams App 认证配置
  - [ ] Azure AD 应用注册（待决策是否需要）
  - [ ] 配置 Teams App Manifest
  - [ ] 实现 Teams SSO 登录
- [ ] 后端认证验证
  - [ ] JWT Token 验证中间件
  - [ ] RBAC 权限控制
  - [ ] PLAT UMS 集成

> **Web Dashboard 认证详情**: 
> - 完整实现文档: `/Users/youshanli/.gemini/antigravity/brain/c88d7f67-8541-49d8-bf64-eb3c7be2d268/walkthrough.md`
> - 使用 Cowman 相同的 Client ID 和 Callback URL
> - 所有环境配置完成: local (`localhost:3000`), dev/test/prod (`cowboy*.castlery.com`)
> - 无需任何 Hydra 配置变更，立即可用

### Phase 3: 后台服务开发 (Card Service)
- [x] **API 服务基础架构** ✅
  - [x] 数据库迁移脚本（5个文件）
    - [x] 001_create_cards_table.sql
    - [x] 002_create_company_values_table.sql（含预置数据）
    - [x] 003_create_card_recipients_table.sql
    - [x] 004_create_card_values_table.sql
    - [x] 005_create_employee_milestones_table.sql
  - [x] 配置数据库连接（PostgreSQL + GORM）
  - [x] 配置管理（环境变量加载和验证）
  - [x] 实现认证中间件（JWT Token 验证）
  - [x] 实现 CORS 中间件
  - [x] 实现日志中间件（Request ID 追踪）
  - [x] 实现错误处理中间件（Panic 恢复）
  - [x] 安装必要的 Go 依赖
- [x] Card Management APIs ✅
  - [x] POST /api/cards - 创建感谢卡
  - [x] GET /api/cards - 获取卡片列表（分页）
  - [x] GET /api/cards/received - 获取接收的卡片
  - [x] GET /api/cards/sent - 获取发送的卡片
  - [x] GET /api/cards/{id} - 获取卡片详情
- [x] Statistics API ✅
  - [x] GET /api/statistics/user/{id} - 用户统计
  - [x] GET /api/statistics/top10 - Top 10 被认可员工
- [ ] Teams Channel Integration
  - [ ] 实现发送消息到 Teams Channel 的功能
  - [ ] 配置 Teams Channel Webhook
  - [ ] 实现消息格式化（Adaptive Card）
  - [ ] 错误处理和重试机制
- [x] 数据模型和存储 ✅
  - [x] 实现 Card 数据模型
  - [x] 实现 Recipient 关联表
  - [x] 实现 Value 关联表
  - [x] 数据库迁移脚本

### Phase 4: Web 后台管理页面开发
- [x] Web App 基础架构 ✅
  - [x] 设置页面路由（react-router-dom）
  - [x] 配置状态管理（Zustand）
  - [x] 实现布局组件（Layout, Header, Sidebar）
  - [x] 集成 AuthService 和认证状态
- [x] 仪表盘页面 ✅
  - [x] 总览统计数据展示
  - [x] 最近发送/接收的卡片列表
  - [x] Top 10 被认可员工
- [x] 卡片管理页面 ✅
  - [x] 发送记录查看（分页、筛选）
  - [x] 接收记录查看（分页、筛选）
  - [x] 卡片详情展示
  - [x] 卡片列表组件复用
- [x] 个人统计页面 ✅
  - [x] 发送/接收数量统计
  - [x] 最常用价值观统计
  - [x] 统计图表可视化
- [x] 筛选和搜索功能 ✅
  - [x] 按价值观筛选
  - [x] 按发送人/接收人筛选
  - [x] 按时间范围筛选
  - [x] 关键词搜索
  - [x] 筛选状态管理

### Phase 5: Teams App 前端开发
- [x] Teams App 基础架构 ✅
  - [x] 创建 Teams App Manifest
  - [x] 配置 App 图标和描述
  - [x] 初始化 Teams SDK
  - [x] 配置 Teams 主题适配
  - [x] 使用固定公司价值观数据
  - [x] Hover 显示价值观详情
- [x] 感谢卡创建表单 ✅
  - [x] 收件人选择器（支持多选）
  - [x] 感谢原因文本输入
  - [x] 价值观选择器（1-3个，带 Tooltip）
  - [x] 卡片预览功能
  - [x] 表单验证
- [x] API 集成 ✅
  - [x] 集成创建 card API
  - [x] 错误处理和用户反馈

### Phase 6: Analytics Service 开发（HR 管理员功能）
- [x] Analytics APIs ✅
  - [x] GET /api/analytics/dashboard - 分析仪表盘
  - [x] POST /api/analytics/export - 数据导出（CSV）
  - [x] GET /api/analytics/recognizers/top - 最活跃认可者
  - [x] GET /api/analytics/teams - 团队分析
  - [x] GET /api/analytics/values/distribution - 价值观分布
- [x] 分析计算逻辑 ✅
  - [x] 统计聚合服务
  - [x] CSV 导出服务
  - [x] 数据缓存策略
- [x] HR 管理页面 ✅
  - [x] HR 分析仪表盘 UI
  - [x] 数据导出功能
  - [x] 团队认可模式可视化
  - [x] 权限控制（仅 HR Admin）

### Phase 7: 集成和测试
- [ ] 单元测试
  - [ ] Card Service 业务逻辑测试
  - [ ] Analytics Service 计算逻辑测试
  - [ ] 前端组件测试
- [ ] 集成测试
  - [ ] Teams App 与 Card Service 集成
  - [ ] Web App 与 Card Service 集成
  - [ ] Web App 与 Analytics Service 集成
- [ ] Teams Channel 集成测试
  - [ ] Teams Channel 消息发送测试
  - [ ] Webhook 可靠性测试
- [ ] 端到端流程测试
  - [ ] 从 Teams App 创建卡片
  - [ ] 验证后台服务接收和存储
  - [ ] 验证 Teams Channel 消息发送
  - [ ] 验证 Web 后台显示
- [ ] 性能测试
  - [ ] 并发创建测试
  - [ ] 大数据集查询测试
- [ ] 安全测试
  - [ ] 认证测试（无 Token 访问）
  - [ ] 授权测试（权限控制）
  - [ ] 数据访问控制测试

### Phase 8: 部署和文档
- [ ] 准备部署配置
  - [ ] Docker 容器化
  - [ ] Kubernetes 配置
  - [ ] 环境变量配置
- [ ] 部署到各环境
  - [ ] Dev 环境部署
  - [ ] Test 环境部署
  - [ ] Production 环境部署
- [ ] 文档编写
  - [ ] 用户使用文档
  - [ ] 运维文档
  - [ ] API 文档
  - [ ] 部署文档
- [ ] 用户验收测试
  - [ ] UAT 测试计划
  - [ ] 问题修复和优化

---

## 当前进度总结

### ✅ 已完成
- **Phase 1**: 项目基础设施搭建 - 100%
- **Phase 2**: Web Dashboard 认证配置 - 100%

### 🎯 下一步
- **Phase 3**: Card Service 后台开发 - Card Management & Statistics APIs (100%)
- **Phase 4**: Web Dashboard 业务页面开发 (90%)
- **Phase 5**: Teams App 前端开发

### 📊 整体进度
- 已完成: 2/8 阶段 (25%)
- 进行中: 0/8 阶段
- 待开始: 6/8 阶段

---

## 技术栈确认

**前端**:
- Teams App: React 18 + TypeScript + Fluent UI + Vite + pnpm
- Web Dashboard: React 18 + TypeScript + Zustand + Vite + pnpm + react-router-dom + axios + oidc-client-ts

**后端**:
- 语言: Go 1.21+
- 框架: Gin
- ORM: Gorm
- 数据库: PostgreSQL 15

**认证**:
- Teams App: Azure AD SSO（待决策）
- Web Dashboard: Hydra OAuth2（复用 Cowman Client ID）

---

## 相关文档

- **项目规划**: `plan.md`
- **需求文档**: `inception/units/`
- **详细设计**: `construction/`
- **Teams App 打包与部署**: [teams_app_bundle.md](file:///Users/youshanli/Workspaces/Programs/Castlery/thank-you-card-1/construction/teams_app_bundle.md)
- **Web Dashboard 认证实现**: `.gemini/antigravity/brain/c88d7f67-8541-49d8-bf64-eb3c7be2d268/walkthrough.md`
- **原始实施计划**: `.gemini/antigravity/brain/c29d65e9-5ccc-432c-b02a-0527888b48e5/implementation_plan.md.resolved`
