import type { ReactNode } from 'react'
import '../assets/css/FormActions.css'

interface FormActionsProps {
  children: ReactNode
}

function FormActions({
  children,
}: FormActionsProps) {
  return (
    <div className="form-actions">
      {children}
    </div>
  )
}

export default FormActions