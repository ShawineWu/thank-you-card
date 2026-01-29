# Web Dashboard 重构说明

## 概述

本次重构将 Send Card 功能从独立页面迁移为 Drawer 形式，并整合了用户相关的所有功能到一个统一的 Profile 页面中。

## 主要变更

### 1. 新增组件

#### UI 组件
- `src/components/ui/drawer.tsx` - Drawer 组件（基于 @radix-ui/react-dialog）
- `src/components/ui/tabs.tsx` - Tabs 组件（基于 @radix-ui/react-tabs）

#### 功能组件
- `src/components/Cards/SendCardDrawer.tsx` - Send Card 的 Drawer 实现
- `src/pages/Profile/index.tsx` - 新的用户 Profile 页面

### 2. 页面结构调整

#### 之前的路由结构
```
/dashboard - Overview
/received - Received Cards
/sent - Sent Cards
/stats - Statistics
/values - Company Values
/analytics - HR Analytics
/settings - Settings
```

#### 现在的路由结构
```
/dashboard - Overview (Dashboard)
/profile - My Profile (包含 4 个 tabs)
  - Overview (统计数据和图表)
  - Received Cards (接收的卡片列表)
  - Sent Cards (发送的卡片列表)
  - Account (账户信息)
/values - Company Values
/analytics - HR Analytics
/settings - Settings
```

### 3. Send Card 功能迁移

- **之前**: `/sent` 页面包含 SendCardForm
- **现在**: Header 中的 "Send Card" 按钮打开 Drawer
- **优势**: 
  - 在任何页面都可以快速发送卡片
  - 不需要离开当前页面
  - 更好的用户体验

### 4. Profile 页面的 4 个 Tabs

#### Overview Tab
- 显示用户统计数据（发送/接收的卡片数、成就、主要价值观）
- 价值观分布图表（横向柱状图）
- 影响力比率图表（饼图）
- 里程碑和成就展示

#### Received Cards Tab
- 显示用户接收的所有卡片
- 支持筛选和搜索
- 分页显示

#### Sent Cards Tab
- 显示用户发送的所有卡片
- 支持筛选和搜索
- 分页显示

#### Account Tab
- 显示用户账户信息
- 用户头像、姓名、邮箱等

### 5. 删除的页面

以下页面已被移除，功能整合到 Profile 页面：
- `src/pages/ReceivedCards/index.tsx`
- `src/pages/SentCards/index.tsx`
- `src/pages/Stats/index.tsx`

### 6. Sidebar 更新

导航菜单项从 6 个减少到 4 个：
- Overview (Dashboard)
- My Profile (新增，整合了 Received/Sent/Stats)
- Company Values
- HR Analytics

## 技术细节

### 新增依赖
```json
{
  "@radix-ui/react-tabs": "^1.1.13"
}
```

### 组件复用
- `CardList` 组件支持 `showSender` 和 `showRecipient` props
- `CardFilters` 组件支持 `showSenderFilter` 和 `showRecipientFilter` props
- 这些 props 使得组件可以在不同场景下灵活使用

### 状态管理
- Profile 页面使用独立的状态管理 received 和 sent cards
- 使用 debounce 优化筛选器性能
- 使用 useCallback 优化数据获取函数

## 使用说明

### 发送卡片
1. 点击 Header 右上角的 "Send Card" 按钮
2. 在 Drawer 中填写表单
3. 选择接收者、公司价值观、填写感谢理由
4. 可选：使用 AI 生成感谢理由
5. 点击 "Send Recognition" 发送

### 查看个人数据
1. 点击 Sidebar 中的 "My Profile"
2. 在 Overview tab 查看统计数据和图表
3. 在 Received Cards tab 查看接收的卡片
4. 在 Sent Cards tab 查看发送的卡片
5. 在 Account tab 查看账户信息

## 迁移指南

如果你有指向旧路由的链接，请更新为：
- `/received` → `/profile` (默认显示 Overview，可切换到 Received Cards tab)
- `/sent` → `/profile` (切换到 Sent Cards tab)
- `/stats` → `/profile` (Overview tab 包含统计信息)

## 未来改进

1. 添加 URL 参数支持，允许直接链接到特定 tab
2. 添加卡片详情页面的深度链接
3. 优化移动端 Drawer 体验
4. 添加键盘快捷键支持（如 Cmd+K 打开 Send Card）
