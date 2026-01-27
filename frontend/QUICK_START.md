# 快速启动指南

## 当前状态

- ✅ 后端运行在: http://localhost:8056
- ✅ 前端运行在: http://localhost:3000
- ✅ API 代理已配置

## 如果页面空白，请检查：

1. **打开浏览器开发者工具** (F12 或 Cmd+Option+I)
   - 查看 Console 标签页是否有错误
   - 查看 Network 标签页，确认 API 请求是否成功

2. **清除浏览器缓存**
   - 按 Cmd+Shift+R (Mac) 或 Ctrl+Shift+R (Windows) 强制刷新
   - 或者在开发者工具中右键刷新按钮，选择"清空缓存并硬性重新加载"

3. **检查后端连接**
   ```bash
   curl http://localhost:8056/api/cards/feed
   ```
   应该返回: `{"items":[],"total":0}`

4. **检查前端日志**
   查看运行 `make frontend-dev` 的终端，看是否有编译错误

## 常见问题

### 页面完全空白
- 检查浏览器控制台是否有 JavaScript 错误
- 确认 React 是否正确加载
- 尝试硬刷新页面 (Cmd+Shift+R)

### API 请求失败
- 确认后端正在运行: `lsof -ti:8056`
- 检查 API 代理配置是否正确

### 样式不显示
- 确认 Tailwind CSS 已正确加载
- 检查浏览器 Network 标签，确认 CSS 文件已加载
