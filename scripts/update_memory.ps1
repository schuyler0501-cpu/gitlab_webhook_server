# 记忆更新脚本 (PowerShell)
# 用于快速更新 .cursor/memory.md 文件

$MemoryFile = ".cursor/memory.md"
$BackupFile = ".cursor/memory.backup.md"

# 创建备份
if (Test-Path $MemoryFile) {
    Copy-Item $MemoryFile $BackupFile
    Write-Host "✅ 已创建备份: $BackupFile" -ForegroundColor Green
}

# 显示当前记忆内容
Write-Host "`n📝 当前记忆文件内容:" -ForegroundColor Cyan
Write-Host "---" -ForegroundColor Gray
if (Test-Path $MemoryFile) {
    Get-Content $MemoryFile
} else {
    Write-Host "记忆文件不存在" -ForegroundColor Yellow
}
Write-Host "---`n" -ForegroundColor Gray

# 提示用户输入
Write-Host "请输入要添加/更新的内容（输入 'q' 退出）:" -ForegroundColor Yellow
$content = Read-Host

if ($content -eq "q") {
    Write-Host "已取消" -ForegroundColor Yellow
    exit 0
}

# 添加时间戳和内容
$timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
$newContent = @"

### $timestamp
$content

"@

Add-Content -Path $MemoryFile -Value $newContent

Write-Host "✅ 记忆已更新！" -ForegroundColor Green

