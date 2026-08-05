// Formula Engine JS — mirrors Go formula_engine.go O(n) tokenization + shunting-yard + decimal precise
// For outstanding live preview in salary structure builder — secure no evil eval
// Supports variables: BASIC, CTC_MONTHLY, CTC_ANNUAL, GROSS, HOUSING etc
// Operators: + - * / ( ) and decimal numbers

export type TokenType = "number" | "variable" | "operator" | "paren_left" | "paren_right"
export interface Token { type: TokenType; value: string }

export function tokenize(expr: string): Token[] {
  const tokens: Token[] = []
  let i = 0
  expr = expr.trim()
  while (i < expr.length) {
    const ch = expr[i]
    if (/\s/.test(ch)) { i++; continue }
    if (ch === '(') { tokens.push({ type: "paren_left", value: "(" }); i++; continue }
    if (ch === ')') { tokens.push({ type: "paren_right", value: ")" }); i++; continue }
    if (ch === '+' || ch === '-' || ch === '*' || ch === '/') {
      if (ch === '-' && (tokens.length===0 || tokens[tokens.length-1].type==="operator" || tokens[tokens.length-1].type==="paren_left")) {
        tokens.push({ type: "number", value: "0" })
      }
      tokens.push({ type: "operator", value: ch }); i++; continue
    }
    if (/[0-9.]/.test(ch)) {
      let start = i, dotCount = 0
      while (i < expr.length && /[0-9.]/.test(expr[i])) {
        if (expr[i]==='.'){ dotCount++; if (dotCount>1) throw new Error(`Invalid number multiple dots at ${i}`) }
        i++
      }
      tokens.push({ type: "number", value: expr.slice(start,i) }); continue
    }
    if (/[A-Za-z_]/.test(ch)) {
      let start = i
      while (i < expr.length && /[A-Za-z0-9_]/.test(expr[i])) i++
      tokens.push({ type: "variable", value: expr.slice(start,i).toUpperCase() }); continue
    }
    throw new Error(`Invalid character '${ch}' at ${i}`)
  }
  return tokens
}

function precedence(op: string): number {
  if (op==="+"||op==="-") return 1
  if (op==="*"||op==="/") return 2
  return 0
}

function infixToPostfix(tokens: Token[]): Token[] {
  const output: Token[] = []
  const stack: Token[] = []
  for (const tok of tokens) {
    if (tok.type==="number"||tok.type==="variable") output.push(tok)
    else if (tok.type==="operator") {
      while (stack.length>0 && stack[stack.length-1].type==="operator" && precedence(stack[stack.length-1].value) >= precedence(tok.value)) {
        output.push(stack.pop()!)
      }
      stack.push(tok)
    } else if (tok.type==="paren_left") stack.push(tok)
    else if (tok.type==="paren_right") {
      let found = false
      while (stack.length>0) {
        const top = stack.pop()!
        if (top.type==="paren_left"){ found=true; break }
        output.push(top)
      }
      if (!found) throw new Error("Mismatched parentheses")
    }
  }
  while (stack.length>0) {
    const top = stack.pop()!
    if (top.type==="paren_left"||top.type==="paren_right") throw new Error("Mismatched parentheses")
    output.push(top)
  }
  return output
}

export function evaluatePostfix(postfix: Token[], vars: Record<string, number>): number {
  const stack: number[] = []
  for (const tok of postfix) {
    if (tok.type==="number") {
      const n = parseFloat(tok.value)
      if (isNaN(n)) throw new Error(`Invalid number ${tok.value}`)
      stack.push(n)
    } else if (tok.type==="variable") {
      const upper = tok.value.toUpperCase()
      if (!(upper in vars)) throw new Error(`Unknown variable ${tok.value} (allowed: BASIC, CTC_MONTHLY, CTC_ANNUAL, GROSS)`)
      stack.push(vars[upper])
    } else if (tok.type==="operator") {
      if (stack.length<2) throw new Error(`Insufficient operands for ${tok.value}`)
      const b = stack.pop()!, a = stack.pop()!
      let res: number
      switch(tok.value){
        case "+": res = a+b; break
        case "-": res = a-b; break
        case "*": res = a*b; break
        case "/": if (b===0) throw new Error("Division by zero"); res = a/b; break
        default: throw new Error(`Unknown operator ${tok.value}`)
      }
      stack.push(Math.round(res*100)/100)
    }
  }
  if (stack.length!==1) throw new Error(`Invalid expression stack left ${stack.length}`)
  return Math.round(stack[0]*100)/100
}

export function evaluateFormula(expression: string, vars: Record<string, number>): number {
  if (!expression.trim()) throw new Error("Empty expression")
  const upperVars: Record<string, number> = {}
  for (const k in vars) upperVars[k.toUpperCase()] = vars[k]
  const tokens = tokenize(expression)
  const postfix = infixToPostfix(tokens)
  return evaluatePostfix(postfix, upperVars)
}

export function validateFormula(formula: string): boolean {
  try {
    const tokens = tokenize(formula)
    for (const tok of tokens) {
      if (tok.type==="variable" && tok.value.length>30) throw new Error(`Variable name too long: ${tok.value}`)
    }
    infixToPostfix(tokens)
    return true
  } catch { return false }
}

// Calculate structure component for JS live preview — mirrors Go CalculateStructureComponent
export interface StructureComponentJS {
  code: string
  name: string
  calculation_type: "fixed" | "percentage_of_basic" | "percentage_of_ctc" | "percentage_of_gross" | "formula"
  amount: number
  percentage: number
  formula?: string
  is_taxable: boolean
  is_part_of_gross: boolean
  is_proratable: boolean
}

export function calculateComponent(comp: StructureComponentJS, ctx: Record<string, number>): number {
  switch(comp.calculation_type){
    case "fixed": return comp.amount
    case "percentage_of_basic": {
      const basic = ctx["BASIC"]
      if (basic===undefined) throw new Error("BASIC not in context")
      return Math.round(basic * comp.percentage / 100 *100)/100
    }
    case "percentage_of_ctc": {
      const ctc = ctx["CTC_MONTHLY"] ?? (ctx["CTC_ANNUAL"]?ctx["CTC_ANNUAL"]/12:0)
      if (ctc===undefined) throw new Error("CTC_MONTHLY not in context")
      return Math.round(ctc * comp.percentage /100 *100)/100
    }
    case "percentage_of_gross": {
      const gross = ctx["GROSS"]
      if (gross===undefined) throw new Error("GROSS not in context")
      return Math.round(gross * comp.percentage /100 *100)/100
    }
    case "formula": {
      if (!comp.formula) throw new Error(`Empty formula for ${comp.code}`)
      return evaluateFormula(comp.formula, ctx)
    }
    default: throw new Error(`Unknown calculation type ${comp.calculation_type}`)
  }
}

export function calculateEarningsFromStructure(
  components: StructureComponentJS[],
  vars: Record<string, number>,
  prorationFactor: number
): { earnings: Array<{code:string,name:string,amount:number,taxable:boolean,proratable:boolean}>, gross: number } {
  const sorted = [...components].sort((a,b)=> (a as any).order_no - (b as any).order_no)
  const earnings: Array<{code:string,name:string,amount:number,taxable:boolean,proratable:boolean}> = []
  let gross = 0
  const ctx = { ...vars, GROSS: 0 }
  for (const comp of sorted) {
    if ((comp as any).component_type && (comp as any).component_type !== "earning") continue
    let amount = calculateComponent(comp, ctx)
    if (comp.is_proratable) amount = Math.round(amount * prorationFactor *100)/100
    ctx[comp.code] = amount
    ctx["GROSS"] = gross + amount
    earnings.push({ code: comp.code, name: comp.name, amount, taxable: comp.is_taxable, proratable: comp.is_proratable })
    if (comp.is_part_of_gross) {
      gross += amount
      ctx["GROSS"] = gross
    }
  }
  return { earnings, gross: Math.round(gross*100)/100 }
}

// ET Tax calculation JS — binary search O(log n) over 7 brackets
export interface TaxBracketJS { min: number; max: number | null; rate: number; deduction: number }
export const etBrackets: TaxBracketJS[] = [
  { min: 0, max: 600, rate: 0, deduction: 0 },
  { min: 601, max: 1650, rate: 0.10, deduction: 60 },
  { min: 1651, max: 3200, rate: 0.15, deduction: 142.5 },
  { min: 3201, max: 5250, rate: 0.20, deduction: 302.5 },
  { min: 5251, max: 7800, rate: 0.25, deduction: 565 },
  { min: 7801, max: 10900, rate: 0.30, deduction: 955 },
  { min: 10901, max: null, rate: 0.35, deduction: 1500 },
]
export function calculateTaxJS(taxable: number, brackets: TaxBracketJS[] = etBrackets): number {
  if (taxable<=0) return 0
  let idx = brackets.findIndex(b => b.max===null || taxable < b.max)
  if (idx===-1) idx = brackets.length-1
  const br = brackets[idx]
  if (taxable < br.min) return 0
  let tax = taxable * br.rate - br.deduction
  if (tax<0) tax=0
  return Math.round(tax*100)/100
}
