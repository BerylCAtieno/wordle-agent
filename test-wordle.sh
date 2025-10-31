#!/bin/bash

# Test 1: Start a new game
echo "=== Test 1: Starting new game ==="
curl -X POST http://localhost:5001/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "test-001",
    "method": "message/send",
    "params": {
      "message": {
        "kind": "message",
        "role": "user",
        "parts": [
          {
            "kind": "text",
            "text": "new game"
          }
        ],
        "messageId": "msg-001"
      },
      "configuration": {
        "blocking": true,
        "acceptedOutputModes": ["text/plain"]
      }
    }
  }' | jq '.'

echo -e "\n\n=== Test 2: Make a guess ==="
# Test 2: Make a guess (use the contextId from Test 1)
curl -X POST http://localhost:5001/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "test-002",
    "method": "message/send",
    "params": {
      "message": {
        "kind": "message",
        "role": "user",
        "parts": [
          {
            "kind": "text",
            "text": "CRANE"
          }
        ],
        "messageId": "msg-002"
      },
      "configuration": {
        "blocking": true,
        "acceptedOutputModes": ["text/plain"]
      }
    }
  }' | jq '.'

echo -e "\n\n=== Test 3: Invalid word ==="
curl -X POST http://localhost:5001/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "test-003",
    "method": "message/send",
    "params": {
      "message": {
        "kind": "message",
        "role": "user",
        "parts": [
          {
            "kind": "text",
            "text": "ZZZZZ"
          }
        ],
        "messageId": "msg-003"
      },
      "configuration": {
        "blocking": true,
        "acceptedOutputModes": ["text/plain"]
      }
    }
  }' | jq '.'