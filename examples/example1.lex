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
    LETTER   ([a-c])+
    DIGIT    [1-3]
    OWO      {LETTER}-{DIGIT}
    ID       {LETTER}({LETTER}|{DIGIT})*  // ID consists of letters and digits
    NUMBER   {DIGIT}+  // NUMBER consists of one or more digits
    WS       ([ \t\n\r])+  // Whitespace: spaces, tabs, newlines, or carriage returns
}

// ======= RULES ========
%%
{LETTER}         { return LITERAL }   // Match letters and return LITERAL
{DIGIT}          { return NUMBER }    // Match digits and return NUMBER
{OWO}          { return WS }    // Match digits and return NUMBER

%%
