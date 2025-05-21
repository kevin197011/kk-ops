#!/bin/bash

# List of template files to fix
TEMPLATES=(
  "ui/templates/dashboard.html"
  "ui/templates/pipeline_detail.html"
  "ui/templates/pipeline_new.html"
  "ui/templates/pipelines.html"
  "ui/templates/server_detail.html"
  "ui/templates/server_new.html"
  "ui/templates/servers.html"
  "ui/templates/user_detail.html"
  "ui/templates/user_new.html"
  "ui/templates/users.html"
  "ui/templates/login.html"
  "ui/templates/register.html"
)

for file in "${TEMPLATES[@]}"; do
  echo "Updating icons in $file..."

  # Replace old sidebar-brand structure with new tech-oriented structure
  sed -i '' 's|<div class="sidebar-brand">[[:space:]]*<i class="bi bi-terminal-plus"></i>[[:space:]]*<span>KK 运维平台</span>[[:space:]]*</div>|<div class="sidebar-brand">\n            <div class="tech-icon">\n                <i class="bi bi-hdd-network"></i>\n                <div class="pulse-ring"></div>\n            </div>\n            <span><span class="tech-highlight">KK-OPS</span> <span class="version-tag">v2.0</span> 运维平台</span>\n        </div>|g' "$file"
done

echo "All icon updates complete!"