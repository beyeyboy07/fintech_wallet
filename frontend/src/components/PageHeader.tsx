import '../assets/css/PageHeader.css'

interface PageHeaderProps {
  title: string
  description?: string
}

function PageHeader({
  title,
  description,
}: PageHeaderProps) {
  return (
    <header className="page-header">
      <h2>{title}</h2>

      {description && <p>{description}</p>}
    </header>
  )
}

export default PageHeader