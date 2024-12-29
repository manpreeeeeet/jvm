package main

func main() {
	jvm := JVM{classes: make([]*Class, 0)}
	_, err := jvm.addClass("Main")
	if err != nil {
		return
	}
	_, err = jvm.executeMethod("Main", "main", "([Ljava/lang/String;)V")
	if err != nil {
		return
	}

}
