import type { ReactNode } from 'react'

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