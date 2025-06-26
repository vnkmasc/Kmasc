# Makefile for encryption C functions
CC = gcc
CFLAGS = -Wall -Wextra -O2 -fPIC
LDFLAGS = -shared
LIBS = -lssl -lcrypto

# Target library
TARGET = libencryption.so

# Source files
SOURCES = encrypt.c
OBJECTS = $(SOURCES:.c=.o)

# Default target
all: $(TARGET)

# Build shared library
$(TARGET): $(OBJECTS)
	$(CC) $(LDFLAGS) -o $@ $(OBJECTS) $(LIBS)

# Compile C source files
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# Clean build artifacts
clean:
	rm -f $(OBJECTS) $(TARGET)

# Install library (optional)
install: $(TARGET)
	sudo cp $(TARGET) /usr/local/lib/
	sudo ldconfig

# Uninstall library
uninstall:
	sudo rm -f /usr/local/lib/$(TARGET)
	sudo ldconfig

.PHONY: all clean install uninstall 