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
    WS ([ \t\n])+
}

// ======= RULES ========
%%
"a"         { return LITERAL }   // Match letters and return LITERAL
"b"          { return NUMBER }    // Match digits and return NUMBER
{WS}         { return WS }
%%
