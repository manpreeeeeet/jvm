### JVM implementation

JVM that can run the following code:
```java
public class Main {

    int multiplier = 3;
    public static int add(int a, int b) {
        return a + b;
    }

    public int multiply(int b) {
        return b * multiplier;
    }

    public int premultiply(int c) {
        return multiply(c) * multiply(c);
    }

    public static void main(String[] args) {
        Main main = new Main();
        int d = main.premultiply(2) * add(1,2);
    }

}
```

Run Instructions: `go build jvm2 && go run jvm2`

This is not a feature complete JVM.
JVM is hardcoded to run `main` method in `Main.class`. You can easily swap the class by changing it in `main.go`.

**Missing features that I think would be interesting to add:**
* Exception Handling.
* OOP support (interfaces,inheritance etc.).
* Array support.