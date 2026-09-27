import type { TextareaHTMLAttributes } from 'react'
import '../assets/css/Textarea.css'

type TextareaProps =
  TextareaHTMLAttributes<HTMLTextAreaElement>

function Textarea({
  className = '',
  ...props
}: TextareaProps) {
  return (
    <textarea
      className={`textarea ${className}`.trim()}
      {...props}
    />
  )
}

export default Textarea