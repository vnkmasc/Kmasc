import Header from '@/components/common/header'

interface Props {
  children: React.ReactNode
}

const StudentLayout: React.FC<Props> = ({ children }) => {
  return (
    <main>
      <Header role='student' />
      <div className='container mt-16 py-6'>{children}</div>
    </main>
  )
}

export default StudentLayout
