# Makefile for MKV encryption C functions
CC = gcc
CFLAGS = -Wall -Wextra -O2 -fPIC
LDFLAGS = -shared
LIBS = -lssl -lcrypto

# Target library
TARGET = libmkv.so

# Source files
SOURCES = mkv/mkv.c mkv/MKV256.c
OBJECTS = $(SOURCES:.c=.o)

# Default target
all: $(TARGET)

# Build shared library
$(TARGET): $(OBJECTS)
	$(CC) $(LDFLAGS) -o $@ $(OBJECTS) $(LIBS)

# Compile C source files
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# Test target
test: $(TARGET)
	@echo "MKV library built successfully: $(TARGET)"
	@echo "Run Go tests with: go test ./..."

# Clean build artifacts
clean:
	rm -f $(OBJECTS) $(TARGET)
	rm -f mkv/*.o

# Install library (optional)
install: $(TARGET)
	sudo cp $(TARGET) /usr/local/lib/
	sudo ldconfig

# Uninstall library
uninstall:
	sudo rm -f /usr/local/lib/$(TARGET)
	sudo ldconfig

.PHONY: all clean install uninstall test 