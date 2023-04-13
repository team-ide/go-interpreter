package golang

import (
	"github.com/team-ide/go-interpreter/node"
	"testing"
)

const golangCode = `
// Importing required Java packages
import java.util.*;
import java.io.*;

// Defining a public class named "HelloWorld"
public class HelloWorld {

  // This is a single-line comment

  /*
   * This is a
   * multi-line comment
   */

  // Here's a variable declaration
  int myNumber = 42;

  // Here's a method declaration
  public void sayHello() {
    System.out.println("Hello, world!");
  }

  // Here's a main method declaration
  public static void main(String[] args) {

    // Here's a conditional statement
    if (args.length > 0) {
      System.out.println("You passed in " + args.length + " command-line arguments.");
    } else {
      System.out.println("You didn't pass in any command-line arguments.");
    }

    // Here's a loop
    for (int i = 0; i < 10; i++) {
      System.out.println("The value of i is: " + i);
    }

    // Here's a switch statement
    int myNumber = 42; // Defining a local variable named myNumber
    switch (myNumber) {
      case 42:
        System.out.println("The answer to the ultimate question of life, the universe, and everything.");
        break;
      default:
        System.out.println("I don't know what that number means.");
        break;
    }

    // Here's a try-catch block
    try {
      int result = 10 / 0;
    } catch (ArithmeticException e) {
      System.out.println("You can't divide by zero!");
    }

    // Here's an array declaration
    int[] myArray = new int[10];

    // Here's a foreach loop
    for (int value : myArray) {
      System.out.println("The value is: " + value);
    }

    // Here's a class instantiation
    HelloWorld helloWorld = new HelloWorld();

    // Here's a method invocation
    helloWorld.sayHello();
  }
}

// Importing required Java packages
import java.util.*;
import java.io.*;

// Defining a public class named "HelloWorld"
public class HelloWorld {

  // Here are all the Java keywords

  // Access modifiers
  public static void main(String[] args) {}
  private int myPrivateInt;
  protected boolean myProtectedBoolean;

  // Class declaration
  class MyClass {}

  // Object creation
  Object myObject = new Object();

  // Method declaration
  public void myMethod() {}

  // Variable declaration
  int myInt;
  boolean myBoolean;
  byte myByte;
  short myShort;
  long myLong;
  float myFloat;
  double myDouble;
  char myChar;

  // Conditional statements
  if (myInt == 0) {}
  else if (myInt == 1) {}
  else {}

  switch (myInt) {
    case 0:
      break;
    case 1:
      break;
    default:
      break;
  }

  // Loop statements
  while (myBoolean) {}
  do {} while (myBoolean);
  for (int i = 0; i < 10; i++) {}

  // Jump statements
  break;
  continue;
  return;

  // Exception handling
  try {}
  catch (Exception e) {}
  finally {}

  // Miscellaneous
  static {}
  final int MY_CONSTANT = 42;
  abstract class MyAbstractClass {}
  interface MyInterface {}
  enum MyEnum {}
}

public class AllJavaKeywords {

  // Access modifiers
  public static final int PUBLIC_VARIABLE = 1;
  private static final int PRIVATE_VARIABLE = 2;
  protected static final int PROTECTED_VARIABLE = 3;

  public static void main(String[] args) {

    // Primitive types
    boolean myBoolean = true;
    byte myByte = 1;
    short myShort = 2;
    int myInt = 3;
    long myLong = 4L;
    float myFloat = 5.0f;
    double myDouble = 6.0;

    // Control flow statements
    if (myBoolean) {
      System.out.println("if statement");
    }

    switch (myByte) {
      case 1:
        System.out.println("switch statement");
        break;
    }

    for (int i = 0; i < 10; i++) {
      System.out.println("for loop");
    }

    while (myBoolean) {
      System.out.println("while loop");
      break;
    }

    do {
      System.out.println("do-while loop");
    } while (myBoolean);

    // Exception handling
    try {
      throw new Exception("Exception message");
    } catch (Exception e) {
      System.out.println("catch block");
    } finally {
      System.out.println("finally block");
    }

    // Object-oriented programming keywords
    class MyClass {}
    interface MyInterface {}
    enum MyEnum { VALUE1, VALUE2 }
    extends MyInterface implements MyInterface {}
    abstract class MyAbstractClass {}
    final class MyFinalClass {}
    synchronized void mySynchronizedMethod() {}
    static void myStaticMethod() {}
    strictfp void myStrictfpMethod() {}
    transient int myTransientVariable;
    volatile int myVolatileVariable;

    // Miscellaneous keywords
    assert myBoolean;
    break;
    case 1:
      break;
    catch (Exception e) {}
    const int MY_CONSTANT = 7;
    continue;
    default:
      break;
    do {
      System.out.println("do-while loop");
    } while (myBoolean);
    else if (myBoolean) {}
    extends MyInterface {}
    false;
    final int myFinalVariable = 8;
    finally {}
    for (int i = 0; i < 10; i++) {}
    goto myLabel;
    if (myBoolean) {}
    implements MyInterface {}
    import java.util.*;
    instanceof MyClass;
    interface MyInterface {}
    native void myNativeMethod();
    new MyClass();
    null;
    package mypackage;
    private int myPrivateVariable;
    protected int myProtectedVariable;
    public int myPublicVariable;
    return;
    strictfp class MyStrictfpClass {}
    super.myMethod();
    switch (myByte) {}
    synchronized void mySynchronizedMethod() {}
    this.myMethod();
    throw new Exception();
    throws Exception;
    transient int myTransientVariable;
    true;
    try {}
    void myMethod() {}
    volatile int myVolatileVariable;
    while (myBoolean) {}

    // Data types
    boolean myBooleanObject = true;
    Byte myByteObject = 1;
    Short myShortObject = 2;
    Integer myIntegerObject = 3;
    Long myLongObject = 4L;
    Float myFloatObject = 5.0f;
    Double myDoubleObject = 6.0;
    Character myCharacterObject = 'a';
    String myStringObject = "Hello, world!";
    Object myObject = new Object();
    int[] myIntArray = new int[]{};
  }
}
`

func TestGolang(t *testing.T) {
	tree, err := Parse(golangCode)
	if tree != nil {
		node.OutTree(golangCode, tree)
	}
	if err != nil {
		panic(err)
	}

}
