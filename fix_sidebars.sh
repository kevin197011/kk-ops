#!/bin/bash

# List of template files with sidebar
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
)

# Loop through each template file
for file in "${TEMPLATES[@]}"; do
  echo "Fixing $file..."

  # Find all sidebar declarations and keep only the first one
  sed -i '' '/width: 260px;/{n;/Note: sidebar-brand styles are now in tech-sidebar.css/{n;/^$/{n;/^$/{n;/background: white;/,/overflow-x: hidden;/{d;}}}}}' "$file"

  # Remove any duplicate sidebar-brand comments
  sed -i '' '/Note: sidebar-brand styles are now in tech-sidebar.css/{x;/Note: sidebar-brand styles/{d;};x;}' "$file"

  # Fix any sidebar icon (some pages use terminal-plus instead of hdd-network)
  sed -i '' 's|<div class="sidebar-brand">\n            <i class="bi bi-terminal-plus"></i>\n            <span>KK 运维平台</span>\n        </div>|<div class="sidebar-brand">\n            <div class="tech-icon">\n                <i class="bi bi-hdd-network"></i>\n                <div class="pulse-ring"></div>\n            </div>\n            <span><span class="tech-highlight">KK-OPS</span> <span class="version-tag">v2.0</span> 运维平台</span>\n        </div>|' "$file"
done

echo "All sidebar issues fixed!"