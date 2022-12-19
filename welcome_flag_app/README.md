# Flag Parse

## Command-Line user interface

A typically command-line application will be an interface similar to the following:

``application [-h] [-n <value>] -silent <arg1> <arg2>``

The user interface has the following:
- **-h** is a Boolean option usually specified to print a help text;
- **-n <value>** expects the user to specify a value for the option, **n**. The app's logic determines the expected data type for the value;
- **-silent** is another Boolean option. Specifying it sets the value to **true**;
- **arg1** and **arg2** are referred to as **positional arguments**. A *positional argument's* data type and interpretation are completely determined by the app. 

