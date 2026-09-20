import type { ReactNode } from 'react'

interface FormFieldProps {
  label: string
  htmlFor: string
  children: ReactNode
  error?: string
  required?: boolean
}

function FormField({
  label,
  htmlFor,
  children,
  error,
  required = false,
}: FormFieldProps) {
  return (
    <div className="form-field">
      <label htmlFor={htmlFor}>
        {label}
        {required && <span className="form-field-required"> *</span>}
      </label>

      {children}

      {error && (
        <span
          id={`${htmlFor}-error`}
          className="form-field-error"
        >
          {error}
        </span>
      )}
    </div>
  )
}

export default FormField