import type { InputHTMLAttributes } from 'react'
import '../assets/css/Input.css'


type InputProps = InputHTMLAttributes<HTMLInputElement>

function Input({
  className = '',
  ...props
}: InputProps) {
  return (
    <input
      className={`input ${className}`.trim()}
      {...props}
    />
  )
}

export default Input