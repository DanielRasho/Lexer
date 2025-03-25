// ======= HEADER =======
%{
    // Token definitions
    const (
        PRINT = iota
        VAR
        ASSIGN
        ADD
        SUB
        NUMBER
        ID
        WS
    )
%}

// ====== NAMED PATTERNS =======
{
    digit   [0-9]
    letter  [a-zA-Z]
    id      {letter}({letter}|{digit})*
    number  ({digit})+
    ws      ([ \t\n])+
}

// ======= RULES ========
%%
"print"    { return PRINT }
"var"      { return VAR }
"="        { return ASSIGN }
"\+"       { return ADD }
"-"        { return SUB }
{ws}       {}  // Ignore whitespace
{id}       { return ID }
{number}   { return NUMBER }
%%

// ======= FOOTER =======
%{
    // Footer section
%}
