# 令牌查询系统

这是一个独立的网页版令牌查询系统，用户可以通过输入令牌来查询其余额信息。

## 功能特性

- 现代化UI设计，响应式布局
- 实时查询令牌余额
- 显示令牌详细信息（名称、剩余额度、已用额度等）
- 额度按配比转换为美元显示（500000配额 = $1）
- 显示配比关系信息
- 错误处理和加载状态提示

## 使用方法

1. 将[tokenSearchWeb]文件夹部署到Web服务器
2. 确保后端API服务正在运行
3. 在浏览器中打开[index.html]文件
4. 输入令牌(sk-开头的字符串)进行查询

## API接口

后端提供以下API接口：

### 查询令牌信息
- **URL**: `/api/token/search`
- **方法**: POST
- **参数**: 
  ```json
  {
    "token": "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "message": "",
    "data": {
      "token_name": "测试令牌",
      "remain_quota": 1000000,
      "used_quota": 50000,
      "unlimited_quota": false,
      "expired_time": 1769990400,
      "status": 1,
      "model_ratio": 0.25
    }
  }
  ```

## 技术栈

- 前端: HTML5, CSS3, JavaScript, Tailwind CSS
- 后端: Go, Gin框架
- 数据库: SQLite/MySQL/PostgreSQL

## 部署说明

1. 确保后端服务已启动并运行
2. 将[tokenSearchWeb]文件夹放置在Web服务器目录下
3. 配置Web服务器以支持静态文件访问
4. 确保API接口可访问

## 注意事项

- 此系统独立于主前端项目，不依赖/web目录下的任何文件
- 需要确保后端API接口正常工作
- 配比关系默认使用gpt-3.5-turbo模型的配比
- 后端会自动去除令牌前缀"sk-"后再查询数据库
- 前端会将额度按500000:1的比例转换为美元显示