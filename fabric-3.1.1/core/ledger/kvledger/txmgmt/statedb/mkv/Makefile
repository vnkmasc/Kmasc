CC = gcc
CFLAGS = -Wall -Wextra -O2 -fPIC -I.
LDFLAGS = -shared

TARGET = libmkv.so
SOURCES = mkv.c MKV256.c
OBJECTS = $(SOURCES:.c=.o)

all: $(TARGET)

$(TARGET): $(OBJECTS)
	$(CC) $(LDFLAGS) -o $@ $(OBJECTS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(OBJECTS) $(TARGET)

.PHONY: all clean 