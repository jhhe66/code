# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 2.8

#=============================================================================
# Special targets provided by cmake.

# Disable implicit rules so canonical targets will work.
.SUFFIXES:

# Remove some rules from gmake that .SUFFIXES does not remove.
SUFFIXES =

.SUFFIXES: .hpux_make_needs_suffix_list

# Suppress display of executed commands.
$(VERBOSE).SILENT:

# A target that is always out of date.
cmake_force:
.PHONY : cmake_force

#=============================================================================
# Set environment variables for the build.

# The shell in which to execute make rules.
SHELL = /bin/sh

# The CMake executable.
CMAKE_COMMAND = /usr/local/bin/cmake

# The command to remove a file.
RM = /usr/local/bin/cmake -E remove -f

# Escaping for special characters.
EQUALS = =

# The program to use to edit the cache.
CMAKE_EDIT_COMMAND = /usr/local/bin/ccmake

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = /home/AustinChen/code/cmake/t2

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = /home/AustinChen/code/cmake/t2

# Include any dependencies generated for this target.
include lib/CMakeFiles/b.dir/depend.make

# Include the progress variables for this target.
include lib/CMakeFiles/b.dir/progress.make

# Include the compile flags for this target's objects.
include lib/CMakeFiles/b.dir/flags.make

lib/CMakeFiles/b.dir/libb.o: lib/CMakeFiles/b.dir/flags.make
lib/CMakeFiles/b.dir/libb.o: libb/libb.c
	$(CMAKE_COMMAND) -E cmake_progress_report /home/AustinChen/code/cmake/t2/CMakeFiles $(CMAKE_PROGRESS_1)
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Building C object lib/CMakeFiles/b.dir/libb.o"
	cd /home/AustinChen/code/cmake/t2/lib && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -o CMakeFiles/b.dir/libb.o   -c /home/AustinChen/code/cmake/t2/libb/libb.c

lib/CMakeFiles/b.dir/libb.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/b.dir/libb.i"
	cd /home/AustinChen/code/cmake/t2/lib && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -E /home/AustinChen/code/cmake/t2/libb/libb.c > CMakeFiles/b.dir/libb.i

lib/CMakeFiles/b.dir/libb.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/b.dir/libb.s"
	cd /home/AustinChen/code/cmake/t2/lib && /usr/bin/cc  $(C_DEFINES) $(C_FLAGS) -S /home/AustinChen/code/cmake/t2/libb/libb.c -o CMakeFiles/b.dir/libb.s

lib/CMakeFiles/b.dir/libb.o.requires:
.PHONY : lib/CMakeFiles/b.dir/libb.o.requires

lib/CMakeFiles/b.dir/libb.o.provides: lib/CMakeFiles/b.dir/libb.o.requires
	$(MAKE) -f lib/CMakeFiles/b.dir/build.make lib/CMakeFiles/b.dir/libb.o.provides.build
.PHONY : lib/CMakeFiles/b.dir/libb.o.provides

lib/CMakeFiles/b.dir/libb.o.provides.build: lib/CMakeFiles/b.dir/libb.o

# Object files for target b
b_OBJECTS = \
"CMakeFiles/b.dir/libb.o"

# External object files for target b
b_EXTERNAL_OBJECTS =

lib/libb.so: lib/CMakeFiles/b.dir/libb.o
lib/libb.so: lib/CMakeFiles/b.dir/build.make
lib/libb.so: lib/CMakeFiles/b.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --red --bold "Linking C shared library libb.so"
	cd /home/AustinChen/code/cmake/t2/lib && $(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/b.dir/link.txt --verbose=$(VERBOSE)

# Rule to build all files generated by this target.
lib/CMakeFiles/b.dir/build: lib/libb.so
.PHONY : lib/CMakeFiles/b.dir/build

lib/CMakeFiles/b.dir/requires: lib/CMakeFiles/b.dir/libb.o.requires
.PHONY : lib/CMakeFiles/b.dir/requires

lib/CMakeFiles/b.dir/clean:
	cd /home/AustinChen/code/cmake/t2/lib && $(CMAKE_COMMAND) -P CMakeFiles/b.dir/cmake_clean.cmake
.PHONY : lib/CMakeFiles/b.dir/clean

lib/CMakeFiles/b.dir/depend:
	cd /home/AustinChen/code/cmake/t2 && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /home/AustinChen/code/cmake/t2 /home/AustinChen/code/cmake/t2/libb /home/AustinChen/code/cmake/t2 /home/AustinChen/code/cmake/t2/lib /home/AustinChen/code/cmake/t2/lib/CMakeFiles/b.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : lib/CMakeFiles/b.dir/depend

