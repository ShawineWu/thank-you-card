# Web Dashboard 重构说明

## 概述

本次重构将 Send Card 功能从独立页面迁移为 Drawer（抽屉）形式，并整合了用户相关的所有功能到一个统一的 Profile（个人资料）页面中。

## 主要变更

### 1. 新增组件

#### UI 组件
- `src/components/ui/drawer.tsx` - Drawer 抽屉组件（基于 @radix-ui/react-dialog）
- `src/components/ui/tabs.tsx` - Tabs 标签页组件（基于 @radix-ui/react-tabs）

#### 功能组件
- `src/components/Cards/SendCardDrawer.tsx` - Send Card 的 Drawer 实现
- `src/pages/Profile/index.tsx` - 新的用户个人资料页面

### 2. 页面结构调整

#### 之前的路由结构
```
/dashboard - 概览
/received - 接收的卡片
/sent - 发送的卡片
/stats - 统计数据
/values - 公司价值观
/analytics - HR 分析
/settings - 设置
```

#### 现在的路由结构
```
/dashboard - 概览（Dashboard）
/profile - 我的资料（包含 4 个标签页）
  - Overview（概览：统计数据和图表）
  - Received Cards（接收的卡片列表）
  - Sent Cards（发送的卡片列表）
  - Account（账户信息）
/values - 公司价值观
/analytics - HR 分析
/settings - 设置
```

### 3. Send Card 功能迁移

- **之前**: `/sent` 页面包含 SendCardForm 表单
- **现在**: Header 顶部导航栏中的 "Send Card" 按钮打开 Drawer 抽屉
- **优势**: 
  - 在任何页面都可以快速发送卡片
  - 不需要离开当前页面
  - 更好的用户体验
  - 节省页面空间

### 4. Profile 页面的 4 个标签页

#### Overview（概览）标签页
- 显示用户统计数据
  - 发送的卡片总数
  - 接收的卡片总数
  - 获得的成就数量
  - 最突出的价值观
- 价值观分布图表（横向柱状图）
- 影响力比率图表（饼图：发送 vs 接收）
- 里程碑和成就展示

#### Received Cards（接收的卡片）标签页
- 显示用户接收的所有感谢卡片
- 支持按关键词、价值观、日期范围筛选
- 支持分页浏览
- 可查看卡片详情

#### Sent Cards（发送的卡片）标签页
- 显示用户发送的所有感谢卡片
- 支持按关键词、价值观、日期范围筛选
- 支持分页浏览
- 可查看卡片详情

#### Account（账户）标签页
- 显示用户账户信息
- 用户头像
- 姓名、邮箱、用户 ID 等基本信息

### 5. 删除的页面

以下页面已被移除，功能整合到 Profile 页面：
- `src/pages/ReceivedCards/index.tsx` → Profile 的 Received Cards 标签页
- `src/pages/SentCards/index.tsx` → Profile 的 Sent Cards 标签页
- `src/pages/Stats/index.tsx` → Profile 的 Overview 标签页

### 6. Sidebar 侧边栏更新

导航菜单项从 6 个精简到 4 个：
- **Overview** - 仪表盘概览
- **My Profile** - 我的资料（新增，整合了 Received/Sent/Stats）
- **Company Values** - 公司价值观
- **HR Analytics** - HR 分析

## 技术细节

### 新增依赖
```json
{
  "@radix-ui/react-tabs": "^1.1.13"
}
```

### 组件复用性增强
- `CardList` 组件新增 `showSender` 和 `showRecipient` props
  - 在 Received Cards 中隐藏接收者（因为就是当前用户）
  - 在 Sent Cards 中隐藏发送者（因为就是当前用户）
- `CardFilters` 组件新增 `showSenderFilter` 和 `showRecipientFilter` props
  - 根据场景灵活显示/隐藏筛选器

### 状态管理
- Profile 页面使用独立的状态管理 received 和 sent cards
- 使用 debounce（防抖）优化筛选器性能，避免频繁请求
- 使用 useCallback 优化数据获取函数，避免不必要的重新渲染

### Drawer 实现
- 基于 @radix-ui/react-dialog 实现
- 从右侧滑入的抽屉效果
- 支持点击遮罩层关闭
- 支持 ESC 键关闭
- 表单验证和错误提示
- 支持 AI 生成感谢理由

## 使用说明

### 如何发送感谢卡片
1. 点击页面右上角 Header 中的 **"Send Card"** 按钮
2. 在打开的 Drawer 抽屉中填写表单：
   - **Recipients（接收者）**: 搜索并选择同事
   - **Company Values（公司价值观）**: 选择 1-3 个相关价值观
   - **Key information（关键信息）**: 可选，输入关键点用于 AI 生成
   - **Recognition Reason（感谢理由）**: 手动输入或使用 AI 生成
3. 点击 **"Generate"** 按钮可使用 AI 生成感谢理由
4. 点击 **"Send Recognition"** 发送卡片

### 如何查看个人数据
1. 点击左侧 Sidebar 中的 **"My Profile"**
2. 在不同标签页中查看：
   - **Overview**: 查看统计数据、图表和成就
   - **Received Cards**: 查看收到的感谢卡片
   - **Sent Cards**: 查看发送的感谢卡片
   - **Account**: 查看账户基本信息

### 如何筛选卡片
在 Received Cards 或 Sent Cards 标签页中：
1. 使用搜索框按关键词搜索
2. 使用下拉菜单按价值观筛选
3. 使用日期选择器按时间范围筛选
4. 点击 "Reset" 按钮清除所有筛选条件

## 迁移指南

### 路由变更
如果你的代码中有指向旧路由的链接，请更新为：
- `/received` → `/profile` (默认显示 Overview，需手动切换到 Received Cards 标签页)
- `/sent` → `/profile` (需手动切换到 Sent Cards 标签页)
- `/stats` → `/profile` (Overview 标签页包含统计信息)

### 代码示例
```typescript
// 之前
<Link to="/received">查看接收的卡片</Link>

// 现在
<Link to="/profile">查看我的资料</Link>
```

## 设计考虑

### 为什么使用 Drawer 而不是 Dialog？
1. **更好的空间利用**: Drawer 从侧边滑入，不会完全遮挡页面内容
2. **更自然的交互**: 用户可以看到背景内容，保持上下文
3. **更适合表单**: 长表单在 Drawer 中滚动更自然
4. **移动端友好**: Drawer 在移动设备上的体验更好

### 为什么整合到 Profile 页面？
1. **逻辑相关性**: Received、Sent、Stats 都是用户个人数据
2. **减少导航层级**: 从 3 个独立页面变为 1 个页面的 3 个标签
3. **更好的信息架构**: 用户可以在一个地方查看所有个人相关信息
4. **提升效率**: 减少页面跳转，提高浏览效率

## 未来改进建议

### 短期改进
1. **URL 参数支持**: 允许直接链接到特定标签页
   ```
   /profile?tab=received
   /profile?tab=sent
   ```
2. **卡片详情深度链接**: 支持分享特定卡片
   ```
   /profile?tab=received&card=123
   ```
3. **键盘快捷键**: 
   - `Cmd/Ctrl + K` 打开 Send Card Drawer
   - `Cmd/Ctrl + P` 跳转到 Profile 页面

### 中期改进
1. **移动端优化**: 
   - Drawer 在小屏幕上全屏显示
   - 标签页改为下拉选择
2. **性能优化**:
   - 虚拟滚动支持大量卡片
   - 图表懒加载
3. **数据导出**: 支持导出个人统计数据为 PDF/Excel

### 长期改进
1. **个性化仪表盘**: 用户可自定义 Overview 显示的内容
2. **社交功能**: 
   - 评论和点赞卡片
   - 分享到团队频道
3. **游戏化**: 
   - 更多成就和徽章
   - 排行榜
   - 每日/每周挑战

## 测试建议

### 功能测试
- [ ] Send Card Drawer 可以正常打开和关闭
- [ ] 表单验证正常工作
- [ ] AI 生成功能正常（需要配置 DEEPSEEK_API_KEY）
- [ ] 卡片发送成功后 Drawer 自动关闭
- [ ] Profile 页面 4 个标签页都能正常切换
- [ ] 筛选和搜索功能正常
- [ ] 分页功能正常
- [ ] 图表正确显示数据

### 兼容性测试
- [ ] Chrome/Edge 浏览器
- [ ] Firefox 浏览器
- [ ] Safari 浏览器
- [ ] 移动端浏览器（iOS Safari, Chrome Mobile）

### 性能测试
- [ ] 大量卡片数据下的加载性能
- [ ] 图表渲染性能
- [ ] 筛选器响应速度

## 问题排查

### Drawer 无法打开
1. 检查 `@radix-ui/react-dialog` 是否正确安装
2. 检查 CSS 动画是否正确加载
3. 检查浏览器控制台是否有错误

### 标签页无法切换
1. 检查 `@radix-ui/react-tabs` 是否正确安装
2. 检查 Tabs 组件的 `defaultValue` 是否正确

### 图表不显示
1. 检查 `recharts` 是否正确安装
2. 检查数据格式是否正确
3. 检查容器高度是否设置

### AI 生成失败
1. 检查后端 `DEEPSEEK_API_KEY` 环境变量是否配置
2. 检查网络连接
3. 查看后端日志获取详细错误信息

## 总结

本次重构显著改善了用户体验：
- ✅ 简化了导航结构（从 7 个菜单项减少到 5 个）
- ✅ 提升了 Send Card 的可访问性（任何页面都可快速发送）
- ✅ 整合了用户相关功能（一站式查看个人数据）
- ✅ 保持了所有原有功能（无功能损失）
- ✅ 提升了代码复用性（组件更灵活）

这是一次成功的重构，为未来的功能扩展打下了良好的基础。
