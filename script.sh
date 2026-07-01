#!/bin/bash

# Start backend
echo "Starting backend..."
cd backend || exit

go run . "$@" &
BACKEND_PID=$!

# Return to project root
cd ..

# Start frontend
echo "Starting frontend..."
cd frontend || exit

node server.js &
FRONTEND_PID=$!

echo "Backend PID: $BACKEND_PID"
echo "Frontend PID: $FRONTEND_PID"

# Wait for both processes
wait