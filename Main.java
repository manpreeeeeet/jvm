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
