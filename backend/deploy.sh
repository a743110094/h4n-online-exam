#!/bin/bash

# 在线考试系统后端部署脚本
# 使用方法: ./deploy.sh

echo "=== 在线考试系统后端部署脚本 ==="

# 检查是否存在可执行文件
if [ ! -f "exam-backend-linux-amd64-pure" ]; then
    echo "错误: 找不到 exam-backend-linux-amd64-pure 文件"
    echo "请先运行编译命令: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o exam-backend-linux-amd64-pure main.go"
    exit 1
fi

# 检查是否存在环境配置文件
if [ ! -f ".env" ]; then
    echo "警告: 找不到 .env 文件，将使用 .env.example 作为模板"
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo "已复制 .env.example 到 .env，请根据实际情况修改配置"
    else
        echo "错误: 找不到 .env.example 文件"
        exit 1
    fi
fi

# 创建部署包
echo "正在创建部署包..."
DEPLOY_DIR="exam-backend-deploy"
rm -rf $DEPLOY_DIR
mkdir -p $DEPLOY_DIR

# 复制必要文件
cp exam-backend-linux-amd64-pure $DEPLOY_DIR/exam-backend
cp .env $DEPLOY_DIR/
cp test.db $DEPLOY_DIR/ 2>/dev/null || echo "注意: 未找到 test.db 文件，将在首次运行时自动创建"

# 创建启动脚本
cat > $DEPLOY_DIR/start.sh << 'EOF'
#!/bin/bash

# 设置可执行权限
chmod +x exam-backend

# 启动后端服务
echo "启动在线考试系统后端服务..."
echo "访问地址: http://localhost:8080"
echo "按 Ctrl+C 停止服务"

./exam-backend
EOF

# 创建系统服务文件（可选）
cat > $DEPLOY_DIR/exam-backend.service << 'EOF'
[Unit]
Description=Online Exam System Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/exam-backend
ExecStart=/opt/exam-backend/exam-backend
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 创建部署说明
cat > $DEPLOY_DIR/README.md << 'EOF'
# 在线考试系统后端部署说明

## 快速部署

1. 将整个文件夹上传到服务器
2. 运行启动脚本:
   ```bash
   chmod +x start.sh
   ./start.sh
   ```

## 系统服务部署（推荐）

1. 将文件复制到系统目录:
   ```bash
   sudo mkdir -p /opt/exam-backend
   sudo cp -r * /opt/exam-backend/
   sudo chown -R www-data:www-data /opt/exam-backend
   ```

2. 安装系统服务:
   ```bash
   sudo cp exam-backend.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable exam-backend
   sudo systemctl start exam-backend
   ```

3. 查看服务状态:
   ```bash
   sudo systemctl status exam-backend
   ```

## 配置说明

- 编辑 `.env` 文件修改配置
- 默认端口: 8080
- 数据库: SQLite (test.db)
- 默认管理员账户: admin/password
- 默认教师账户: teacher1/password  
- 默认学生账户: student1/password

## 防火墙配置

```bash
# Ubuntu/Debian
sudo ufw allow 8080

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

## 反向代理配置 (Nginx)

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
EOF

# 设置权限
chmod +x $DEPLOY_DIR/start.sh
chmod +x $DEPLOY_DIR/exam-backend

echo "部署包创建完成！"
echo "部署包位置: $DEPLOY_DIR/"
echo "文件列表:"
ls -la $DEPLOY_DIR/
echo ""
echo "=== 部署步骤 ==="
echo "1. 将 $DEPLOY_DIR 文件夹上传到服务器"
echo "2. 在服务器上运行: chmod +x start.sh && ./start.sh"
echo "3. 访问 http://服务器IP:8080 测试"
echo "4. 详细说明请查看 $DEPLOY_DIR/README.md"