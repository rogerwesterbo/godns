#!/usr/bin/env bash
# Script to clean up old GitHub Actions caches
# Requires: gh CLI (GitHub CLI)

set -e

echo "🧹 Cleaning up old GitHub Actions caches..."

# Check if gh is installed
if ! command -v gh &> /dev/null; then
    echo "❌ Error: GitHub CLI (gh) is not installed"
    echo "Install it from: https://cli.github.com/"
    exit 1
fi

# Check if gh-actions-cache extension is installed
if ! gh extension list | grep -q "actions/gh-actions-cache"; then
    echo "📦 Installing gh-actions-cache extension..."
    gh extension install actions/gh-actions-cache
fi

echo ""
echo "📊 Current cache statistics:"
gh actions-cache list -L 10

echo ""
read -p "🗑️  Delete ALL caches? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Deleting all caches..."
    echo ""
    
    total_deleted=0
    
    # Keep deleting until no caches remain
    while true; do
        # Get all cache keys from the list (field 1, skip header)
        cache_list=$(gh actions-cache list -L 100 2>/dev/null) || true
        
        # Count total caches (excluding header)
        cache_count=$(echo "$cache_list" | tail -n +2 | wc -l | tr -d ' ')
        
        if [ "$cache_count" -eq 0 ]; then
            echo ""
            echo "✅ No more caches to delete"
            break
        fi
        
        # Get first cache key (use awk to get first field which handles multiple spaces)
        first_key=$(echo "$cache_list" | tail -n +2 | head -1 | awk '{print $1}')
        
        if [ -n "$first_key" ]; then
            echo "[$total_deleted/$cache_count] Deleting: $first_key"
            # Delete and capture result, don't let errors stop the script
            if gh actions-cache delete "$first_key" --confirm 2>&1 | grep -q "Deleted"; then
                ((total_deleted++)) || true
            else
                echo "  (delete failed, continuing...)"
                ((total_deleted++)) || true
            fi
        else
            echo "Could not extract cache key, stopping"
            break
        fi
    done
    
    echo ""
    echo "✅ Total deleted: $total_deleted caches"
    exit 0
else
    echo "❌ Cancelled"
    exit 0
fi

echo ""
echo "📊 Remaining caches:"
gh actions-cache list -L 10 2>/dev/null || echo "No caches remaining"
