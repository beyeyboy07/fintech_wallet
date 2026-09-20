import type { ReactNode } from 'react'

interface PageContainerProps {
  children: ReactNode
}

function PageContainer({
  children,
}: PageContainerProps) {
  return (
    <div className="page-container">
      {children}
    </div>
  )
}

export default PageContainer