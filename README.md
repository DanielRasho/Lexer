<h1 align="center">YAAALex 🚀</h1>
<h3 align="center">(Yet Another Another Another Lexer) </h3>

## Getting Started 🎬

```bash
task build                              // Builds YAAAlex app
task run <YALex file> <Output path>     // Runs YAAAlex with an definition file and and output file
task testLex <YALex file >              // Builds and compiles a lexer file, and run it with a dummy main.
task test                               // Run tests
task clean                              // Removes executables
```

## The YALex File 📄
Since YALex initial definition was meant for C, we tweak it a little bit to be easer to work with using Go. Below is the structure for a YALEX go file. You can find more examples on `examples/`

```
// Use "//" for comments
{ 
    // ======= HEADER =======
    // The entire contents of this section will be COPIED to the BEGINING of the generated Lexer.go file
}
{
    // ====== NAMED PATTERNS =======
    // Definition 
    // - A Pattern should be defined in a single line
    // - Use "{}" to refer to Named patterns defined before
    // - use "\" to scape "{}" if you want them within a Regex expresion

    let LETTER = [a-zA-Z]
    let DIGIT = [0-9]
    let ID = {LETTER}({LETTER}|{DIGIT})*  // ID is a combination of LETTER and DIGIT
    let NUMBER = {DIGIT}+  // A NUMBER consists of one or more DIGITS
} 
{
    //  ===== TOKENS ID =======
    // Definition of the possible token types the compiled lexer can output, this wil be later
    // used for the action definition when a pattern is matched.

    //  ** The order in which they are written also defines ITS PRIORITY. 
    // If the lexer happens to recognize 2 possible token types for a lexeme it will take 
    // the one with HIGHEST PRIORITY (decleared first here).

    LITERAL
    NUMBER
    COND 
}
{
    // ======= RULES ========
    // Define how to react when certain patterns are matched.
    // All Rules are composed by a "pattern" and an "action".
    // The action, is ANY GO CODE that will be executed when the lexeme is recognized
    // - They may end with a return statement using any ID defined in the TOKENS ID section
    // - If there is not return statement, the Lexer wont yield any token when that pattern is matched.
    // - Use "{}" to refer to named patterns defined before

    {LETTER} { return LITERAL }
    "if" { return COND }
}
{
    // ======== FOOTER =======
    // The entire contents of this section will be COPIED to the END of the generated Lexer.go file
}
```