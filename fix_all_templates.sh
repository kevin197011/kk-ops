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
  echo "Fixing $file..."

  # 1. Remove duplicated sidebar styles
  sed -i '' -E '/Note: sidebar-brand styles are now in tech-sidebar.css/{n;n;/background: white;/,/overflow-x: hidden;/d;}' "$file"

  # 2. Fix sidebar HTML structure
  sed -i '' 's|<div class="sidebar-brand">.*<i class="bi bi-terminal-plus"></i>.*<span>KK 运维平台</span>.*</div>|<div class="sidebar-brand">\
            <div class="tech-icon">\
                <i class="bi bi-hdd-network"></i>\
                <div class="pulse-ring"></div>\
            </div>\
            <span><span class="tech-highlight">KK-OPS</span> <span class="version-tag">v2.0</span> 运维平台</span>\
        </div>|' "$file"

  # 3. Make sure there's only one sidebar declaration
  sed -i '' -E '/\.sidebar \{/{n;n;n;/\.sidebar \{/,/\}/d;}' "$file"

  echo "✓ Fixed $file"
done

echo "All templates fixed successfully!"