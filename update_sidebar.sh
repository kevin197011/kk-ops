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

# Add the CSS link to each file
for file in "${TEMPLATES[@]}"; do
  echo "Updating $file..."

  # Add CSS link if it doesn't exist
  if ! grep -q "tech-sidebar.css" "$file"; then
    sed -i '' 's|<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.3/font/bootstrap-icons.css">|<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.3/font/bootstrap-icons.css">\n    <link rel="stylesheet" href="/static/css/tech-sidebar.css">|' "$file"
  fi

  # Remove existing sidebar-brand CSS if present
  sed -i '' '/.sidebar-brand {/,/}/d' "$file"
  sed -i '' '/.sidebar-brand i {/,/}/d' "$file"

  # Add comment noting sidebar-brand styles are in separate CSS
  sed -i '' 's|.sidebar {|.sidebar {\n            background: white;\n            box-shadow: 0 1px 3px 0 rgba(0,0,0,0.1), 0 1px 2px 0 rgba(0,0,0,0.06);\n            position: fixed;\n            top: 0;\n            bottom: 0;\n            left: 0;\n            z-index: 100;\n            width: 260px;\n            padding-top: 0;\n            overflow-x: hidden;\n        }\n\n        /* Note: sidebar-brand styles are now in tech-sidebar.css */\n\n        |' "$file"

  # Update the sidebar-brand HTML structure
  sed -i '' 's|<div class="sidebar-brand">\n            <i class="bi bi-terminal-plus"></i>\n            <span>KK 运维平台</span>\n        </div>|<div class="sidebar-brand">\n            <div class="tech-icon">\n                <i class="bi bi-terminal-plus"></i>\n                <div class="pulse-ring"></div>\n            </div>\n            <span><span class="tech-highlight">KK-OPS</span> <span class="version-tag">v2.0</span> 运维平台</span>\n        </div>|' "$file"
done

echo "All files updated successfully!"