#!/bin/bash

# 保存有效对话脚本 (Linux/Mac)
# 当用户使用 [有效对话] 标识时，自动保存对话内容

if [ $# -lt 2 ]; then
    echo "用法: $0 <用户消息> <AI回复>"
    exit 1
fi

USER_MESSAGE="$1"
AI_MESSAGE="$2"
CONVERSATIONS_FILE=".cursor/conversations.md"
BACKUP_FILE=".cursor/conversations.backup.md"

# 创建备份
if [ -f "$CONVERSATIONS_FILE" ]; then
    cp "$CONVERSATIONS_FILE" "$BACKUP_FILE"
fi

# 获取当前北京时间
# 注意：需要系统支持 TZ 环境变量
export TZ='Asia/Shanghai'
TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")

# 构建新的对话记录
NEW_ENTRY="

## 对话记录 - $TIMESTAMP

### 用户
$USER_MESSAGE

### AI 回复
$AI_MESSAGE

---
"

# 追加到文件
echo "$NEW_ENTRY" >> "$CONVERSATIONS_FILE"

echo "✅ 有效对话已保存到: $CONVERSATIONS_FILE"
echo "📅 时间: $TIMESTAMP"

