import type { ReactNode } from 'react'
import '../assets/css/Card.css'

interface CardProps {
  children: ReactNode
  className?: string
}

function Card({
  children,
  className = '',
}: CardProps) {
  return (
    <div className={`card ${className}`.trim()}>
      {children}
    </div>
  )
}

export default Card