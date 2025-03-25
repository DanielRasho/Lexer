// ======= HEADER =======
%{
    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file
    //  ------ TOKENS ID -----
    // Define the token types that the lexer will recognize
    const (
        IF = iota
        ELSE
        WHILE
        RETURN 
        ASIGN
        PLUS
        MINUS
        MULT
        DIV
        LPAREN
        RPAREN
        LBRACE
        RBRACE
        ID
        NUMBER
        WS
    )
%}

// ====== NAMED PATTERNS =======
{
    // Define named patterns using regular expressions
    digit   [0-2]
    letter  [a-cA-B]
    id      {letter}({letter}|{digit})*
    number  ({digit})+  // ID consists of letters and digits
    ws      ([ \t\n\t])+  // NUMBER consists of one or more digits
}

// ======= RULES ========
%%
{id}      { return ID }
{number}  { return NUMBER }
{ws}      {}
%%
