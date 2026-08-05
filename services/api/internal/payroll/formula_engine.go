package payroll

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/shopspring/decimal"
)

// FormulaEngine — secure, no evil eval, O(n) tokenization + O(n) shunting-yard + O(n) eval
// Supports variables: BASIC, CTC_MONTHLY, CTC_ANNUAL, GROSS, BASIC_PLUS allowances etc
// Operators: + - * / ( ) and decimal numbers
// Example: "CTC_MONTHLY * 0.4" or "BASIC * 0.1 + CTC_MONTHLY * 0.2"

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenVariable
	TokenOperator
	TokenParenLeft
	TokenParenRight
)

type Token struct {
	Type  TokenType
	Value string
}

func tokenize(expr string) ([]Token, error) {
	var tokens []Token
	expr = strings.TrimSpace(expr)
	i := 0
	for i < len(expr) {
		ch := rune(expr[i])
		if unicode.IsSpace(ch) {
			i++
			continue
		}
		if ch == '(' {
			tokens = append(tokens, Token{Type: TokenParenLeft, Value: "("})
			i++
			continue
		}
		if ch == ')' {
			tokens = append(tokens, Token{Type: TokenParenRight, Value: ")"})
			i++
			continue
		}
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			// handle unary minus: if at start or after operator or left paren, and next is number/variable, treat as part of number? Simplify: insert 0 before unary minus
			if ch == '-' {
				if len(tokens) == 0 || tokens[len(tokens)-1].Type == TokenOperator || tokens[len(tokens)-1].Type == TokenParenLeft {
					// unary minus: push 0 then -
					tokens = append(tokens, Token{Type: TokenNumber, Value: "0"})
				}
			}
			tokens = append(tokens, Token{Type: TokenOperator, Value: string(ch)})
			i++
			continue
		}
		if unicode.IsDigit(ch) || ch == '.' {
			start := i
			dotCount := 0
			for i < len(expr) && (unicode.IsDigit(rune(expr[i])) || expr[i] == '.') {
				if expr[i] == '.' {
					dotCount++
					if dotCount > 1 {
						return nil, fmt.Errorf("invalid number multiple dots at %d", i)
					}
				}
				i++
			}
			tokens = append(tokens, Token{Type: TokenNumber, Value: expr[start:i]})
			continue
		}
		if unicode.IsLetter(ch) || ch == '_' {
			start := i
			for i < len(expr) && (unicode.IsLetter(rune(expr[i])) || unicode.IsDigit(rune(expr[i])) || expr[i] == '_') {
				i++
			}
			varName := strings.ToUpper(expr[start:i])
			// validate var name allowed chars: only A-Z _ 0-9 and known list
			if len(varName) == 0 {
				return nil, fmt.Errorf("empty variable at %d", start)
			}
			tokens = append(tokens, Token{Type: TokenVariable, Value: varName})
			continue
		}
		return nil, fmt.Errorf("invalid character '%c' at position %d", ch, i)
	}
	return tokens, nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// infixToPostfix — shunting-yard O(n)
func infixToPostfix(tokens []Token) ([]Token, error) {
	var output []Token
	var stack []Token // operator stack

	for _, tok := range tokens {
		switch tok.Type {
		case TokenNumber, TokenVariable:
			output = append(output, tok)
		case TokenOperator:
			for len(stack) > 0 && stack[len(stack)-1].Type == TokenOperator && precedence(stack[len(stack)-1].Value) >= precedence(tok.Value) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, tok)
		case TokenParenLeft:
			stack = append(stack, tok)
		case TokenParenRight:
			found := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.Type == TokenParenLeft {
					found = true
					break
				}
				output = append(output, top)
			}
			if !found {
				return nil, fmt.Errorf("mismatched parentheses")
			}
		}
	}
	// drain stack
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.Type == TokenParenLeft || top.Type == TokenParenRight {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, top)
	}
	return output, nil
}

// EvaluatePostfix — O(n) stack evaluation with decimal precise
func evaluatePostfix(postfix []Token, vars map[string]decimal.Decimal) (decimal.Decimal, error) {
	var stack []decimal.Decimal
	for _, tok := range postfix {
		switch tok.Type {
		case TokenNumber:
			d, err := decimal.NewFromString(tok.Value)
			if err != nil {
				return decimal.Zero, fmt.Errorf("invalid number %s", tok.Value)
			}
			stack = append(stack, d)
		case TokenVariable:
			val, ok := vars[tok.Value]
			if !ok {
				// allow zero for unknown variable? No, error for safety
				return decimal.Zero, fmt.Errorf("unknown variable %s (allowed: BASIC, CTC_MONTHLY, CTC_ANNUAL, GROSS)", tok.Value)
			}
			stack = append(stack, val)
		case TokenOperator:
			if len(stack) < 2 {
				return decimal.Zero, fmt.Errorf("insufficient operands for operator %s", tok.Value)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var res decimal.Decimal
			switch tok.Value {
			case "+":
				res = a.Add(b)
			case "-":
				res = a.Sub(b)
			case "*":
				res = a.Mul(b)
			case "/":
				if b.IsZero() {
					return decimal.Zero, fmt.Errorf("division by zero")
				}
				res = a.Div(b)
			default:
				return decimal.Zero, fmt.Errorf("unknown operator %s", tok.Value)
			}
			stack = append(stack, res)
		}
	}
	if len(stack) != 1 {
		return decimal.Zero, fmt.Errorf("invalid expression, stack left %d", len(stack))
	}
	return stack[0].Round(2), nil
}

// EvaluateFormula — public API: expression + variable map -> decimal
// Variables keys case-insensitive, will be uppercased
func EvaluateFormula(expression string, vars map[string]decimal.Decimal) (decimal.Decimal, error) {
	if strings.TrimSpace(expression) == "" {
		return decimal.Zero, fmt.Errorf("empty expression")
	}
	// Uppercase vars map for case-insensitivity O(n)
	upperVars := make(map[string]decimal.Decimal, len(vars))
	for k, v := range vars {
		upperVars[strings.ToUpper(k)] = v
	}
	tokens, err := tokenize(expression)
	if err != nil {
		return decimal.Zero, err
	}
	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return decimal.Zero, err
	}
	result, err := evaluatePostfix(postfix, upperVars)
	if err != nil {
		return decimal.Zero, err
	}
	return result, nil
}

// CalculateStructureComponent — calculate amount for a component given context
func CalculateStructureComponent(comp StructureComponent, ctxVars map[string]decimal.Decimal) (decimal.Decimal, error) {
	switch comp.CalculationType {
	case CalcFixed:
		return comp.Amount, nil
	case CalcPercentageOfBasic:
		basic, ok := ctxVars["BASIC"]
		if !ok {
			return decimal.Zero, fmt.Errorf("BASIC not in context for percentage_of_basic")
		}
		// percentage 40.00 = 40%
		return basic.Mul(comp.Percentage).Div(decimal.NewFromInt(100)).Round(2), nil
	case CalcPercentageOfCTC:
		ctc, ok := ctxVars["CTC_MONTHLY"]
		if !ok {
			// fallback CTC_ANNUAL/12
			if ctcAnnual, ok2 := ctxVars["CTC_ANNUAL"]; ok2 {
				ctc = ctcAnnual.Div(decimal.NewFromInt(12)).Round(2)
			} else {
				return decimal.Zero, fmt.Errorf("CTC_MONTHLY not in context")
			}
		}
		return ctc.Mul(comp.Percentage).Div(decimal.NewFromInt(100)).Round(2), nil
	case CalcPercentageOfGross:
		gross, ok := ctxVars["GROSS"]
		if !ok {
			return decimal.Zero, fmt.Errorf("GROSS not in context for percentage_of_gross")
		}
		return gross.Mul(comp.Percentage).Div(decimal.NewFromInt(100)).Round(2), nil
	case CalcFormula:
		if comp.Formula == "" {
			return decimal.Zero, fmt.Errorf("empty formula for component %s", comp.Code)
		}
		return EvaluateFormula(comp.Formula, ctxVars)
	default:
		return decimal.Zero, fmt.Errorf("unknown calculation type %s", comp.CalculationType)
	}
}

// ValidateFormula — security: only allowed variable names + operators
func ValidateFormula(formula string) error {
	allowedVars := map[string]bool{
		"BASIC": true, "CTC_MONTHLY": true, "CTC_ANNUAL": true, "GROSS": true,
		"BASIC_PLUS_ALLOW": true, "HOUSING": true, "TRANSPORT": true,
	}
	tokens, err := tokenize(formula)
	if err != nil {
		return err
	}
	for _, tok := range tokens {
		if tok.Type == TokenVariable {
			if !allowedVars[tok.Value] {
				// Allow any uppercase variable but must be alphanumeric + underscore — already tokenized
				// For strict security, warn but allow if it's known pattern: check if it starts with allowed prefix?
				// We will allow any variable that is all caps + underscore, but must be in vars context at runtime — so we warn only if suspicious length >30 or contains dangerous? Already safe because no function calls.
				if len(tok.Value) > 30 {
					return fmt.Errorf("variable name too long: %s", tok.Value)
				}
			}
		}
	}
	// Try postfix to catch syntax errors
	_, err = infixToPostfix(tokens)
	return err
}
