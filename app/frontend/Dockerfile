FROM node:24-alpine3.21

WORKDIR /app

# Copy package files first for better caching
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Use JSON format for CMD and run production server
CMD ["npm", "run", "start"]