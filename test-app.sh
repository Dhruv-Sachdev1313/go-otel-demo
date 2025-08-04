#!/bin/bash

echo "Testing Go OpenTelemetry Demo Application"
echo "========================================="

BASE_URL="http://localhost:8080"

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/health"
echo -e "\n"

# Generate some traffic with errors
echo "2. Generating traffic with potential errors..."
for i in {1..5}; do
    echo "Request $i:"
    curl -s "$BASE_URL/error"
    echo
    sleep 1
done

echo -e "\n3. Testing cart operations..."

# Add items to different user carts
echo "Adding items to user1's cart..."
curl -s "$BASE_URL/cart/add?user_id=user1&item=laptop"
echo
curl -s "$BASE_URL/cart/add?user_id=user1&item=mouse"
echo

echo "Adding items to user2's cart..."
curl -s "$BASE_URL/cart/add?user_id=user2&item=keyboard"
echo
curl -s "$BASE_URL/cart/add?user_id=user2&item=monitor"
echo
curl -s "$BASE_URL/cart/add?user_id=user2&item=headphones"
echo

# Check cart contents
echo -e "\n4. Checking cart contents..."
echo "User1's cart:"
curl -s "$BASE_URL/cart/get?user_id=user1"
echo
echo "User2's cart:"
curl -s "$BASE_URL/cart/get?user_id=user2"
echo

# Remove an item
echo -e "\n5. Removing item from user1's cart..."
curl -s "$BASE_URL/cart/remove?user_id=user1&index=0"
echo

echo "User1's cart after removal:"
curl -s "$BASE_URL/cart/get?user_id=user1"
echo

echo -e "\n6. Generating continuous load (20 requests)..."
for i in {1..20}; do
    # Mix of different requests
    case $((i % 4)) in
        0) curl -s "$BASE_URL/health" > /dev/null ;;
        1) curl -s "$BASE_URL/error" > /dev/null ;;
        2) curl -s "$BASE_URL/cart/add?user_id=user$((i%3+1))&item=item$i" > /dev/null ;;
        3) curl -s "$BASE_URL/cart/get?user_id=user$((i%3+1))" > /dev/null ;;
    esac
    echo -n "."
done
echo -e "\n"

echo "Test completed! Check your SigNoz dashboard for metrics."
echo "Metrics being sent:"
echo "- http_errors_total (Counter)"
echo "- http_request_duration_seconds (Histogram)"  
echo "- cart_items_count (Gauge)"
