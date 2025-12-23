# 保存有效对话脚本 (PowerShell)
# 当用户使用 [有效对话] 标识时，自动保存对话内容

param(
    [Parameter(Mandatory=$true)]
    [string]$UserMessage,
    
    [Parameter(Mandatory=$true)]
    [string]$AIMessage
)

$ConversationsFile = ".cursor/conversations.md"
$BackupFile = ".cursor/conversations.backup.md"

# 创建备份
if (Test-Path $ConversationsFile) {
    Copy-Item $ConversationsFile $BackupFile -Force
}

# 获取当前北京时间
$beijingTime = [System.TimeZoneInfo]::ConvertTimeBySystemTimeZoneId(
    [DateTime]::Now,
    "China Standard Time"
)
$timestamp = $beijingTime.ToString("yyyy-MM-dd HH:mm:ss")

# 读取现有内容
$existingContent = ""
if (Test-Path $ConversationsFile) {
    $existingContent = Get-Content $ConversationsFile -Raw
}

# 构建新的对话记录
$newEntry = @"

## 对话记录 - $timestamp

### 用户
$UserMessage

### AI 回复
$AIMessage

---

"@

# 追加到文件
Add-Content -Path $ConversationsFile -Value $newEntry -Encoding UTF8

Write-Host "✅ 有效对话已保存到: $ConversationsFile" -ForegroundColor Green
Write-Host "📅 时间: $timestamp" -ForegroundColor Cyan

