# Quick Start Guide - Web Dashboard 重构版

## 安装依赖

```bash
cd src/frontend/web-dashboard
pnpm install
```

## 开发模式运行

```bash
pnpm run dev
```

访问 `https://localhost:5173` (注意是 HTTPS)

## 构建生产版本

```bash
pnpm run build
```

构建产物在 `dist/` 目录

## 主要功能测试

### 1. 测试 Send Card Drawer

1. 登录后，点击右上角的 **"Send Card"** 按钮
2. Drawer 应该从右侧滑入
3. 填写表单：
   - 搜索并选择接收者
   - 选择 1-3 个公司价值观
   - 输入感谢理由（或使用 AI 生成）
4. 点击 "Send Recognition" 发送
5. 成功后 Drawer 自动关闭，显示成功提示

### 2. 测试 Profile 页面

1. 点击左侧菜单的 **"My Profile"**
2. 测试 4 个标签页：
   - **Overview**: 查看统计卡片、图表、成就
   - **Received Cards**: 查看接收的卡片列表，测试筛选和分页
   - **Sent Cards**: 查看发送的卡片列表，测试筛选和分页
   - **Account**: 查看账户信息

### 3. 测试导航

1. 从 Dashboard 点击 "Send Card" → 应该打开 Drawer
2. 从 Profile 点击 "Send Card" → 应该打开 Drawer
3. 从任何页面点击 "My Profile" → 应该跳转到 Profile 页面

## 新增的路由

| 路由 | 页面 | 说明 |
|------|------|------|
| `/dashboard` | Dashboard | 仪表盘概览 |
| `/profile` | Profile | 用户资料（4 个标签页） |
| `/values` | Company Values | 公司价值观管理 |
| `/analytics` | HR Analytics | HR 分析 |
| `/settings` | Settings | 设置 |

## 移除的路由

以下路由已被移除，功能整合到 `/profile`：
- ~~`/received`~~ → `/profile` (Received Cards tab)
- ~~`/sent`~~ → `/profile` (Sent Cards tab)
- ~~`/stats`~~ → `/profile` (Overview tab)

## 环境变量

确保后端配置了以下环境变量（用于 AI 生成功能）：

```bash
DEEPSEEK_API_KEY=your_api_key_here
```

如果没有配置，AI 生成功能会失败，但不影响手动输入。

## 常见问题

### Q: Drawer 打开后无法关闭？
A: 点击遮罩层或按 ESC 键可以关闭。如果还是无法关闭，检查浏览器控制台是否有错误。

### Q: 图表不显示？
A: 确保有数据。如果是新用户，可能需要先发送或接收一些卡片。

### Q: AI 生成失败？
A: 检查后端是否配置了 `DEEPSEEK_API_KEY`。如果没有，可以手动输入感谢理由。

### Q: 标签页切换没有反应？
A: 检查浏览器控制台是否有错误。确保 `@radix-ui/react-tabs` 已正确安装。

## 开发提示

### 修改 Drawer 样式

编辑 `src/components/ui/drawer.tsx`：

```typescript
// 修改宽度
className="... max-w-2xl ..."  // 改为 max-w-3xl 或其他值

// 修改动画
className="... duration-300 ..."  // 改为 duration-500 等
```

### 修改 Profile 默认标签页

编辑 `src/pages/Profile/index.tsx`：

```typescript
<Tabs defaultValue="overview" ...>  // 改为 "received", "sent", "account"
```

### 添加新的统计卡片

在 `src/pages/Profile/index.tsx` 的 Overview tab 中添加：

```typescript
<Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
  <CardHeader className="flex flex-row items-center justify-between pb-2">
    <CardTitle className="text-sm font-medium text-slate-400">
      新指标
    </CardTitle>
    <YourIcon className="h-4 w-4 text-primary" />
  </CardHeader>
  <CardContent>
    <div className="text-3xl font-bold">{yourValue}</div>
    <p className="text-xs text-slate-500 mt-1">描述</p>
  </CardContent>
</Card>
```

## 技术栈

- **React 19** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Tailwind CSS** - 样式
- **Radix UI** - 无障碍组件
- **Recharts** - 图表库
- **React Router** - 路由
- **Zustand** - 状态管理
- **Axios** - HTTP 客户端

## 下一步

1. 阅读 `REFACTORING-ZH.md` 了解详细的重构说明
2. 查看 `src/pages/Profile/index.tsx` 了解 Profile 页面实现
3. 查看 `src/components/Cards/SendCardDrawer.tsx` 了解 Drawer 实现
4. 根据需求进行自定义修改

## 反馈

如有问题或建议，请联系开发团队。
