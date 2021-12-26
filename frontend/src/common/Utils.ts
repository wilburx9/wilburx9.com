export function redact(s: string, size: number = 1): string {
  let max = Math.max(1, size)
  return `${s.slice(0, max)}****${s.slice(-max)}`
}