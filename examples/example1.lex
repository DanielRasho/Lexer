// ======= HEADER =======
%{
    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file
    //  ------ TOKENS ID -----
    // Define the token types that the lexer will recognize
    const (
        LITERAL = iota
        NUMBER
        COND
        WS
    )
%}

// ====== NAMED PATTERNS =======
{
    // Define named patterns using regular expressions
    LETTER   [a-zA-Z]
    DIGIT    [0-9]
    ID       {LETTER}({LETTER}|{DIGIT})*  // ID consists of letters and digits
    NUMBER   {DIGIT}+  // NUMBER consists of one or more digits
    WS       [ \t\n\r]+  // Whitespace: spaces, tabs, newlines, or carriage returns
    COND     (if|else|while)  // Keywords like 'if', 'else', 'while'
}

// ======= RULES ========
%%
{LETTER}         { return LITERAL }   // Match letters and return LITERAL
{DIGIT}          { return NUMBER }    // Match digits and return NUMBER
{COND}           { return COND }      // Match conditional keywords and return COND
{WS}             { return WS }        // Match whitespace and return WS
{ID}             { return LITERAL }   // Match identifiers and return LITERAL

%%
