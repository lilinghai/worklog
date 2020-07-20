use std::io;
fn main() {
    //macro not function !
    //; represent the statement ends
    println!("Guess the number");
    println!("Pls input your guess.");
    //variable default is inmutable; mut means mutable
    //String::new() is standard lib
    let mut guess = String::new();
    
    //read_line() return io::Result instance
    io::stdin()
        .read_line(&mut guess)
        .expect("failed to read line");
    println!("You guessed: {}", guess);
}
