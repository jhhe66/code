# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 3.8

# Delete rule output on recipe failure.
.DELETE_ON_ERROR:


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

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = /root/code/cmake/t2

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = /root/code/cmake/t2

# Include any dependencies generated for this target.
include CMakeFiles/a_static.dir/depend.make

# Include the progress variables for this target.
include CMakeFiles/a_static.dir/progress.make

# Include the compile flags for this target's objects.
include CMakeFiles/a_static.dir/flags.make

CMakeFiles/a_static.dir/liba/liba.c.o: CMakeFiles/a_static.dir/flags.make
CMakeFiles/a_static.dir/liba/liba.c.o: liba/liba.c
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --progress-dir=/root/code/cmake/t2/CMakeFiles --progress-num=$(CMAKE_PROGRESS_1) "Building C object CMakeFiles/a_static.dir/liba/liba.c.o"
	/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -o CMakeFiles/a_static.dir/liba/liba.c.o   -c /root/code/cmake/t2/liba/liba.c

CMakeFiles/a_static.dir/liba/liba.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/a_static.dir/liba/liba.c.i"
	/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -E /root/code/cmake/t2/liba/liba.c > CMakeFiles/a_static.dir/liba/liba.c.i

CMakeFiles/a_static.dir/liba/liba.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/a_static.dir/liba/liba.c.s"
	/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -S /root/code/cmake/t2/liba/liba.c -o CMakeFiles/a_static.dir/liba/liba.c.s

CMakeFiles/a_static.dir/liba/liba.c.o.requires:

.PHONY : CMakeFiles/a_static.dir/liba/liba.c.o.requires

CMakeFiles/a_static.dir/liba/liba.c.o.provides: CMakeFiles/a_static.dir/liba/liba.c.o.requires
	$(MAKE) -f CMakeFiles/a_static.dir/build.make CMakeFiles/a_static.dir/liba/liba.c.o.provides.build
.PHONY : CMakeFiles/a_static.dir/liba/liba.c.o.provides

CMakeFiles/a_static.dir/liba/liba.c.o.provides.build: CMakeFiles/a_static.dir/liba/liba.c.o


# Object files for target a_static
a_static_OBJECTS = \
"CMakeFiles/a_static.dir/liba/liba.c.o"

# External object files for target a_static
a_static_EXTERNAL_OBJECTS =

lib/liba.a: CMakeFiles/a_static.dir/liba/liba.c.o
lib/liba.a: CMakeFiles/a_static.dir/build.make
lib/liba.a: CMakeFiles/a_static.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --bold --progress-dir=/root/code/cmake/t2/CMakeFiles --progress-num=$(CMAKE_PROGRESS_2) "Linking C static library lib/liba.a"
	$(CMAKE_COMMAND) -P CMakeFiles/a_static.dir/cmake_clean_target.cmake
	$(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/a_static.dir/link.txt --verbose=$(VERBOSE)

# Rule to build all files generated by this target.
CMakeFiles/a_static.dir/build: lib/liba.a

.PHONY : CMakeFiles/a_static.dir/build

CMakeFiles/a_static.dir/requires: CMakeFiles/a_static.dir/liba/liba.c.o.requires

.PHONY : CMakeFiles/a_static.dir/requires

CMakeFiles/a_static.dir/clean:
	$(CMAKE_COMMAND) -P CMakeFiles/a_static.dir/cmake_clean.cmake
.PHONY : CMakeFiles/a_static.dir/clean

CMakeFiles/a_static.dir/depend:
	cd /root/code/cmake/t2 && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /root/code/cmake/t2 /root/code/cmake/t2 /root/code/cmake/t2 /root/code/cmake/t2 /root/code/cmake/t2/CMakeFiles/a_static.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : CMakeFiles/a_static.dir/depend

