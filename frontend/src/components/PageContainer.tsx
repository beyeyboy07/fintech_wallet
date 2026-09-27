import type { ReactNode } from 'react'
import '../assets/css/PageContainer.css'

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