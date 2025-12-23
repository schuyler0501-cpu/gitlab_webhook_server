#!/bin/bash

# 记忆更新脚本
# 用于快速更新 .cursor/memory.md 文件

MEMORY_FILE=".cursor/memory.md"
BACKUP_FILE=".cursor/memory.backup.md"

# 创建备份
if [ -f "$MEMORY_FILE" ]; then
    cp "$MEMORY_FILE" "$BACKUP_FILE"
    echo "✅ 已创建备份: $BACKUP_FILE"
fi

# 显示当前记忆内容
echo "📝 当前记忆文件内容:"
echo "---"
cat "$MEMORY_FILE" 2>/dev/null || echo "记忆文件不存在"
echo "---"
echo ""

# 提示用户输入
echo "请输入要添加/更新的内容（输入 'q' 退出）:"
read -r content

if [ "$content" = "q" ]; then
    echo "已取消"
    exit 0
fi

# 添加时间戳和内容
timestamp=$(date "+%Y-%m-%d %H:%M:%S")
echo "" >> "$MEMORY_FILE"
echo "### $timestamp" >> "$MEMORY_FILE"
echo "$content" >> "$MEMORY_FILE"
echo "" >> "$MEMORY_FILE"

echo "✅ 记忆已更新！"

