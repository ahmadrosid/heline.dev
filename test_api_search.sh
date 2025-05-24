#!/bin/bash

# Test script to verify search functionality through API endpoints
# This script uses curl to test the actual HTTP API

API_BASE="http://localhost:8000"
TEST_QUERIES=(
    "function"
    "void"
    ": void"
    "): void"
    "function("
    "=>"
    "[]"
    "{}"
    "()"
)

echo "🔍 Testing Search API with curl commands"
echo "========================================="
echo ""

# Function to test a single query
test_query() {
    local query="$1"
    echo "📋 Testing query: '$query'"
    echo "Command: curl -X GET \"$API_BASE/api/search\" -G --data-urlencode \"q=$query\""
    
    # Execute the curl command and capture response
    response=$(curl -s -X GET "$API_BASE/api/search" -G --data-urlencode "q=$query")
    
    # Check if curl was successful
    if [ $? -ne 0 ]; then
        echo "❌ Curl command failed"
        echo ""
        return 1
    fi
    
    # Parse the response using jq if available, otherwise use grep
    if command -v jq >/dev/null 2>&1; then
        total=$(echo "$response" | jq -r '.hits.total // 0')
        has_highlighting=$(echo "$response" | grep -c "<mark>" || echo "0")
    else
        # Fallback to grep if jq is not available
        total=$(echo "$response" | grep -o '"total":[0-9]*' | cut -d':' -f2 || echo "0")
        has_highlighting=$(echo "$response" | grep -c "<mark>" || echo "0")
    fi
    
    # Display results
    if [ "$total" -gt 0 ]; then
        echo "✅ Found $total results"
        if [ "$has_highlighting" -gt 0 ]; then
            echo "🎯 Has highlighting with <mark> tags"
        else
            echo "⚠️  Results found but no highlighting detected"
        fi
    else
        echo "❌ No results found"
    fi
    
    # Show first few characters of response for debugging
    echo "📄 Response preview: $(echo "$response" | head -c 100)..."
    echo ""
}

# Function to check if server is running
check_server() {
    echo "🔧 Checking if server is running on $API_BASE..."
    if curl -s "$API_BASE/api/search?q=test" >/dev/null 2>&1; then
        echo "✅ Server is running"
        echo ""
        return 0
    else
        echo "❌ Server is not running or not accessible"
        echo "   Make sure to start the server with: make dev"
        echo ""
        return 1
    fi
}

# Function to test with verbose output
test_query_verbose() {
    local query="$1"
    echo "🔍 Verbose test for query: '$query'"
    echo "Command: curl -v -X GET \"$API_BASE/api/search\" -G --data-urlencode \"q=$query\""
    echo ""
    
    curl -v -X GET "$API_BASE/api/search" -G --data-urlencode "q=$query" 2>&1 | head -50
    echo ""
    echo "----------------------------------------"
    echo ""
}

# Main execution
main() {
    # Check if server is running
    if ! check_server; then
        exit 1
    fi
    
    # Test each query
    for query in "${TEST_QUERIES[@]}"; do
        test_query "$query"
    done
    
    echo "🎯 Summary"
    echo "=========="
    echo "All queries have been tested. Key points to verify:"
    echo "1. ': void' should return results (this was the original issue)"
    echo "2. Special character queries should have <mark> highlighting"
    echo "3. No curl errors should occur"
    echo ""
    echo "If you want to test manually:"
    echo "curl -X GET \"$API_BASE/api/search\" -G --data-urlencode \"q=: void\""
    echo ""
    
    # Offer verbose testing
    read -p "🔍 Run verbose test for ': void' query? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        test_query_verbose ": void"
    fi
}

# Run the main function
main