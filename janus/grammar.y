
%{
%}

%token STRING

%%

file:
	STRING
	| STRING '+' STRING;
	| STRING "==" STRING;
	| /* empty */
	;

%%

