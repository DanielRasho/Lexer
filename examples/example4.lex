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
        ASSIGN
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
"if"      { return IF }
"func"      { return FUNC }
"else"    { return ELSE }
"while"   { return WHILE }
"return"  { return RETURN }
"="       { return ASSIGN }
"\+"       { return PLUS }
"-"       { return MINUS }
"\*"       { return MULT }
"/"       { return DIV }
"\("    { return RPAREN}
"\)"    { return LPAREN}
{id}      { return ID }
{number}  { return NUMBER }
{ws}      {}
%%

%{
    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file
    //  ------ TOKENS ID -----
    // Define the token types that the lexer will recognize
    //This is a footer
%}