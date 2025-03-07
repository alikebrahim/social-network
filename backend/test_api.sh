#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:3000"

echo -e "${YELLOW}==== Social Network API Test Suite ====${NC}"
echo "Testing endpoints at $BASE_URL"
echo

# Function to check if response is successful
check_response() {
  if [[ "$1" == *"error"* ]]; then
    echo -e "${RED}✘ Failed: $2${NC}"
    echo "Response: $1"
    return 1
  else
    echo -e "${GREEN}✓ Success: $2${NC}"
    return 0
  fi
}

# Test 1: Create a user (if authentication is implemented)
echo -e "${YELLOW}1. Creating test user...${NC}"
RANDOM_SUFFIX=$RANDOM
USER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register -H "Content-Type: application/json" -d "{
  \"email\":\"testuser${RANDOM_SUFFIX}@example.com\",
  \"password\":\"password123\",
  \"first_name\":\"Test\",
  \"last_name\":\"User\",
  \"date_of_birth\":\"1990-01-01\",
  \"nickname\":\"testuser\",
  \"about_me\":\"Just a test user\",
  \"profile_type\":\"public\"
}")
check_response "$USER_RESPONSE" "Create user"
USER_ID=$(echo $USER_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
echo "User ID: $USER_ID"
echo

# Test 2: Create a group
echo -e "${YELLOW}2. Creating group...${NC}"
GROUP_RESPONSE=$(curl -s -X POST $BASE_URL/groups -H "Content-Type: application/json" -d '{
  "title":"Test Group",
  "description":"Group for API testing"
}')
check_response "$GROUP_RESPONSE" "Create group"
GROUP_ID=$(echo $GROUP_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
echo "Group ID: $GROUP_ID"
echo

# Test 3: Get group details
echo -e "${YELLOW}3. Retrieving group details...${NC}"
GROUP_GET_RESPONSE=$(curl -s $BASE_URL/groups/$GROUP_ID)
check_response "$GROUP_GET_RESPONSE" "Get group details"
echo

# Test 4: List all groups
echo -e "${YELLOW}4. Listing all groups...${NC}"
GROUPS_LIST_RESPONSE=$(curl -s $BASE_URL/groups)
check_response "$GROUPS_LIST_RESPONSE" "List groups"
echo

# Test 5: Search for groups
echo -e "${YELLOW}5. Searching for groups...${NC}"
SEARCH_RESPONSE=$(curl -s "$BASE_URL/groups/search?q=Test")
check_response "$SEARCH_RESPONSE" "Search groups"
echo

# Test 6: Invite a user to the group
echo -e "${YELLOW}6. Inviting user to group...${NC}"
INVITE_RESPONSE=$(curl -s -X POST $BASE_URL/groups/$GROUP_ID/invite -H "Content-Type: application/json" -d "{
  \"inviter_id\":1,
  \"invitee_id\":${USER_ID:-2}
}")
check_response "$INVITE_RESPONSE" "Invite user to group"
echo

# Test 7: List group members
echo -e "${YELLOW}7. Listing group members...${NC}"
MEMBERS_RESPONSE=$(curl -s $BASE_URL/groups/$GROUP_ID/members)
check_response "$MEMBERS_RESPONSE" "List group members"
echo

# Test 8: Create an event in the group
echo -e "${YELLOW}8. Creating group event...${NC}"
EVENT_RESPONSE=$(curl -s -X POST $BASE_URL/groups/$GROUP_ID/events -H "Content-Type: application/json" -d '{
  "title":"Test Event",
  "description":"Event for testing",
  "event_time":"2025-05-01T18:00:00Z"
}')
check_response "$EVENT_RESPONSE" "Create group event"
EVENT_ID=$(echo $EVENT_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
echo "Event ID: $EVENT_ID"
echo

# Test 9: List group events
echo -e "${YELLOW}9. Listing group events...${NC}"
EVENTS_RESPONSE=$(curl -s $BASE_URL/groups/$GROUP_ID/events)
check_response "$EVENTS_RESPONSE" "List group events"
echo

# Test 10: Respond to event
echo -e "${YELLOW}10. Responding to event...${NC}"
EVENT_RESPONSE_RESPONSE=$(curl -s -X POST $BASE_URL/groups/$GROUP_ID/events/$EVENT_ID/response -H "Content-Type: application/json" -d '{
  "response":"going"
}')
check_response "$EVENT_RESPONSE_RESPONSE" "Respond to event"
echo

# Test 11: Get group chat history
echo -e "${YELLOW}11. Getting group chat history...${NC}"
CHAT_HISTORY_RESPONSE=$(curl -s $BASE_URL/groups/$GROUP_ID/chat)
check_response "$CHAT_HISTORY_RESPONSE" "Get group chat history"
echo

echo -e "${YELLOW}==== Test Summary ====${NC}"
echo "All tests completed! Check the results above for any failures."
echo "To test WebSocket functionality, use the websocket_test.html file."