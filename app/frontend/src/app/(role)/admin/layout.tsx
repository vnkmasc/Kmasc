import Header from '@/components/common/header'

interface Props {
  children: React.ReactNode
}

const AdminLayout: React.FC<Props> = ({ children }) => {
  return (
    <main>
      <Header role='admin' />
      <div className='container mt-16 py-6'>{children}</div>
    </main>
  )
}

export default AdminLayout
